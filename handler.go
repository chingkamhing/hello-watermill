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

const toEmail = "chingkamhing@gmail.com"

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

		time.Sleep(1 * time.Second)
	}
}

func handlerEmail(msg *message.Message) error {
	body := string(msg.Payload)
	log.Printf("Send email: %s message: %s", msg.UUID, body)
	emailWorker.Submit(func() {
		sendEmail(toEmail, body)
	})
	return nil
}

func handlerImos(msg *message.Message) error {
	xml := string(msg.Payload)
	log.Printf("Post to imos: %s message: %s", msg.UUID, xml)
	emailWorker.Submit(func() {
		postImos(xml)
	})
	return nil
}

func sendEmail(to, body string) error {
	_ = to
	_ = body
	log.Printf("Send email to: %s body: %s", to, body)
	time.Sleep(5 * time.Second)
	return nil
}

func postImos(xml string) error {
	_ = xml
	log.Printf("Post xml to IMOS: %s", xml)
	time.Sleep(2 * time.Second)
	return nil
}
