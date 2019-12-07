package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-central-1"))

/*
func scantest() (*dynamodb.ScanOutput, error) {
	// Prepare the input for the query.
	input := &dynamodb.ScanInput{
		TableName: aws.String("unimeret-location"),

		//FilterExpression: aws.String("begins_with(hash, :sw)"),
	}

	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	if result.Items == nil {
		return nil, nil
	}

	// The result.Item object returned has the underlying type
	// map[string]*AttributeValue. We can use the UnmarshalMap helper
	// to parse this straight into the fields of a struct. Note:
	// UnmarshalListOfMaps also exists if you are working with multiple
	// items.
	//loc := new(location)

	//err = dynamodbattribute.UnmarshalMap(result.Items, link)
	//if err != nil {
	//	return nil, err
	//}


	return result, nil
}
*/
