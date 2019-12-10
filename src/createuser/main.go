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

func createUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	const userPoolID = "eu-central-1_PeQQEcL03"
	const appClientId = "5n3ndtvdqbgfvnagbrjtlfdnvd"
	const clientSecret = "ep6j8m81mjjpj02olb70q8oucdmj4028gqtvk1an51esukpg8o6"

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
	password := inputUser.Password

	if emailID == "" || userName == "" || password == "" {
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

	newUserData := &cognitoidentityprovider.SignUpInput{

		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(emailID),
			},
		},
	}

	newUserData.SetClientId(appClientId)
	newUserData.SetUsername(userName)
	newUserData.SetPassword(password)
	newUserData.SetSecretHash(SecretHash(userName, appClientId, clientSecret))

	//newUserData.SetUserPoolId(userPoolID)
	//

	//cognitoClient.
	_, err = cognitoClient.SignUp(newUserData) //AdminCreateUser(newUserData)
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
