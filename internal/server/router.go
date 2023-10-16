package server

import (
	handlers "github.com/Meduzza143/metric/internal/server/handlers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(`/update/{type}/{name}/{value}`, handlers.UpdateHandle).Methods("POST")
	r.HandleFunc(`/value/{type}/{name}`, handlers.GetMetric).Methods("GET")
	r.HandleFunc(`/`, handlers.GetAll).Methods("GET")
	return r
}
