package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/IBM/sarama"
)

type Message struct {
	UserId     int    `json:"user_id"`
	PostId     string `json:"post_id"`
	UserAction string `json:"user_action"`
}

func main() {
	brokers := []string{"172.18.0.2:9093", "172.18.0.3:9094", "172.18.0.4:9095"}


	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
		os.Exit(1)
	}
	defer producer.Close()

	postId := [5]string{"POST00001", "POST00002", "POST00003", "POST00004", "POST00005"}
	userId := [5]int{100001, 100002, 100003, 100004, 100005}
	userAction := [5]string{"love", "like", "hate", "smile", "cry"}

	for {
		message := Message{
			UserId:     userId[rand.Intn(len(userId))],
			PostId:     postId[rand.Intn(len(postId))],
			UserAction: userAction[rand.Intn(len(userAction))],
		}

		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("Failed to marshal message: %v", err)
			os.Exit(1)
		}

		msg := &sarama.ProducerMessage{
			Topic: "post-likes",
			Key:   sarama.StringEncoder(message.PostId),
			Value: sarama.StringEncoder(jsonMessage),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Fatalf("Failed to send message: %v", err)
			os.Exit(1)
		}

		log.Printf("Message sent: %s (partition=%d, offset=%d)\n", string(jsonMessage), partition, offset)

		// 5秒待機
		time.Sleep(5 * time.Second)
	}
}
