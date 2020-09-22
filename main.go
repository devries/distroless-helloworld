package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		remoteAddr := r.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			remoteAddr = r.RemoteAddr
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response_string := fmt.Sprintf("Hello %s from %s", remoteAddr, hostname)

		io.WriteString(w, response_string)
	})

	logHandler := loggingHandler(mux)

	log.Printf("Starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), logHandler))
}

type statusRecorder struct {
	http.ResponseWriter
	status    int
	byteCount int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *statusRecorder) Write(p []byte) (int, error) {
	bc, err := rec.ResponseWriter.Write(p)
	rec.byteCount += bc

	return bc, err
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec := statusRecorder{w, 200, 0}
		next.ServeHTTP(&rec, req)
		remoteAddr := req.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			remoteAddr = req.RemoteAddr
		}
		log.Printf("%s - \"%s %s %s\" (%s) %d %d", remoteAddr, req.Method, req.URL.Path, req.Proto, req.Host, rec.status, rec.byteCount)
	})
}
