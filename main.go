package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := Helper()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	fmt.Printf("Response from webhook: %s\n", resp)
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       resp, // we still include the list of PRs in the response.
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
