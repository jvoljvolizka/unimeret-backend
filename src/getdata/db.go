package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-central-1"))

func getItem(urlid string) (*shortURL, error) {
	// Prepare the input for the query.
	input := &dynamodb.GetItemInput{
		TableName: aws.String("unimeret-location"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(urlid),
			},
		},
	}

	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	// The result.Item object returned has the underlying type
	// map[string]*AttributeValue. We can use the UnmarshalMap helper
	// to parse this straight into the fields of a struct. Note:
	// UnmarshalListOfMaps also exists if you are working with multiple
	// items.
	loc := new(location)
	err = dynamodbattribute.UnmarshalMap(result.Item, link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

/*
func putItem(link *shortURL) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("URLs"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(link.ID),
			},
			"URL": {
				S: aws.String(link.URL),
			},
			"DelToken": {
				S: aws.String(link.DelToken),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}*/
