package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"log"
)

const (
	awsRegion = "us-west-2"
)

func main() {
	// Create a new session
	sess := session.Must(session.NewSession())

	// Create SFN New Client
	sfnClient := sfn.New(sess, aws.NewConfig().WithRegion(awsRegion))

	// SFN Arn
	stateMachinearn := "arn:aws:states:us-west-2:123456789012:stateMachine:HelloWorld-StateMachine"

	// Input Payload JSON
	input := `{"name": "John"}`

	// Start Execution
	startExecutionInput := &sfn.StartExecutionInput{
		StateMachineArn: aws.String(stateMachinearn),
		Input:           aws.String(input),
	}

	result, err := sfnClient.StartExecution(startExecutionInput)
	if err != nil {
		log.Println("Error starting execution: ", err)
		return
	}
	log.Println("Execution Started: ", result)
}
