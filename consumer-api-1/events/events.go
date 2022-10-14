package events

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

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
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

}
func (m FunctionStruct) ADD_DATA_TO_SHEET(data interface{}) {
	const (
		secret = "./final.json"
	)

	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(secret), option.WithScopes(sheets.SpreadsheetsScope))

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}
	var vr sheets.ValueRange
	myval := []interface{}{"One", "Two", "Three"}
	vr.Values = append(vr.Values, myval)
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	readRange := os.Getenv("SPREADSHEET_RANGE")
	_, errT := srv.Spreadsheets.Values.Append(spreadsheetId, readRange, &vr).ValueInputOption("RAW").Do()
	if errT != nil {
		fmt.Println(errT)
	}
	log.Println("PUSHED")

}
