package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

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

func reterr(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func authUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//userPoolID := os.Getenv("USERPOOL")
	appClientId := os.Getenv("CLIENTID")
	clientSecret := os.Getenv("CLIENTSECRET")

	var inputUser user
	body := []byte(req.Body)

	err := json.Unmarshal(body, &inputUser)

	if err != nil {
		return reterr(err)
	}

	if inputUser.Username == "" || inputUser.Password == "" {
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

	params := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(inputUser.Username),
			"PASSWORD":    aws.String(inputUser.Password),
			"SECRET_HASH": aws.String(SecretHash(inputUser.Username, appClientId, clientSecret)),
		},
		ClientId: aws.String(appClientId),
	}

	authResponse, err := cognitoClient.InitiateAuth(&params)

	if err != nil {
		return reterr(err)
	}

	jsondata, err := json.Marshal(authResponse)

	if err != nil {
		return reterr(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsondata),
	}, nil

}

func main() {
	lambda.Start(authUser)
}
