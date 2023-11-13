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

func ResponseWritter(exWriter ExtendedWriter, data []byte) {

	l := logger.GetLogger()
	l.Info().Int("status", exWriter.status).Int("body size", len(data)).Msg("response")

	var answer []byte
	if exWriter.acceptEncoding == "gzip" {
		answer = zipper.GzipBytes(data)
	} else {
		answer = data
	}

	exWriter.addHeaders()
	exWriter.writeCurrentHeader()
	exWriter.Write(answer)
}

func (exWriter *ExtendedWriter) writeCurrentHeader() {
	exWriter.WriteHeader(exWriter.status)
}

func (exWriter *ExtendedWriter) addHeaders() {
	exWriter.Header().Set("This server is ", "MINE")

	switch exWriter.acceptFormat {
	case "json":
		{
			exWriter.Header().Set("Content-Type", "application/json")
		}
	default:
		{
			exWriter.Header().Set("Content-Type", "text/html")
		}
	}

	switch exWriter.acceptEncoding {
	case "gzip":
		{
			exWriter.Header().Set("Content-Encoding", "gzip")
		}
	default:
		{
			exWriter.Header().Set("Content-Encoding", "identity")
		}
	}
}
