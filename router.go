package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/spf13/cobra"
)

func runRouter(cmd *cobra.Command, args []string) {
	// create pubsub channel
	logger := watermill.NewStdLogger(false, false)
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
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	// prepare topic handlers
	for topic, handler := range topicHandlers {
		router.AddNoPublisherHandler(
			fmt.Sprintf("%v_handler", topic),
			topic,
			pubSub,
			message.NoPublishHandlerFunc(handler),
		)
	}
	// publish message to the topics
	go publishMessages(pubSub)
	// Now that all handlers are registered, we're running the Router.
	// Run is blocking while the router is running.
	ctx := context.Background()
	if err := router.Run(ctx); err != nil {
		log.Fatalf("router run error: %v", err)
	}
}
