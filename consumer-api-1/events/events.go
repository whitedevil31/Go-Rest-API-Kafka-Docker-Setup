package events

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitResponse struct {
	FormId    primitive.ObjectID `bson:"formId" json:"formId" `
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt" `
	Answers   []Answers          `bson:"answers" json:"answers" `
}
type Answers struct {
	QuestionId primitive.ObjectID `bson:"questionId" json:"questionId" `
	AnswerId   primitive.ObjectID `bson:"answerId" json:"answerId" `
	AnswerText string             `bson:"answerText" json:"answerText" `
}
type FunctionStruct struct{}

func (m FunctionStruct) EVENT_NAME_1(data interface{}) {

	bytesData, _ := json.Marshal(data)
	eventInfo := &SubmitResponse{}
	json.Unmarshal(bytesData, eventInfo)
	fmt.Println("FUNC 1 CONSUMER 1")
	fmt.Println(eventInfo.FormId)
}
func (m FunctionStruct) EVENT_NAME_2(data interface{}) {

	bytesData, _ := json.Marshal(data)
	eventInfo := &SubmitResponse{}
	json.Unmarshal(bytesData, eventInfo)
	fmt.Println("FUNC 2 CONSUMER 1")
	fmt.Println(eventInfo.FormId)

}
func (m FunctionStruct) EVENT_NAME_3(data interface{}) {

	bytesData, _ := json.Marshal(data)
	eventInfo := &SubmitResponse{}
	json.Unmarshal(bytesData, eventInfo)
	fmt.Println("FUNC 3 CONSUMER 1")
	fmt.Println(eventInfo.FormId)

}
