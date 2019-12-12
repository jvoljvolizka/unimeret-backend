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

func reterr(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func changePass(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var inputReq changePassRequest
	body := []byte(req.Body)

	err := json.Unmarshal(body, &inputReq)

	if err != nil {
		return reterr(err)
	}

	if inputReq.AccessToken == "" || inputReq.OldPass == "" || inputReq.NewPass == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error" : "empty field"}`,
		}, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		return reterr(err)
	}

	// Create Cognito service client

	cognitoClient := cognitoidentityprovider.New(sess)

	params := cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(inputReq.AccessToken),
		PreviousPassword: aws.String(inputReq.OldPass),
		ProposedPassword: aws.String(inputReq.NewPass),
	}

	_, err = cognitoClient.ChangePassword(&params)

	if err != nil {
		return reterr(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{ "result" : "success" }`,
	}, nil

}

func main() {
	lambda.Start(changePass)
}
