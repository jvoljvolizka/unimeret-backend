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

const precision int = 6

func getdata(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var input inlocation
	body := []byte(req.Body)

	err := json.Unmarshal(body, &input)

	if err != nil {
		return reterr(err)
	}

	hash := geohash.Encode(input.Lat, input.Lng)

	sendhash := hash[:precision]
	locations, err := scantest(sendhash)

	if err != nil {
		return reterr(err)
	}

	var resp string
	for _, loc := range locations {
		resp = resp + " - " + loc.Hash
		//resp = append(resp, loc.Hash)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       resp,
	}, nil
}

func main() {
	lambda.Start(getdata)
}
