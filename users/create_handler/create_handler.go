package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/laaksomavrick/twitter-go/common"
)

type createUserDto struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// Can this take the createUserDto directly as an arg?
func Handler(ctx context.Context, event common.Event) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer
	var dto createUserDto

	// does this replace the usage of a validator?
	decoder := json.NewDecoder(strings.NewReader(event.Body))

	err := decoder.Decode(&dto)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404}, err
	}

	json.NewEncoder(&buf).Encode(dto)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
