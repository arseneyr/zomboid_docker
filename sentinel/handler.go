package main

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type lambdaHandler func(context.Context, events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error)

func eventToRequest(ctx context.Context, event events.LambdaFunctionURLRequest) (req *http.Request, err error) {
	u := url.URL{
		Scheme:   "https",
		Host:     event.RequestContext.DomainName,
		RawPath:  event.RawPath,
		RawQuery: event.RawQueryString,
	}
	p, err := url.PathUnescape(event.RawPath)
	if err != nil {
		return
	}
	u.Path = p
	var body io.Reader = strings.NewReader(event.Body)
	if event.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}
	req, err = http.NewRequestWithContext(ctx, event.RequestContext.HTTP.Method, u.String(), body)
	return
}

func responseRecorderToEvent(res *httptest.ResponseRecorder) events.LambdaFunctionURLResponse {
	result := res.Result()
	headers := make(map[string]string, len(result.Header))
	for k, v := range result.Header {
		headers[k] = strings.Join(v, ", ")
	}
	return events.LambdaFunctionURLResponse{
		StatusCode:      result.StatusCode,
		Headers:         headers,
		IsBase64Encoded: true,
		Body:            base64.StdEncoding.EncodeToString(res.Body.Bytes()),
	}
}

func handlerToLambda(handler http.Handler) lambdaHandler {
	return func(ctx context.Context, event events.LambdaFunctionURLRequest) (res events.LambdaFunctionURLResponse, err error) {
		r := httptest.NewRecorder()
		req, err := eventToRequest(ctx, event)
		if err != nil {
			return
		}
		handler.ServeHTTP(r, req)
		res = responseRecorderToEvent(r)
		return
	}
}
