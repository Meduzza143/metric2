package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/Meduzza143/metric/internal/logger"
)

/*
   Сервер опционально принимать запросы в сжатом формате (при наличии соответствующего HTTP-заголовка Content-Encoding).
   Отдавать сжатый ответ клиенту, который поддерживает обработку сжатых ответов (с HTTP-заголовком Accept-Encoding).
*/

type WrappedWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func newWriter(rw http.ResponseWriter) *WrappedWriter {
	gz := gzip.NewWriter(rw)
	return &WrappedWriter{rw, gz}
}
func (wr *WrappedWriter) Header() http.Header {
	return wr.rw.Header()
}
func (wr *WrappedWriter) Write(data []byte) (int, error) {
	return wr.gw.Write(data)
}
func (wr *WrappedWriter) WriteHeader(status int) {
	wr.rw.WriteHeader(status)
}
func (wr *WrappedWriter) Flush() {
	wr.gw.Flush()
	wr.gw.Close()
}

/*
	Делаем новую обертку для http.Writer, но вместо обычного врайтера , используем в нем gzipWriter
*/

func UnpackMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		l := logger.GetLogger()
		if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
			l.Info().Str("Content-Encoding", "gzip").Msg("Content-Encoding")
			//gzipWriter := new(WrappedWriter)
			gzipWriter := newWriter(w)

			next(gzipWriter, req)
			defer gzipWriter.Flush()
		} else {
			next(w, req)
		}

	})

}
