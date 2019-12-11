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

type changePassRequest struct {
	AccessToken string `json:"AccessToken"`
	OldPass     string `json:"oldpass"`
	NewPass     string `json:"newpass"`
}

func changePass(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var inputReq changePassRequest
	body := []byte(req.Body)

	err := json.Unmarshal(body, &inputReq)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	if inputReq.AccessToken == "" || inputReq.OldPass == "" || inputReq.NewPass == "" {
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

	params := cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(inputReq.AccessToken),
		PreviousPassword: aws.String(inputReq.OldPass),
		ProposedPassword: aws.String(inputReq.NewPass),
	}

	passChange, err := cognitoClient.ChangePassword(&params)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       passChange.String(),
	}, nil

}

func main() {
	lambda.Start(changePass)
}
