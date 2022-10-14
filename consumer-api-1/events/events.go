package events

import (
	"context"
	"encoding/json"

	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/whitedevil31/atlan-backend/consumer-api-1/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
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

func (m FunctionStruct) SEND_EMAIL(data interface{}) {

	mg := mailgun.NewMailgun(os.Getenv("MAILER_DOMAIN"), os.Getenv("MAILER_KEY"))
	sender := os.Getenv("MAILER_ID")
	subject := "YOUR FORM RESPONSE HAS BEEN RECORDED"
	body := "We have received your form submission and here are your details. `Details of the form in a tabular format`"
	recipient := "np6771@srmist.edu.in"

	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, id, err := mg.Send(ctx, message)

	if err != nil {
		logger.ErrorLogger.Println("SEND EMAIL FUNCTION FAILED !")
	}

	logger.InfoLogger.Printf("SUBMISSION RESULT EMAILED ! | ID = %s", id)

}
func (m FunctionStruct) ADD_DATA_TO_SHEET(data interface{}) {
	bytesData, _ := json.Marshal(data)
	eventInfo := &SubmitResponse{}
	json.Unmarshal(bytesData, eventInfo)
	var vr sheets.ValueRange
	myval := []interface{}{}
	for i := 0; i < len(eventInfo.Answers); i++ {
		myval = append(myval, eventInfo.Answers[i].AnswerText)
	}
	const (
		secret = "./final.json"
	)

	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(secret), option.WithScopes(sheets.SpreadsheetsScope))

	if err != nil {
		logger.ErrorLogger.Println("SHEET SERVICE CANNOT BE INITIATED!")
	}

	vr.Values = append(vr.Values, myval)
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	readRange := os.Getenv("SPREADSHEET_RANGE")
	sheetData, errT := srv.Spreadsheets.Values.Append(spreadsheetId, readRange, &vr).ValueInputOption("RAW").Do()
	if errT != nil {
		logger.ErrorLogger.Println("FAILED TO ADD DATA TO SHEETS")
	}

	logger.InfoLogger.Println("ADDED DATA TO ROW " + sheetData.Updates.UpdatedRange)

}
