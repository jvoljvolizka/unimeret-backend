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

type newUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func createUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	userPoolID := "eu-central-1_PeQQEcL03"

	var inputUser newUser
	body := []byte(req.Body)

	err := json.Unmarshal(body, &inputUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	userName := inputUser.Username
	emailID := inputUser.Email

	if emailID == "" || userName == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Email or Username can't be empty",
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

	newUserData := &cognitoidentityprovider.AdminCreateUserInput{
		DesiredDeliveryMediums: []*string{
			aws.String("EMAIL"),
		},
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(emailID),
			},
		},
	}

	newUserData.SetUserPoolId(userPoolID)
	newUserData.SetUsername(userName)

	_, err = cognitoClient.AdminCreateUser(newUserData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "User created",
	}, nil

}

func main() {
	lambda.Start(createUser)
}
