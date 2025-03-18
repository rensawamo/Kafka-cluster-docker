package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		messageData := map[string]interface{}{
			"topic":     msg.Topic,
			"partition": msg.Partition,
			"offset":    msg.Offset,
			"key":       string(msg.Key),
			"value":     string(msg.Value),
			"timestamp": msg.Timestamp,
		}

		jsonData, err := json.MarshalIndent(messageData, "", "  ")
		if err != nil {
			log.Printf("Error marshaling message data: %v", err)
			continue
		}

		log.Printf("Received message:\n%s\n", jsonData)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	error := godotenv.Load()
	if error != nil {
		log.Fatalf("Error loading .env file")
	}

	broker1 := os.Getenv("BROKER1")
	broker2 := os.Getenv("BROKER2")
	broker3 := os.Getenv("BROKER3")

	brokers := []string{
		fmt.Sprintf("%s:9093", broker1),
		fmt.Sprintf("%s:9094", broker2),
		fmt.Sprintf("%s:9095", broker3),
	}
	fmt.Println(brokers)

	groupID := "consumer-group"

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0 // specify appropriate Kafka version
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	ctx := context.Background()

	for {
		err := consumerGroup.Consume(ctx, []string{"post-likes"}, exampleConsumerGroupHandler{})
		if err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
	}
}
