package producerController

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"github.com/whitedevil31/atlan-backend/producer-api/config"
	"github.com/whitedevil31/atlan-backend/producer-api/logger"
	"github.com/whitedevil31/atlan-backend/producer-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var wg sync.WaitGroup

func AddForm(w http.ResponseWriter, r *http.Request) {

	addForm := &utils.Form{}
	utils.ParseBody(r, addForm)
	c := config.GetDB()
	collection := c.Database("atlan-backend").Collection("forms")
	formId := primitive.NewObjectID()
	createdBy := primitive.NewObjectID()
	for i := 0; i < len(addForm.Questions); i++ {
		addForm.Questions[i].QuestionId = primitive.NewObjectID()
	}
	_, AddFormsError := collection.InsertOne(context.Background(), &utils.Form{
		FormTitle:       addForm.FormTitle,
		FormDescription: addForm.FormDescription,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedBy:       createdBy,
		Questions:       addForm.Questions,
	})
	if AddFormsError != nil {
		logger.ErrorLogger.Println("ERROR IN ADDING FORM")
		utils.JSONError(w, "SOMETHING_WENT_WRONG", utils.GetCode("SOMETHING_WENT_WRONG"))
		return

	}
	result := utils.AddFormResponse{}
	result.Message = "Form added successfullly!"
	logger.InfoLogger.Println("FORM ADDED " + "FORM ID " + result.FormId.Hex())
	result.FormId = formId
	res, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func AddResponse(w http.ResponseWriter, r *http.Request) {

	c := config.GetDB()
	addResponse := &utils.SubmitResponse{}
	utils.ParseBody(r, addResponse)
	ch := make(chan utils.GetEventsResponse)
	go GetEventData(addResponse.FormId, ch)
	wg.Add(1)
	getEventResult := <-ch

	collection := c.Database("atlan-backend").Collection("submissions")

	submissionId := primitive.NewObjectID()
	for i := 0; i < len(addResponse.Answers); i++ {
		addResponse.Answers[i].AnswerId = primitive.NewObjectID()
	}
	insertData := &utils.SubmitResponse{
		FormId:    addResponse.FormId,
		CreatedAt: time.Now(),
		Answers:   addResponse.Answers}

	_, AddSubmissionError := collection.InsertOne(context.Background(), insertData)

	if AddSubmissionError != nil {
		logger.ErrorLogger.Println("ERROR IN ADDING RESPONSE")
		utils.JSONError(w, "SOMETHING_WENT_WRONG", utils.GetCode("SOMETHING_WENT_WRONG"))
		return

	}
	wg.Wait()
	result := utils.AddSubmissionResponse{}
	eventData := getEventResult.Events
	formDataBytes, _ := json.Marshal(insertData)
	kafkaWriter := &kafka.Writer{
		Addr:  kafka.TCP(os.Getenv("KAFKA_BROKER")),
		Topic: os.Getenv("TOPIC"),
	}
	for _, item := range eventData {
		eventName, ok := item["eventName"].(string)

		if !ok {
			logger.ErrorLogger.Println("ERROR IN SENDING EVENTS TO KAFKA")
		}
		kafkaError := kafkaWriter.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(eventName),
			Value: formDataBytes,
		})
		if kafkaError != nil {
			logger.ErrorLogger.Println(kafkaError)
		} else {
			logger.InfoLogger.Println("MESSAGE SENT FROM PRODUCER FOR RESPONSE " + result.SubmissionId.Hex())
		}
	}
	logger.InfoLogger.Println("RESPONSE ADDED " + "RESPONSE ID" + result.SubmissionId.Hex())
	result.Message = "Response added successfullly!"
	result.SubmissionId = submissionId
	res, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
func AddEvent(w http.ResponseWriter, r *http.Request) {

	addEvent := &utils.Event{}
	utils.ParseBody(r, addEvent)
	c := config.GetDB()
	collection := c.Database("atlan-backend").Collection("events")

	_, addEventError := collection.InsertOne(context.Background(), &utils.Event{
		FormId:    addEvent.FormId,
		EventName: addEvent.EventName,
	})
	if addEventError != nil {
		logger.ErrorLogger.Println("ERROR IN ADDING EVENT TO FORM" + addEvent.FormId.Hex())
		utils.JSONError(w, "SOMETHING_WENT_WRONG", utils.GetCode("SOMETHING_WENT_WRONG"))
		return

	}
	result := utils.AddEventResponse{}
	result.Message = "Event added successfullly!"
	logger.InfoLogger.Println("EVENT ADDED TO FORM" + addEvent.FormId.Hex())
	res, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
func GetEventData(formId primitive.ObjectID, ch chan utils.GetEventsResponse) (utils.GetEventsResponse, error) {
	c := config.GetDB()
	var events []primitive.M
	result := utils.GetEventsResponse{}
	collection := c.Database("atlan-backend").Collection("events")
	cursor, getEventError := collection.Find(context.Background(), bson.D{{Key: "formId", Value: formId}})
	if getEventError != nil {
		logger.ErrorLogger.Println("ERROR IN FETCHING EVENTS FOR FORM ID " + formId.Hex())
		return result, errors.New("SOMETHING_WENT_WRONG")

	} else {
		getEventError = cursor.All(context.Background(), &events)

		if getEventError != nil {
			logger.ErrorLogger.Println("SOMETHING WENT WRONG")
			return result, errors.New("SOMETHING_WENT_WRONG")

		}

	}
	logger.InfoLogger.Println("EVENTS FETCHED FOR FORM ID " + formId.Hex())

	result.Events = events

	ch <- result
	wg.Done()
	return result, nil
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c := config.GetDB()
	formId, err := primitive.ObjectIDFromHex(vars["formId"])
	if err != nil {
		logger.ErrorLogger.Println("INVALID FORM ID " + formId.Hex())
		utils.JSONError(w, "INVALID_ID", utils.GetCode("INVALID_ID"))
		return
	}
	var events []primitive.M
	result := utils.GetEventsResponse{}
	collection := c.Database("atlan-backend").Collection("events")
	cursor, getEventError := collection.Find(context.Background(), bson.D{{Key: "formId", Value: formId}})
	if getEventError != nil {
		logger.ErrorLogger.Println("ERROR IN FETCHING EVENTS FOR FORM ID " + formId.Hex())
		utils.JSONError(w, "SOMETHING_WENT_WRONG", utils.GetCode("SOMETHING_WENT_WRONG"))
		return

	} else {
		getEventError = cursor.All(context.Background(), &events)

		if getEventError != nil {
			logger.ErrorLogger.Println("SOMETHING WENT WRONG")
			utils.JSONError(w, "SOMETHING_WENT_WRONG", utils.GetCode("SOMETHING_WENT_WRONG"))
			return

		}

	}
	logger.InfoLogger.Println("EVENTS FETCHED FOR FORM ID " + formId.Hex())

	result.Events = events
	res, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
