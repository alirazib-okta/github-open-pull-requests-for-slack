package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	responseBody, err := GetListOfPRs()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       responseBody,
	}

	return response, nil
}

func main() {
	lambda.Start(handler)
}
