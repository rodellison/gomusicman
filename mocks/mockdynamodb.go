package mocks

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDBSvcClient struct {
	dynamodbiface.DynamoDBAPI
}

var (
	MockDynamoGetItem func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
)

//This is the mocked version of the real function
//It returns the variable above, which is a function that can be overloaded in our test routines
func (m *MockDynamoDBSvcClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return MockDynamoGetItem(input)
}
