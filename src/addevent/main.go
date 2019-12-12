package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/mmcloughlin/geohash"
)

//| userid     | eventname      | max_person _nullable_ |start_date | finish_date | body |

type inEvent struct {
	EventName   string  `json:"EventName"`
	MaxPerson   string  `json:"MaxPerson"`
	StartDate   string  `json:"StartDate"`
	EndDate     string  `json:"EndDate"`
	Body        string  `json:"Body"`
	AccessToken string  `json:"AccessToken"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

type storeEvent struct {
	EventName string `json:"EventName"`
	MaxPerson string `json:"MaxPerson"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
	Body      string `json:"Body"`
	UserID    string `json:"UserID"`
	Lochash   string `json:"lochash"`
	LocID     string `json:"locID"`
}

func reterr(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func getdata(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var input inEvent
	body := []byte(req.Body)

	err := json.Unmarshal(body, &input)

	if err != nil {
		return reterr(err)
	}

	/*	if input.AccessToken == "" || input.Body == "" || input.EndDate == "" || input.EventName == "" || input.Lat == "" || input.Lng == "" || input.MaxPerson == "" || input.StartDate == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "You forgot something mate",
			}, nil
		}
	*/

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		return reterr(err)
	}

	locHash := geohash.Encode(input.Lat, input.Lng)

	cognitoClient := cognitoidentityprovider.New(sess)

	params := cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(input.AccessToken),
	}

	userInfo, err := cognitoClient.GetUser(&params)

	if err != nil {
		return reterr(err)
	}

	var newEvent storeEvent

	newEvent.Body = input.Body
	newEvent.EndDate = input.EndDate
	newEvent.EventName = input.EventName
	newEvent.MaxPerson = input.MaxPerson
	newEvent.StartDate = input.StartDate
	newEvent.UserID = *userInfo.UserAttributes[0].Value
	newEvent.Lochash = locHash
	newEvent.LocID = *userInfo.UserAttributes[0].Value + "-" + locHash

	err = putItem(&newEvent)

	if err != nil {
		return reterr(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{ "result" : "success" }`,
	}, nil
}

func main() {
	lambda.Start(getdata)
}
