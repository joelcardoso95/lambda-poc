package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"log"
)

const (
	awsRegion = "us-east-1"
)

type User struct {
	Username string `json:"userName"`
	UserId   string `json:"userId"`
}

func main() {
	// Create a new context
	ctx := context.Background()
	// Start the Lambda handler
	log.Println("Starting Lambda Handler")
	lambda.StartWithOptions(handler, lambda.WithContext(ctx))
}

func handler(ctx context.Context, events events.DynamoDBEvent) {
	log.Println("Received events: ", events.Records)

	// Create a new session
	sess := session.Must(session.NewSession())
	//
	// Create SFN New Client
	sfnClient := sfn.New(sess, aws.NewConfig().WithRegion(awsRegion))
	//
	// SFN Arn
	stateMachinearn := "arn:aws:states:us-east-1:#yourAccount:stateMachine:poc"
	//
	for _, record := range events.Records {
		log.Println("Processing Record: ", record.EventName)

		if record.EventName == "MODIFY" {
			// reading the new image and setting user properties
			newImage := record.Change.NewImage
			user := User{
				Username: newImage["user_name"].String(),
				UserId:   newImage["user_id"].String(),
			}
			inputJson, err := json.Marshal(user)
			if err != nil {
				log.Println("Error marshalling NewImage to JSON: ", err)
				panic(err)
			}
			input := string(inputJson)
			log.Println("Input: ", input)

			// Start Execution
			startExecutionInput := &sfn.StartExecutionInput{
				StateMachineArn: aws.String(stateMachinearn),
				Input:           aws.String(input),
			}

			result, err := sfnClient.StartExecution(startExecutionInput)
			if err != nil {
				log.Println("Error starting execution: ", err)
				panic(err)
			}
			log.Println("Execution Started: ", result)
		}
	}
}

// Example of sending a task success
//func sendTaskSuccess(sfnClient *sfn.SFN, taskToken string, result *sfn.StartExecutionOutput) {
//	resultJson, err := json.Marshal(result)
//	if err != nil {
//		log.Println("Error marshalling result to JSON: ", err)
//		panic(err)
//	}
//
//	sendTaskSuccessInput := &sfn.SendTaskSuccessInput{
//		TaskToken: aws.String(taskToken),
//		Output:    aws.String(string(resultJson)),
//	}
//
//	_, err = sfnClient.SendTaskSuccess(sendTaskSuccessInput)
//	if err != nil {
//		log.Println("Error sending task success: ", err)
//		panic(err)
//	}
//
//	log.Println("Task Success Sent")
//}
