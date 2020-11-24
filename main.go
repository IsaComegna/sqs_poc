package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func GetQueues(svc *sqs.SQS) (*sqs.ListQueuesOutput, error) {
	// Create an SQS service client
	// snippet-start:[sqs.go.list_queues.call]

	result, err := svc.ListQueues(nil)
	// snippet-end:[sqs.go.list_queues.call]
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetQueueURL(svc *sqs.SQS, queue *string) (*sqs.GetQueueUrlOutput, error) {
	// Create an SQS service client
	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SendMsg(svc *sqs.SQS, queueURL *string) error {
	// Create an SQS service client
	// snippet-start:[sqs.go.send_message.call]
	_, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("The Whistler"),
			},
			"Author": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("John Grisham"),
			},
			"WeeksOn": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
		QueueUrl:    queueURL,
	})
	// snippet-end:[sqs.go.send_message.call]
	if err != nil {
		return err
	}

	return nil
}

func GetMessages(svc *sqs.SQS, queueURL *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
	// Create an SQS service client
	// snippet-start:[sqs.go.receive_messages.call]
	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   timeout,
	})
	// snippet-end:[sqs.go.receive_messages.call]
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	svc := sqs.New(sess, aws.NewConfig().WithEndpoint("http://localhost:9324"))

	result, err := GetQueues(svc)
	if err != nil {
		fmt.Println("Got an error retrieving queue URLs:")
		fmt.Println(err)
		return
	}

	// snippet-start:[sqs.go.list_queues.display]
	for i, url := range result.QueueUrls {
		fmt.Printf("%d: %s\n", i, *url)
	}
	// snippet-end:[sqs.go.list_queues.display]

	queueURL := "http://localhost:9324/queue/default"

	err = SendMsg(svc, &queueURL)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return
	}

	fmt.Println("Sent message to queue ")

	var timeout int64 = 3

	msgResult, err := GetMessages(svc, &queueURL, &timeout)
	if err != nil {
		fmt.Println("Got an error receiving messages:")
		fmt.Println(err)
		return
	}

	fmt.Println("Message ID:     " + *msgResult.Messages[0].MessageId)

	// snippet-start:[sqs.go.receive_messages.print_handle]
	fmt.Println("Message Handle: " + *msgResult.Messages[0].ReceiptHandle)
}
