package main

import (
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)


type LogHandler struct {
	Logger *log.Logger
}


func (h LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v %v %v\n", r.Method, r.URL, r.Proto))
	
	for name, values := range r.Header {
		val := strings.Join(values, "")
		sb.WriteString(fmt.Sprintf("%v: %v\n", name, val))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v. Logging %v bytes of body.\n", err, len(body))
	}

	var msg string
	if len(body) > 0 {
		defer r.Body.Close()
		sb.WriteString("\n")
		if isPrintable(body) {
			msg = "Logging ASCII message"
			sb.WriteString(string(body))
		} else {
			msg = "Logging binary message (hexdump)"
			sb.WriteString(hex.Dump(body))
		}
	}

	h.Logger.Printf("%v\n%v\n\n", msg, sb.String())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}


func main() {
    optListen := flag.String("l", "127.0.0.1:8000", "Endpoint to listen on (client)")
    optCertPem := flag.String("cert-pem", "", "Path to x509 certificate to present to TLS clients")
    optCertKey := flag.String("cert-key", "", "Path to x509 certificate key")
	optLogPath := flag.String("log", "", "Path to log file")
    flag.Parse()

	if *optLogPath == "" {
		log.Fatalf("Option 'log' must be specified!\n")
	}

	file, err := os.OpenFile(*optLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file for writing: %v", err)
	}
	defer file.Close()

	useTLS := *optCertPem != "" && *optCertKey != ""
	if useTLS {
		serveHTTPS(*optListen, *optCertPem, *optCertKey, file)
	} else {
		serveHTTP(*optListen, file)
	}
}


func serveHTTPS(endpoint, certPath, keyPath string, logWriter io.Writer) {
	server := &http.Server{
		Addr: endpoint,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS10,
		},
		Handler: LogHandler{Logger: log.New(logWriter, "", log.Ldate | log.Ltime)},
	}

	log.Printf("Listening on %v (HTTPS).\n", endpoint)

	if err := server.ListenAndServeTLS(certPath, keyPath); err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}
}


func serveHTTP(endpoint string, logWriter io.Writer) {
	server := &http.Server{
		Addr: endpoint,
		Handler: LogHandler{Logger: log.New(logWriter, "", log.Ldate | log.Ltime)},
	}

	log.Printf("Listening on %v (HTTP).\n", endpoint)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}
}


func isPrintable(buf []byte) bool {
	for _, b := range(buf) {
		if b < 32 {
			return false
		}
	}

	return true
}
