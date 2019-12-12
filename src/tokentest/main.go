package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type token struct {
	AccessToken string `json:"AccessToken"`
}

func tokenRead(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var inputToken token
	body := []byte(req.Body)

	err := json.Unmarshal(body, &inputToken)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	if inputToken.AccessToken == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "You forgot something mate",
		}, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	// Create Cognito service client

	cognitoClient := cognitoidentityprovider.New(sess)

	params := cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(inputToken.AccessToken),
	}

	userInfo, err := cognitoClient.GetUser(&params)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	jsondata, _ := json.Marshal(userInfo)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsondata),
	}, nil

}

func main() {
	lambda.Start(tokenRead)
}
