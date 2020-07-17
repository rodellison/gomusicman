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
	SongKickValidParm   string
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
	// Create the dynamodb service client, to be used for querying for artist or venue input corrections
	DynamoDBSvcClient = dynamodb.New(sess, &cfg)

	//Making the Tablename an environmental variable so that it can be changed outside of program
	TableName = os.Getenv("DYNAMO_DB_TABLENAME")

}

func QueryMusicManParmTable(strValue string) string {

	fmt.Println("Checking MusicManParmTable for entry: ", strValue)
	params := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"SongKickInvalidParm": {
				S: aws.String(strValue),
			},
		},
	}

	result, err := DynamoDBSvcClient.GetItem(params)
	if err != nil {
		fmt.Println("Error performing DynamoDB GetItem: ", err.Error())
		return strValue
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
			return strValue

		}
	}

}
