package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/whitedevil31/atlan-backend/consumer-api-1/events"
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
		GroupTopics: []string{"POST_FORM_SUBMIT"},
		Brokers:     []string{"kafka:9092"},
		GroupID:     "GROUP1",
		MaxBytes:    100,
	}
	fmt.Println(os.Getenv("MONGO"))
	reader := kafka.NewReader(conf)
	for {

		fmt.Println("WAITING CONSUMER 1.....")
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Print(err.Error())
			continue
		}

		data := &SubmitResponse{}
		json.Unmarshal([]byte(m.Value), data)
		events := new(events.FunctionStruct)

		functionCall := reflect.ValueOf(events).MethodByName(string(m.Key))
		functionCall.Call([]reflect.Value{reflect.ValueOf(data)})

	}
}
