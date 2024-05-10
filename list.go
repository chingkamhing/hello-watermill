package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/spf13/cobra"
)

func runList(cmd *cobra.Command, args []string) {
	sdkConfig, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(awsSqsRegion),
		awsconfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     awsAccessKey,
				SecretAccessKey: awsAccessSecret,
			},
		}),
	)
	if err != nil {
		log.Fatalf("LoadDefaultConfig() error: %v", err)
	}
	fmt.Println("Let's list SQS for your account.")
	var queueUrls []string
	sqsClient := sqs.NewFromConfig(sdkConfig)
	paginator := sqs.NewListQueuesPaginator(sqsClient, &sqs.ListQueuesInput{})
	for paginator.HasMorePages() {
		output, e := paginator.NextPage(context.TODO())
		if e != nil {
			err = e
			break
		}
		queueUrls = append(queueUrls, output.QueueUrls...)
	}
	if err != nil {
		log.Fatalf("List queue error: %v", err)
	}
	if len(queueUrls) <= 0 {
		fmt.Println("You don't have any queues!")
		os.Exit(0)
	}
	for _, queueUrl := range queueUrls {
		fmt.Printf("\t%v\n", queueUrl)
	}
}
