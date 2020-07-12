package clients

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"os"
)

var (
	DynamoDBSvcClient dynamodbiface.DynamoDBAPI
	TableName         string
)

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


//func getEventIDsForOldEvents takes an input endData (of form 20200101), and scans for items in the DynamoDB table
//where the Event EndData is prior to the input and returns this collection .
func getEventIDsForOldEvents(endDate string) (returnItems *dynamodb.ScanOutput, err error) {

	// Create the Expression to fill the scan input struct with.
	// Get all events whos EndDate is less than, (earlier) that the endDate string provided. This effectively gets all the items
	//whos event has already happened. They are the ones to be purged...
	filt := expression.Name("EndDate").LessThan(expression.Value(endDate))

	//Create a projection, to get back particular attributes..
	proj := expression.NamesList(expression.Name("EventID"), expression.Name("StartDate"), expression.Name("EndDate"))

	//Now build the expression of the projection we want, with the filter applied
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		fmt.Println("Got error building Scan input expression:")
		fmt.Println(err.Error())
		return nil, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(TableName),
	}

	// Make the DynamoDB Query API call
	result, err := DynamoDBSvcClient.Scan(params)
	if err != nil {
		fmt.Println("DynamoDb Scan Query API call failed:")
		fmt.Println((err.Error()))
		return nil, err
	}

	return result, nil
}
