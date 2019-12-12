package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

/*
type storeEvent struct {
	EventName string `json:"EventName"`
	MaxPerson string `json:"MaxPerson"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
	Body      string `json:"Body"`
	UserID    string `json:"UserID"`
	lochash   string `json:"lochash"`
	locID     string `json:"locID"`
}

*/

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-central-1"))

func putItem(inputEvent *storeEvent) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("unimeret-events"),
		Item: map[string]*dynamodb.AttributeValue{
			"EventName": {
				S: aws.String(inputEvent.EventName),
			},
			"MaxPerson": {
				S: aws.String(inputEvent.MaxPerson),
			},
			"StartDate": {
				S: aws.String(inputEvent.StartDate),
			},
			"EndDate": {
				S: aws.String(inputEvent.EndDate),
			},
			"Body": {
				S: aws.String(inputEvent.Body),
			},
			"userid": {
				S: aws.String(inputEvent.UserID),
			},
			"lochash": {
				S: aws.String(inputEvent.Lochash),
			},
			"locID": {
				S: aws.String(inputEvent.LocID),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}
