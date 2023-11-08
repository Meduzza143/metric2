package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/Meduzza143/metric/internal/logger"
)

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		l := logger.GetLogger()
		l.Info().Str("URI", req.URL.Path).Str("Method", req.Method).Str("Remote address", req.RemoteAddr).Msg("request")
		reqStart := time.Now()

		for i, v := range req.Header {
			l.Debug().Strs(i, v).Msg("server request header")
		}

		// Read the Body content
		var buf []byte
		if req.Body != nil {
			buf, _ = io.ReadAll(req.Body)
		}
		// Restore the io.ReadCloser to its original state
		req.Body = io.NopCloser(bytes.NewBuffer(buf))

		l.Debug().Str("BODY", string(buf)).Msg("request body")

		respdata := responseData{
			status: 0,
			size:   0,
		}
		loggingWriter := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   &respdata,
		}

		//next(w, req)
		next(&loggingWriter, req)

		l.Info().Int("status", respdata.status).Int("size", respdata.size).Msg("response")

		reqDuration := time.Now().Sub(reqStart)
		l.Info().Dur("request running time", reqDuration).Msg("request")
	})
}
