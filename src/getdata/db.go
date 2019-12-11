package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-central-1"))

//	svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-central-1"))

// snippet-start:[dynamodb.go.scan_items.vars]

// snippet-end:[dynamodb.go.scan_items]

func scantest(startHash string) ([]*location, error) {

	tableName := "unimeret-location"

	// Create the Expression to fill the input struct with.
	// Get all movies in that year; we'll pull out those with a higher rating later
	filt := expression.Name("hash").BeginsWith(startHash)

	proj := expression.NamesList(expression.Name("hash"), expression.Name("locid"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := db.Scan(params)
	if err != nil {
		return nil, err
	}

	numItems := 0

	var locations []*location
	for _, i := range result.Items {
		item := location{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return nil, err
		}

		// Which ones had a higher rating than minimum?
		// Or it we had filtered by rating previously:
		//   if item.Year == year {
		numItems++
		locations = append(locations, &item)

	}
	return locations, nil
}
