package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

/*
   Сервер опционально принимать запросы в сжатом формате (при наличии соответствующего HTTP-заголовка Content-Encoding).
   Отдавать сжатый ответ клиенту, который поддерживает обработку сжатых ответов (с HTTP-заголовком Accept-Encoding).
*/

/*
	Делаем новую обертку для http.Writer, но вместо обычного врайтера , используем в нем gzipWriter
*/

func UnpackMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//l := logger.GetLogger()
		if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
			//l.Info().Str("Content-Encoding", "gzip").Msg("Content-Encoding")

			gzReader, err := gzip.NewReader(req.Body)
			if err != nil {
				http.Error(w, "Failed to decode gzip body", http.StatusBadRequest)
				return
			}
			defer gzReader.Close()
			req.Body = gzReader
		}

		next(w, req)
	})
}
