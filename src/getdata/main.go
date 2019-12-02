package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcloughlin/geohash"
)

type location struct {
	lat   string `json:"lat"`
	lon   string `json:"lon"`
	locid string `json:"locid"`
}

type inlocation struct {
	lat float64 `json:"lat"`
	lng float64 `json:"lng"`
}

func getdata(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input inlocation
	body := []byte(req.Body)

	err := json.Unmarshal(body, &input)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	hash := geohash.Encode(input.lat, input.lng)

	//radius := 100

	/*check, _ := getItem(newURL.ID)

	if check != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "id already exists",
		}, nil
	}

	err = putItem(&newURL)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}*/

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       hash,
	}, nil
}

func main() {
	lambda.Start(getdata)
}
