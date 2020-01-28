package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler handles things
type Handler struct {
}

// Handle an event
func (h *Handler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	handler := Handler{}
	lambda.Start(handler.Handle)
}

// To build zip file:

//rho:examples rho$ GOOS=linux go build -o sample-function sample-function.go
// rho:examples rho$ zip sample-function.zip sample-function
