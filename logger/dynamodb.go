package logger

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		fmt.Printf("An error occured when starting the AWS Session: %s\n", err)
		os.Exit(1)
	}

	d.Svc = dynamodb.New(session)
	d.Table = aws.String(dynamoTable)
}
