package handlers

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/unbeman/ya-prac-mcas/internal/utils"
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

func DecryptMiddleware(privateKey *rsa.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, request *http.Request) {
			if privateKey != nil {
				chyper, err := io.ReadAll(request.Body)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}
				request.Body.Close()
				data, err := utils.GetDecryptedMessage(privateKey, chyper, request.Header.Get("Encrypted-Key"))
				if err != nil {
					log.Error(err)
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}

				request.Body = io.NopCloser(bytes.NewReader(data))
				request.ContentLength = int64(len(data))
			}
			next.ServeHTTP(writer, request)
		}
		return http.HandlerFunc(fn)
	}
}

func IPCheckerMiddleware(trustedSubnet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, request *http.Request) {
			if trustedSubnet != nil {
				clientIP := request.Header.Get("X-Real-IP")

				if err := utils.CheckIPBelongsNetwork(clientIP, trustedSubnet); err != nil {
					http.Error(writer, err.Error(), http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(writer, request)
		}
		return http.HandlerFunc(fn)
	}
}
