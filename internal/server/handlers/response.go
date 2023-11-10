package handlers

import (
	"net/http"

	"github.com/Meduzza143/metric/internal/logger"
	"github.com/Meduzza143/metric/internal/zipper"
)

type (
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func ResponseWritter(w http.ResponseWriter, status int, data []byte, settings RespSettings) {

	l := logger.GetLogger()
	l.Info().Int("status", status).Int("body size", len(data)).Msg("response")

	var answer []byte
	if settings.acceptEncoding == "gzip" {
		answer = zipper.GzipBytes(data)
	} else {
		answer = data
	}

	//l.Debug().Str("answer body", string(answer)).Msg("response")

	addHeaders(w)

	// for i, v := range w.Header() {
	// 	l.Debug().Strs(i, v).Msg("server response header")
	// }

	w.WriteHeader(status)
	w.Write(answer)

}

func addHeaders(w http.ResponseWriter) {
	w.Header().Set("This server is ", "MINE")

	switch respSet.acceptFormat {
	case "json":
		{
			w.Header().Set("Content-Type", "application/json")
		}
	default:
		{
			w.Header().Set("Content-Type", "text/html")
		}
	}

	switch respSet.acceptEncoding {
	case "gzip":
		{
			w.Header().Set("Content-Encoding", "gzip")
		}
	default:
		{
			w.Header().Set("Content-Encoding", "identity")
		}
	}
}
