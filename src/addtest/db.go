package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-central-1"))

func putItem(loc *location) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("unimeret-location"),
		Item: map[string]*dynamodb.AttributeValue{
			"locid": {
				S: aws.String(loc.Locid),
			},
			"hash": {
				S: aws.String(loc.Hash),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}
