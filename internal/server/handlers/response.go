package handlers

import (
	"net/http"

	"github.com/Meduzza143/metric/internal/logger"
)

func ResponseWritter(w http.ResponseWriter, status int, body string) {

	l := logger.GetLogger()
	l.Info().Int("status", status).Int("body size", len(body)).Msg("response")

	w.WriteHeader(status)
	w.Write([]byte(body))
}
