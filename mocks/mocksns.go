package mocks

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type MockSNSSvcClient struct {
	snsiface.SNSAPI
}

var (
	MockDoPublishEvent func(input *sns.PublishInput) (*sns.PublishOutput, error)
)

//This is the mocked version of the real function
//It returns the variable above, which is a function that can be overloaded in our test routines
func (s *MockSNSSvcClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	return MockDoPublishEvent(input)
}
