package logger

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

var (
	ErrorDeployerRegionUndefined      = errors.New("DEPLOYER_AWS_REGION must be defined")
	ErrorDeployerDynamoTableUndefined = errors.New("DEPLOYER_DYNAMODB_TABLE must be defined")
	ErrorOccurredWhenStartingSession  = errors.New("An error occurred when starting the AWS Session")
	ErrorParsingCommands              = errors.New("An error occurred when parsing command")
	ErrorCallingPutItem               = errors.New("Got error calling PutItem")
)

type DynamoLogger struct {
	Svc   *dynamodb.DynamoDB
	Table *string
}

func (d *DynamoLogger) LogCommand(command CommandLog) error {
	av, err := dynamodbattribute.MarshalMap(command)
	if err != nil {
		return errors.Wrap(err, ErrorParsingCommands.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: d.Table,
	}

	_, err = d.Svc.PutItem(input)
	if err != nil {
		return errors.Wrap(err, ErrorCallingPutItem.Error())
	}

	return nil
}

func NewDynamoLogger() (*DynamoLogger, error) {
	awsRegion, ok := os.LookupEnv("DEPLOYER_AWS_REGION")
	if !ok {
		return nil, ErrorDeployerRegionUndefined
	}

	dynamoTable, ok := os.LookupEnv("DEPLOYER_DYNAMODB_TABLE")
	if !ok {
		return nil, ErrorDeployerDynamoTableUndefined
	}

	config := &aws.Config{
		Region: aws.String(awsRegion),
	}

	ses, err := session.NewSession(config)
	if err != nil {
		return nil, errors.Wrap(err, ErrorOccurredWhenStartingSession.Error())
	}

	d := &DynamoLogger{}
	d.Svc = dynamodb.New(ses)
	d.Table = aws.String(dynamoTable)

	return d, nil
}
