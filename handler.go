package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MsgHandler func(msg *message.Message) error
type TopicHandler map[string]MsgHandler

var topicHandlers = TopicHandler{
	"wk.email.send": handlerEmail,
	"wk.imos.post":  handlerImos,
}

func (handler TopicHandler) Topics() []string {
	topics := make([]string, 0, len(topicHandlers))
	for topic := range topicHandlers {
		topics = append(topics, topic)
	}
	return topics
}

func publishMessages(publisher message.Publisher) {
	topics := topicHandlers.Topics()
	numTopics := len(topics)
	for {
		msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, Wah Kwong!"))
		r := rand.Intn(numTopics)
		topic := topics[r]
		if err := publisher.Publish(topic, msg); err != nil {
			log.Fatalf("Publish error: %v", err)
		}

		time.Sleep(time.Second)
	}
}

func handlerEmail(msg *message.Message) error {
	fmt.Printf("Send email: %s message: %s\n", msg.UUID, string(msg.Payload))
	return nil
}

func handlerImos(msg *message.Message) error {
	fmt.Printf("Post to imos: %s message: %s\n", msg.UUID, string(msg.Payload))
	return nil
}
