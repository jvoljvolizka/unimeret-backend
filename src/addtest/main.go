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

func reterr(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func getdata(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var input inlocation
	body := []byte(req.Body)

	err := json.Unmarshal(body, &input)

	if err != nil {
		return reterr(err)
	}

	hash := geohash.Encode(input.Lat, input.Lng)

	var newloc location
	newloc.Hash = hash
	newloc.Locid = "T-" + hash

	err = putItem(&newloc)

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
