package main

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func main() {

	conf := kafka.ReaderConfig{
		GroupTopics: []string{"topic5", "topic6"},
		Brokers:     []string{"kafka:9092"},
		GroupID:     "group1",
		MaxBytes:    100,
	}
	reader := kafka.NewReader(conf)
	for {
		fmt.Println("WAITING.....")
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		if m.Topic == "topic5" {
			fmt.Println("TOPIC 5 EVENT TRIGGERED")
		} else if m.Topic == "topic6" {
			fmt.Println("TOPIC 6 EVENT TRIGGERED@")
		}

	}
}
