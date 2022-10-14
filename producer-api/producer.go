package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/whitedevil31/atlan-backend/producer-api/logger"
	"github.com/whitedevil31/atlan-backend/producer-api/routes"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterStudentRoutes(r)
	http.Handle("/", r)
	host := ":8080"
	logger.InfoLogger.Println("APP RUNNING ON " + host)
	http.ListenAndServe(host, r)

}
