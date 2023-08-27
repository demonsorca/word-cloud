package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	kafka "github.com/segmentio/kafka-go"
)

var FruitMap = make(map[string]int)
var GameMap = make(map[string]int)
var MovieMap = make(map[string]int)

func KafkaProduceHandler(c buffalo.Context) (err error) {
	for {
		var responseMap map[string]int
		x, err := io.ReadAll(c.Request().Body)
		newMap := make(map[string]interface{})
		err = json.Unmarshal(x, &newMap)
		if err != nil {
			log.Fatal("Error in unmarshaling request data: ", err)

		}
		if _, ok := newMap["message"]; !ok {
			err = errors.New("message not found")
			break
		}
		topic := "word-cloud-group"
		textData := newMap["message"].(string)
		fmt.Println("Request newMap", topic, textData)
		wordType := newMap["type"].(string)
		switch wordType {
		case "fruit":
			if _, ok := FruitMap[textData]; ok {
				FruitMap[textData] += 1
			} else {
				FruitMap[textData] = 1
			}
			responseMap = FruitMap
		case "game":
			if _, ok := GameMap[textData]; ok {
				GameMap[textData] += 1
			} else {
				GameMap[textData] = 1
			}
			responseMap = GameMap
		case "movie":
			if _, ok := MovieMap[textData]; ok {
				MovieMap[textData] += 1
			} else {
				MovieMap[textData] = 1
			}
			responseMap = MovieMap
		}
		// err = produceToKafka(topic, textData)
		if err != nil {
			log.Fatal("Error producing to Kafka: ", err)
			break
		}
		resBytes, err := json.Marshal(responseMap)
		if err != nil {
			log.Fatal("Error in Marshalling: ", err)
			break
		}
		err = c.Render(http.StatusOK, r.JSON(map[string]string{"message": "message produced successfully", "result": string(resBytes)}))
		// consumeFromKafka(topic)
		break
	}
	return
}

func produceToKafka(topic string, textData string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"http://localhost:9092"},
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
		Brokers: []string{"http://localhost:9092"},
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
