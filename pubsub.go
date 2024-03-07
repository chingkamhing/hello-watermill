package main

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/spf13/cobra"
)

func runPubSub(cmd *cobra.Command, args []string) {
	// create pubsub channel
	logger := watermill.NewStdLogger(false, false)
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		logger,
	)
	defer pubSub.Close()
	// prepare topic handlers
	for topic, handler := range topicHandlers {
		messagesEmail, err := pubSub.Subscribe(context.Background(), topic)
		if err != nil {
			log.Fatalf("Subscribe topic %q error: %v", topic, err)
		}
		go func(messages <-chan *message.Message, handler MsgHandler) {
			for msg := range messages {
				handler(msg)
				// we need to Acknowledge that we received and processed the message; otherwise, it will be resent over and over again.
				msg.Ack()
			}
		}(messagesEmail, handler)
	}
	// publish message to the topics
	publishMessages(pubSub)
}
