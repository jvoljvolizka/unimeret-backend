package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SecretHash(username, clientID, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func authUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	const userPoolID = "eu-central-1_PeQQEcL03"
	const appClientId = "5n3ndtvdqbgfvnagbrjtlfdnvd"
	const clientSecret = "ep6j8m81mjjpj02olb70q8oucdmj4028gqtvk1an51esukpg8o6"

	var inputUser user
	body := []byte(req.Body)

	err := json.Unmarshal(body, &inputUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	userName := inputUser.Username
	password := inputUser.Password

	if userName == "" || password == "" {
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

	params := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(userName),
			"PASSWORD":    aws.String(password),
			"SECRET_HASH": aws.String(SecretHash(userName, appClientId, clientSecret)),
		},
		ClientId: aws.String(appClientId),
	}

	//newUserData.SetUserPoolId(userPoolID)
	//

	//cognitoClient.
	authResponse, err := cognitoClient.InitiateAuth(&params)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       authResponse.String(),
	}, nil

}

func main() {
	lambda.Start(authUser)
}
