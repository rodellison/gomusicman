package clients

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

var (
	SNSSvcClient snsiface.SNSAPI
)

func init() {

	//Get Session, credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create the eventbridge events service client, to be used for putting events
	SNSSvcClient = sns.New(sess)

}

// func PublishSNSMessage uses an SDK service client to send an SNS Publish request
func PublishSNSMessage(snsTopic, snsSubject, snsMessage string) (err error) {

	pubInput := &sns.PublishInput{
		Message:  aws.String(snsMessage),
		Subject:  aws.String(snsSubject),
		TopicArn: aws.String(snsTopic),
	}

	_, err = SNSSvcClient.Publish(pubInput)
	if err != nil {
		return err
	} else {
		return nil
	}

}
