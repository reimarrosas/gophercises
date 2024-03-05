package urlshortener

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	WarnLog  *log.Logger
	InfoLog  *log.Logger
	ErrorLog *log.Logger
)

type loggedResponseWriter struct {
	rw         http.ResponseWriter
	statusCode int
}

func (lrw *loggedResponseWriter) Header() http.Header {
	return lrw.rw.Header()
}

func (lrw *loggedResponseWriter) Write(buf []byte) (int, error) {
	return lrw.rw.Write(buf)
}

func (lrw *loggedResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.rw.WriteHeader(statusCode)
}

func InitLogger(filepath string) error {

	var w io.Writer = os.Stdin

	if filepath != "" {
		f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil
		}

		w = io.MultiWriter(w, f)
	}

	flags := log.Ldate | log.Ltime

	WarnLog = log.New(w, "WARN: ", flags)
	InfoLog = log.New(w, "INFO: ", flags)
	ErrorLog = log.New(w, "ERROR: ", flags)

	return nil
}

func LoggerMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()

		lrw := loggedResponseWriter{rw: w}
		f.ServeHTTP(&lrw, r)

		d := time.Since(s)

		code, text := lrw.statusCode, http.StatusText(lrw.statusCode)
		InfoLog.Printf("%s %s %d %s - %s", r.Method, r.URL.Path, code, text, d)
	}
}
