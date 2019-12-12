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

type newUser struct {
	Email    string `json:"email"`
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

func createUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//userPoolID := os.Getenv("USERPOOL")
	appClientId := os.Getenv("CLIENTID")
	clientSecret := os.Getenv("CLIENTSECRET")

	var inputUser newUser
	body := []byte(req.Body)

	err := json.Unmarshal(body, &inputUser)

	if err != nil {
		return reterr(err)
	}

	if inputUser.Email == "" || inputUser.Username == "" || inputUser.Password == "" {
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

	newUserData := &cognitoidentityprovider.SignUpInput{

		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(inputUser.Email),
			},
		},
	}

	newUserData.SetClientId(appClientId)
	newUserData.SetUsername(inputUser.Username)
	newUserData.SetPassword(inputUser.Password)
	newUserData.SetSecretHash(SecretHash(inputUser.Username, appClientId, clientSecret))

	_, err = cognitoClient.SignUp(newUserData) //AdminCreateUser(newUserData)

	if err != nil {
		return reterr(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{ "result" : "success" }`,
	}, nil

}

func main() {
	lambda.Start(createUser)
}
