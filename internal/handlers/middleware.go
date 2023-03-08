package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

const GZIPType string = "gzip"

func isSupportsGZIP(encodings []string) bool {
	for _, encode := range encodings {
		if strings.Contains(encode, GZIPType) {
			return true
		}
	}
	return false
}

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (wr gzipWriter) Write(b []byte) (int, error) {
	return wr.Writer.Write(b)
}

func GZipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if isSupportsGZIP(request.Header.Values("Content-Encoding")) {
			gzReader, err := gzip.NewReader(request.Body)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gzReader.Close()
			request.Body = gzReader
		}
		if isSupportsGZIP(request.Header.Values("Accept-Encoding")) {
			gzWriter, err := gzip.NewWriterLevel(writer, gzip.BestSpeed)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gzWriter.Close()

			writer = gzipWriter{ResponseWriter: writer, Writer: gzWriter}
			writer.Header().Set("Content-Encoding", GZIPType)
		}
		next.ServeHTTP(writer, request)
	})
}
