package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type MsgHandler func(msg *message.Message) error
type TopicHandler map[string]MsgHandler

const toEmail = "chingkamhing@gmail.com"

var topicHandlers = TopicHandler{
	"test-one": handlerOne,
	"test-two": handlerTwo,
}

var gochannelConfig = gochannel.Config{
	OutputChannelBuffer:            20,
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
	// sleep a short while to make sure router is running
	time.Sleep(100 * time.Millisecond)
	for i := 0; i < numMessage; i++ {
		msg := message.NewMessage(watermill.NewUUID(), []byte(fmt.Sprintf("Hello, Wah Kwong! (%v)", i)))
		topic := "test-one"
		if err := publisher.Publish(topic, msg); err != nil {
			log.Fatalf("Publish error: %v", err)
		}
		log.Printf("Publish %v to %s topic %s", i, topic, msg.UUID)

		time.Sleep(200 * time.Millisecond)
	}
}

func handlerOne(msg *message.Message) error {
	body := string(msg.Payload)
	err := sendEmail(toEmail, body)
	log.Printf("Handler one: %s message: %s", msg.UUID, body)
	return err
}

func handlerTwo(msg *message.Message) error {
	xml := string(msg.Payload)
	err := postImos(xml)
	log.Printf("Handler two: %s message: %s", msg.UUID, xml)
	return err
}

func sendEmail(to, body string) error {
	_ = to
	_ = body
	time.Sleep(5 * time.Second)
	return nil
}

func postImos(xml string) error {
	_ = xml
	time.Sleep(1 * time.Second)
	return nil
}
