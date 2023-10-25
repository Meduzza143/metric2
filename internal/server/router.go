package server

import (
	handlers "github.com/Meduzza143/metric/internal/server/handlers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	//r.Use(handlers.LogMiddleware)

	//r.HandleFunc(`/update/{type}/{name}/{value}`, handlers.GetMetric(handlers.TestMiddleware(http.Request, http.Response))).Methods("POST")
	r.HandleFunc(`/update/{type}/{name}/{value}`, handlers.UpdateHandle).Methods("POST")
	//r.HandleFunc(`/update/{type}/{name}/{value}`, handlers.LogMiddleware(handlers.UpdateHandle)).Methods("POST")

	//r.HandleFunc(`/update/{type}/{name}/{value}`, handlers.UpdateHandle).Methods("POST")
	r.HandleFunc(`/value/{type}/{name}`, handlers.GetMetric).Methods("GET")
	//r.Handle()
	//r.HandleFunc(`/`, handlers.GetAll).Methods("GET")
	r.HandleFunc(`/`, handlers.LogMiddleware(handlers.GetAll)).Methods("GET")
	// r.HandleFunc(`/update/{type}/{name}/{value}`, handlers.UpdateHandle).Methods("POST")
	// r.HandleFunc(`/value/{type}/{name}`, handlers.GetMetric).Methods("GET")
	// r.HandleFunc(`/`, handlers.GetAll).Methods("GET")
	return r
}
