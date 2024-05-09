package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amazonsqs/connection"
	"github.com/ThreeDotsLabs/watermill-amazonsqs/sqs"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

var driverSupportWorkGroup = map[string]struct{}{
	"firestorepubsub": {},
	"gcppubsub":       {},
	"kafka":           {},
	"nats":            {},
	"rabbitmq":        {},
	"redis":           {},
	"sql":             {},
	"sqs":             {}, //FIXME, not sure yet
}

var awsAccessKey = os.Getenv("AWS_ACCESS_KEY")
var awsAccessSecret = os.Getenv("AWS_ACCESS_SECRET")
var awsSqsRegion = os.Getenv("AWS_SQS_REGION")
var awsSqsEndpoint = os.Getenv("AWS_SQS_ENDPOINT")

func runRouter(cmd *cobra.Command, args []string) {
	// create pubsub channel
	logger := watermill.NewStdLogger(debug, false)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		log.Fatalf("New router error: %v", err)
	}
	defer router.Close()
	// SignalsHandler will gracefully shutdown Router when SIGTERM is received.
	// You can also close the router by just calling `r.Close()`.
	router.AddPlugin(plugin.SignalsHandler)
	// For simplicity, we are using the gochannel Pub/Sub here,
	// You can replace it with any Pub/Sub implementation, it will work the same.
	publisher, subscriber, err := createPubsuber(logger)
	if err != nil {
		log.Fatalf("Create Publisher/Subscriber error: %v", err)
	}
	// only driver support workgroup can have multiple workers
	_, supportWorkGroup := driverSupportWorkGroup[pubsubDriver]
	if !supportWorkGroup {
		numWorkers = 1
	}
	// prepare topic handlers
	for topic, handler := range topicHandlers {
		for i := 0; i < numWorkers; i++ {
			router.AddNoPublisherHandler(
				fmt.Sprintf("%v_handler_%v", topic, i),
				topic,
				subscriber,
				message.NoPublishHandlerFunc(handler),
			)
		}
	}
	// publish message to the topics
	go publishMessages(publisher)
	// Now that all handlers are registered, we're running the Router.
	// Run is blocking while the router is running.
	ctx := context.Background()
	if err := router.Run(ctx); err != nil {
		log.Fatalf("router run error: %v", err)
	}
}

func createPubsuber(logger watermill.LoggerAdapter) (message.Publisher, message.Subscriber, error) {
	logger.Info("Create pubsub", watermill.LogFields{"driver": pubsubDriver})
	switch pubsubDriver {
	case "sqs":
		cfg, err := awsconfig.LoadDefaultConfig(
			context.Background(),
			connection.SetEndPoint(awsSqsEndpoint),
			awsconfig.WithRegion(awsSqsRegion),
			awsconfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     awsAccessKey,
					SecretAccessKey: awsAccessSecret,
				},
			}),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("LoadDefaultConfig(): %w", err)
		}
		publisher, err := sqs.NewPublisher(sqs.PublisherConfig{
			AWSConfig:              cfg,
			CreateQueueIfNotExists: true,
			Marshaler:              sqs.DefaultMarshalerUnmarshaler{},
		}, logger)
		if err != nil {
			return nil, nil, fmt.Errorf("NewPublisher(): %w", err)
		}
		subscriber, err := sqs.NewSubscriber(sqs.SubscriberConfig{
			AWSConfig:                    cfg,
			CreateQueueInitializerConfig: sqs.QueueConfigAtrributes{},
			Unmarshaler:                  sqs.DefaultMarshalerUnmarshaler{},
		}, logger)
		if err != nil {
			return nil, nil, fmt.Errorf("NewSubscriber(): %w", err)
		}
		return publisher, subscriber, nil
	case "redis":
		redisClient := redis.NewClient(&redis.Options{
			Addr: redisAddr,
			DB:   redisDb,
		})
		subscriber, err := redisstream.NewSubscriber(
			redisstream.SubscriberConfig{
				Client:        redisClient,
				Unmarshaller:  redisstream.DefaultMarshallerUnmarshaller{},
				ConsumerGroup: "wk_api_consumer_group",
			},
			logger,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("pubsub new redis subscriber: %w", err)
		}
		publisher, err := redisstream.NewPublisher(
			redisstream.PublisherConfig{
				Client:     redisClient,
				Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
			},
			logger,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("pubsub new redis publisher: %w", err)
		}
		return publisher, subscriber, nil
	case "gochannel":
		fallthrough
	default:
		// config watermill to use go channel as pubsub adapter
		pubSuber := gochannel.NewGoChannel(gochannelConfig, logger)
		return pubSuber, pubSuber, nil
	}
}
