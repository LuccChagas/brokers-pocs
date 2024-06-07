package main

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received message: %s\n", string(msg.Value))
		sess.MarkMessage(msg, "")
	}

	return nil
}

func main() {
	brokers := []string{"localhost:9092"}
	groupID := "consumer-group"

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Panicf("error to creating consumer group client %v", err)
	}

	ctx := context.Background()

	for {
		err := consumerGroup.Consume(ctx, []string{"post-likes"}, exampleConsumerGroupHandler{})
		if err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
	}
}
