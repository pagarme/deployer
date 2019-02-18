package logger

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CommandLog struct {
	Username     string   `json:"username"`
	Timestamp    string   `json:"timestamp"`
	Command      string   `json:"command"`
	Args         []string `json:"args"`
	Status       string   `json:"status"`
	StatusReason string   `json:"statusReason"`
	ExecutionID  string   `json:"executionId"`
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
