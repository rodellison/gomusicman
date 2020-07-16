package clients

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rodellison/gomusicman/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
		DynamoDBSvcClient = &mocks.MockDynamoDBSvcClient{}
}

//This test returns a known valid replacement for an item
func TestRequestSongKickFoundItem(t *testing.T) {

	mocks.MockDynamoGetItem = func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: nil,
			Item: map[string]*dynamodb.AttributeValue{
				"SongKickInvalidParm": &dynamodb.AttributeValue{
					S: aws.String("deaf leopard"),
				},
				"SongKickValidParm": &dynamodb.AttributeValue{
					S: aws.String("Def Leppard"),
				},
			},
		}, nil
	}

	expectedVal := "Def Leppard"
	strArtist := QueryMusicManParmTable("deaf leopard")
	assert.Equal(t, expectedVal, strArtist)

}

//This test ensures that if a replacement value is NOT found, to just return the original passed in value
func TestRequestSongKickDidNotFindItem(t *testing.T) {

	mocks.MockDynamoGetItem = func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: nil,
			Item: map[string]*dynamodb.AttributeValue{
				"SongKickInvalidParm": &dynamodb.AttributeValue{
					S: aws.String("YouWontFindMe"),
				},
				"SongKickValidParm": &dynamodb.AttributeValue{
					S: aws.String(""),
				},
			},
		}, nil
	}

	expectedVal := "YouWontFindMe"
	strArtist := QueryMusicManParmTable("YouWontFindMe")
	assert.Equal(t, expectedVal, strArtist)

}