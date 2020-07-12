package mocks

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDBSvcClient struct {
	dynamodbiface.DynamoDBAPI
}

var (
	MockDynamoScan func(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
)

//This is the mocked version of the real function
//It returns the variable above, which is a function that can be overloaded in our test routines
func (m *MockDynamoDBSvcClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return MockDynamoScan(input)
}
