package syslogger

import (
	"log"
	"net/http"
)

type Handler struct {
	Next http.Handler
}

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bytesSent uint64
}

func (handler *Handler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	var writer = &loggerResponseWriter{wr, http.StatusOK, 0}
	handler.Next.ServeHTTP(writer, req)
	log.Printf("%s %s %s %d %d", req.RemoteAddr, req.Method, req.RequestURI, writer.statusCode, writer.bytesSent)
}

func (writer *loggerResponseWriter) WriteHeader(statusCode int) {
	writer.statusCode = statusCode
	writer.ResponseWriter.WriteHeader(statusCode)
}

func (writer *loggerResponseWriter) Write(buffer []byte) (sent int, err error) {
	sent, err = writer.ResponseWriter.Write(buffer)

	writer.bytesSent += uint64(sent)
	return
}