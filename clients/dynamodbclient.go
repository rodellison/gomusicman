package clients

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"os"
)

var (
	DynamoDBSvcClient dynamodbiface.DynamoDBAPI
	TableName         string
)

type MusicManParm struct {
	SongKickInvalidParm string
	SongKickValidParm string
}

func init() {

	//During testing, we'll override the endpoint to ensure testing against local DynamoDB Docker image
	cfg := aws.Config{
		//		Endpoint: aws.String("http://localhost:8000"),
		Region:     aws.String("us-east-1"),
		MaxRetries: aws.Int(3),
	}

	//Get Session, credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create the eventbridge events service client, to be used for putting events
	DynamoDBSvcClient = dynamodb.New(sess, &cfg)

	//Making the Tablename an environmental variable so that it can be changed outside of program
	TableName = os.Getenv("DYNAMO_DB_TABLENAME")

}

func QueryMusicManParmTable( strArtistValue string) (string) {

	strDynamoDBTableName :=os.Getenv("DYNAMO_DB_TABLENAME");

	params := &dynamodb.GetItemInput{
		TableName:                 aws.String(strDynamoDBTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"SongKickInvalidParm": {
				S: aws.String(strArtistValue),
			},
		},
	}

	result, err := DynamoDBSvcClient.GetItem(params)
	if err != nil {
		fmt.Println(err.Error())
		return strArtistValue
	}

	var MusicManParmResult MusicManParm
	err = dynamodbattribute.UnmarshalMap(result.Item, &MusicManParmResult)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	} else {
		//Check if there is a Valid Entry provided from the Result - i.e. We got a hit, so lets swap out the bad value for the good.
		if MusicManParmResult.SongKickValidParm != "" {
			return MusicManParmResult.SongKickValidParm
		} else {
			return strArtistValue

		}
	}

}
