package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type MsgHandler func(msg *message.Message) error
type TopicHandler map[string]MsgHandler

var topicHandlers = TopicHandler{
	"wk.email.send": handlerEmail,
	"wk.imos.post":  handlerImos,
}

var gochannelConfig = gochannel.Config{
	OutputChannelBuffer:            4,
	Persistent:                     false,
	BlockPublishUntilSubscriberAck: false,
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
		log.Printf("Publish to %s topic %s", topic, msg.UUID)

		time.Sleep(10 * time.Millisecond)
	}
}

func handlerEmail(msg *message.Message) error {
	log.Printf("Send email: %s message: %s", msg.UUID, string(msg.Payload))
	time.Sleep(5 * time.Second)
	return nil
}

func handlerImos(msg *message.Message) error {
	log.Printf("Post to imos: %s message: %s", msg.UUID, string(msg.Payload))
	time.Sleep(1 * time.Second)
	return nil
}
