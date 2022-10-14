package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/whitedevil31/atlan-backend/consumer-api-2/events"
	"github.com/whitedevil31/atlan-backend/consumer-api-2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Answers struct {
	QuestionId primitive.ObjectID `bson:"questionId" json:"questionId" `
	AnswerId   primitive.ObjectID `bson:"answerId" json:"answerId" `
	AnswerText string             `bson:"answerText" json:"answerText" `
}
type SubmitResponse struct {
	FormId    primitive.ObjectID `bson:"formId" json:"formId" `
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt" `
	Answers   []Answers          `bson:"answers" json:"answers" `
}
type FunctionStruct struct{}

func main() {

	conf := kafka.ReaderConfig{
		GroupTopics: []string{os.Getenv("TOPIC")},
		Brokers:     []string{os.Getenv("BROKER")},
		GroupID:     os.Getenv("GROUP_ID"),
		MaxBytes:    100,
	}
	reader := kafka.NewReader(conf)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.ErrorLogger.Println("FAILED TO INITIATE KAFKA CONSUMER")
		}
		logger.InfoLogger.Println("CONSUMED " + string(m.Key) + " EVENT")
		data := &SubmitResponse{}
		json.Unmarshal([]byte(m.Value), data)
		events := new(events.FunctionStruct)

		functionCall := reflect.ValueOf(events).MethodByName(string(m.Key))
		functionCall.Call([]reflect.Value{reflect.ValueOf(data)})

	}
}
