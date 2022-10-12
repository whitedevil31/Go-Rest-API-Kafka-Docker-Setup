package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var c *mongo.Client
var collection *mongo.Collection

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

type ErrorFormat struct {
	Message string `json:"message" bson:"message"`
}

func JSONError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	t := ErrorFormat{Message: err}
	json.NewEncoder(w).Encode(t)
}
func GetCode(err string) int {
	switch err {
	case "EMAIL_IN_USE":
		return http.StatusConflict
	case "VALIDATION_FAILED":
		return http.StatusBadRequest
	case "INTERNAL_SERVER_ERROR":
		return http.StatusInternalServerError
	case "SOMETHING_WENT_WRONG":
		return http.StatusInternalServerError
	case "USER_NOT_EXIST":
		return http.StatusNotFound
	case "INCORRECT_PASSWORD":
		return http.StatusNotFound
	case "UNAUTHORISED_ACCESS":
		return http.StatusUnauthorized
	case "INVALID_LOGIN":
		return http.StatusUnauthorized
	case "TOKEN_ERROR":
		return http.StatusFailedDependency
	case "INVALID_ID":
		return http.StatusBadRequest
	case "RESULT_NOT_FOUND":
		return http.StatusNotFound
	case "STUDENT_NOT_FOUND":
		return http.StatusNotFound
	case "TOKEN_NOT_FOUND":
		return http.StatusNotAcceptable
	default:
		return 500
	}

}

type Questions struct {
	QuestionId          primitive.ObjectID `bson:"questionId" json:"questionId" `
	QuestionTitle       string             `bson:"questionTitle" json:"questionTitle" `
	QuestionDescription string             `bson:"questionDescription" json:"questionDescription" `
	AnswerType          string             `bson:"answerType" json:"answerType" `
}
type Form struct {
	FormTitle       string             `bson:"formTitle" json:"formTitle"  `
	FormDescription string             `bson:"formDescription" json:"formDescription"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt" `
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt" `
	CreatedBy       primitive.ObjectID `bson:"createdBy" json:"createdBy" `
	Questions       []Questions        `bson:"questions" json:"questions" `
}
type AddFormResponse struct {
	Message string             `bson:"message" json:"message" `
	FormId  primitive.ObjectID `bson:"formId" json:"formId" `
}
type AddSubmissionResponse struct {
	Message      string             `bson:"message" json:"message" `
	SubmissionId primitive.ObjectID `bson:"submissionId" json:"submissionId" `
}
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
type Event struct {
	FormId    primitive.ObjectID `bson:"formId" json:"formId" `
	EventName string             `bson:"eventName" json:"eventName" `
}
type AddEventResponse struct {
	Message string `bson:"message" json:"message" `
}
type GetEventsResponse struct {
	Events []primitive.M `bson:"events"  json:"events"`
}
