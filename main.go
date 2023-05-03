package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get the list of PRs
	listOfPRs, err := GetListOfPRs()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Post to slack webhook
	resp, err := PostToSlackWebhook(listOfPRs)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	fmt.Printf("Response from webhook: %s\n", resp)
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       listOfPRs, // we still include the list of PRs in the response.
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
