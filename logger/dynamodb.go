package logger

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoLogger struct {
	Svc   *dynamodb.DynamoDB
	Table *string
}

func (d *DynamoLogger) Init() {
	awsRegion, ok := os.LookupEnv("DEPLOYER_AWS_REGION")
	if !ok {
		fmt.Printf("DEPLOYER_AWS_REGION must be defined")
		os.Exit(1)
	}

	dynamoTable, ok := os.LookupEnv("DEPLOYER_DYNAMODB_TABLE")
	if !ok {
		fmt.Printf("DEPLOYER_DYNAMODB_TABLE must be defined")
		os.Exit(1)
	}

	config := &aws.Config{
		Region: aws.String(awsRegion),
	}

	ses, err := session.NewSession(config)
	if err != nil {
		fmt.Printf("An error occured when starting the AWS Session: %s\n", err)
		os.Exit(1)
	}

	d.Svc = dynamodb.New(ses)
	d.Table = aws.String(dynamoTable)
}

func (d *DynamoLogger) LogCommand(command CommandLog) error {
	av, err := dynamodbattribute.MarshalMap(command)
	if err != nil {
		fmt.Printf("An error occured when parsing command: %s\n", err)
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: d.Table,
	}

	_, err = d.Svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return nil
}
