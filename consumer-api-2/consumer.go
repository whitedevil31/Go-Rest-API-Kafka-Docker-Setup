package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/segmentio/kafka-go"
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

func (m FunctionStruct) EVENT_NAME_1(s *SubmitResponse) {
	fmt.Println("FUNC 1 CONSUMER 2" + " " + string(s.Answers[0].AnswerText))
}

func (m FunctionStruct) EVENT_NAME_2(s *SubmitResponse) {
	fmt.Println("FUNC 2 CONSUMER 2" + string(s.Answers[0].AnswerText))
}
func (m FunctionStruct) EVENT_NAME_3(s *SubmitResponse) {
	fmt.Println("FUNC 3 CONSUMER 2" + string(s.Answers[0].AnswerText))
}
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

		fmt.Println("WAITING CONSUMER 2.....")
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Print(err.Error())
			continue
		}

		functionStruct := FunctionStruct{}
		eventName := string(m.Key)
		data := &SubmitResponse{}
		json.Unmarshal([]byte(m.Topic), data)
		functionCall := reflect.ValueOf(functionStruct).MethodByName(eventName)
		functionCall.Call([]reflect.Value{reflect.ValueOf(data)})

	}
}
