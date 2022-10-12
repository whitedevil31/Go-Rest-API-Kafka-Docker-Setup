package routes

import (
	"github.com/gorilla/mux"
	producerController "github.com/whitedevil31/atlan-backend/producer-api/controller"
)

var RegisterStudentRoutes = func(router *mux.Router) {
	mainRouter := router.PathPrefix("/api").Subrouter()

	mainRouter.HandleFunc("/add_form", producerController.AddForm).Methods("POST")
	mainRouter.HandleFunc("/submit_form", producerController.AddResponse).Methods("POST")
	mainRouter.HandleFunc("/add_event", producerController.AddEvent).Methods("POST")
	//mainRouter.HandleFunc("/get_events/{formId}", producerController.GetEvents).Methods("GET")
}
