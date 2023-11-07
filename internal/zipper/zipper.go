package zipper

import (
	"bytes"
	"compress/gzip"
)

func GzipBytes(data []byte) []byte {
	var buf bytes.Buffer
	gzwriter := gzip.NewWriter(&buf)
	gzwriter.Write(data)
	//gzwriter.Flush()
	//defer gzwriter.Close()
	gzwriter.Close()
	return buf.Bytes()
}

func UnGzipBytes(data []byte) []byte {
	r, _ := gzip.NewReader(bytes.NewReader(data))
	var b bytes.Buffer
	b.ReadFrom(r)
	r.Close()
	//defer r.Close()
	return b.Bytes()
}

//TODO: implement

// type WrappedWriter struct {
// 	rw http.ResponseWriter
// 	gw *gzip.Writer
// }

// func NewWriter(rw http.ResponseWriter) *WrappedWriter {
// 	gz := gzip.NewWriter(rw)
// 	return &WrappedWriter{rw, gz}
// }
// func (wr *WrappedWriter) Header() http.Header {
// 	return wr.rw.Header()
// }
// func (wr *WrappedWriter) Write(data []byte) (int, error) {
// 	return wr.gw.Write(data)
// }
// func (wr *WrappedWriter) WriteHeader(status int) {
// 	wr.rw.WriteHeader(status)
// }
// func (wr *WrappedWriter) Flush() {
// 	wr.gw.Flush()
// 	wr.gw.Close()
// }
