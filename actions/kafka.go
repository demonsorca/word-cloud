package actions

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gobuffalo/buffalo"
	kafka "github.com/segmentio/kafka-go"
)

func KafkaProduceHandler(c buffalo.Context) error {
	topic := "word-cloud-topic"
	textData := "Your input text data here."

	err := produceToKafka(topic, textData)
	if err != nil {
		log.Fatal("Error producing to Kafka: ", err)
	}

	consumeFromKafka(topic)
	return err
}

func produceToKafka(topic string, textData string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka-broker-url:port"},
		Topic:   topic,
	})

	words := strings.Fields(textData) // Split text into words

	for _, word := range words {
		message := kafka.Message{
			Key:   []byte(fmt.Sprintf("word-%s", word)),
			Value: []byte(word),
		}
		err := writer.WriteMessages(context.Background(), message)
		if err != nil {
			return err
		}
	}

	return nil
}

func consumeFromKafka(topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka-broker-url:port"},
		Topic:   topic,
		GroupID: "word-cloud-group",
	})

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Error reading message: ", err)
			break
		}
		word := string(message.Value)
		// Process the word for word cloud generation
		fmt.Println("Received word:", word)
	}
}
