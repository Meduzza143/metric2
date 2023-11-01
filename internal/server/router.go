package server

import (
	handlers "github.com/Meduzza143/metric/internal/server/handlers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	//r.Use(handlers.LogMiddleware) TODO: раскурить таки
	r.HandleFunc(`/update/{type}/{name}/{value}`, handlers.UnpackMiddleware(handlers.LogMiddleware(handlers.UpdateHandle))).Methods("POST")
	r.HandleFunc(`/update/`, handlers.UnpackMiddleware(handlers.LogMiddleware(handlers.UpdateHandle))).Methods("POST") //for json handling
	r.HandleFunc(`/value/{type}/{name}`, handlers.UnpackMiddleware(handlers.LogMiddleware(handlers.GetMetric))).Methods("GET")
	r.HandleFunc(`/`, handlers.UnpackMiddleware(handlers.LogMiddleware(handlers.GetAll))).Methods("GET")

	return r
}
