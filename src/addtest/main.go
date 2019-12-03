package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcloughlin/geohash"
)

type location struct {
	Hash  string `json:"hash"`
	Locid string `json:"locid"`
}

type inlocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

const precision int = 6

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

	hash := geohash.Encode(input.Lat, input.Lng)

	var newloc location
	newloc.Hash = hash
	newloc.Locid = "T-" + hash

	err = putItem(&newloc)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}
	//6 150m~
	//hash := fmt.Sprintf("%f", input.lat)
	//hash := input.lat

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
		Body:       "got it ! btw your life is still just a sad void",
	}, nil
}

func main() {
	lambda.Start(getdata)
}
