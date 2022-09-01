package main

import (
	"context"
	_ "embed"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

//go:embed index.html
var htmlTemplateString string
var htmlTemplate *template.Template = template.Must(template.New("page").Parse(htmlTemplateString))
var ec2Client *ec2.Client
var ec2InstanceId string

func initAWSClient() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	ec2Client = ec2.NewFromConfig(cfg)
	instance, set := os.LookupEnv("INSTANCE_ID")
	if !set {
		panic("INSTANCE_ID not set")
	}
	ec2InstanceId = instance
}

func isInstanceRunning(ctx context.Context) bool {
	output, err := ec2Client.DescribeInstanceStatus(ctx, &ec2.DescribeInstanceStatusInput{
		InstanceIds: []string{ec2InstanceId},
	})
	if err != nil {
		panic(err)
	}
	for _, status := range output.InstanceStatuses {
		if status.InstanceState.Name == ec2Types.InstanceStateNamePending || status.InstanceState.Name == ec2Types.InstanceStateNameRunning {
			return true
		}
	}
	return false
}

func startInstance(ctx context.Context) {
	_, err := ec2Client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: []string{ec2InstanceId},
	})
	if err != nil {
		log.Println("unable to start instance", err)
	}
}

func generateHtml(writer io.Writer, started bool) {
	templateData := struct {
		Checked string
	}{
		Checked: "",
	}
	if started {
		templateData.Checked = "checked"
	}
	err := htmlTemplate.Execute(writer, templateData)
	if err != nil {
		panic(err)
	}
}

func handleRoot(resWriter http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet && req.URL.Path == "/" {
		generateHtml(resWriter, isInstanceRunning(req.Context()))
	} else if req.Method == http.MethodPost && req.URL.Path == "/start" {
		startInstance(req.Context())
	} else {
		http.NotFound(resWriter, req)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	initAWSClient()
	lambda.Start(handlerToLambda(mux))
}
