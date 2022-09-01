package main

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestInstanceStatus(t *testing.T) {
	initAWSClient()
	assert.True(t, isInstanceRunning(context.TODO()))
}

func TestTemplate(t *testing.T) {
	var buf bytes.Buffer
	generateHtml(&buf, false)
	t.Log(buf.String())
}

func TestHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	handler := handlerToLambda(mux)

	tests := []struct {
		request events.LambdaFunctionURLRequest
		expect  events.LambdaFunctionURLResponse
	}{
		{
			request: events.LambdaFunctionURLRequest{
				RawPath:        "/",
				RawQueryString: "",
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "GET",
					},
					DomainName: "yo",
				},
			},
			expect: events.LambdaFunctionURLResponse{
				StatusCode: http.StatusOK,
			},
		},
	}

	for _, test := range tests {
		res, _ := handler(context.TODO(), test.request)
		assert.Equal(t, test.expect, res)
	}
}
