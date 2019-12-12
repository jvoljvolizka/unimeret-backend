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

type token struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}

func SecretHash(username, clientID, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func refresh(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	const userPoolID = "eu-central-1_PeQQEcL03"
	const appClientId = "5n3ndtvdqbgfvnagbrjtlfdnvd"
	const clientSecret = "ep6j8m81mjjpj02olb70q8oucdmj4028gqtvk1an51esukpg8o6"

	var inToken token
	body := []byte(req.Body)

	err := json.Unmarshal(body, &inToken)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
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
		AccessToken: aws.String(inToken.AccessToken),
	}

	userInfo, err := cognitoClient.GetUser(&params)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	killme := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(inToken.RefreshToken),
			"SECRET_HASH":   aws.String(SecretHash(*userInfo.Username, appClientId, clientSecret)),
		},
		ClientId: aws.String(appClientId),
	}

	//newUserData.SetUserPoolId(userPoolID)
	//

	//cognitoClient.
	authResponse, err := cognitoClient.InitiateAuth(&killme)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	jsondata, _ := json.Marshal(authResponse)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsondata),
	}, nil

}

func main() {
	lambda.Start(refresh)
}
