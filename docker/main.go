package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("200"))
	log.Print("health check: working")
}

func getClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func index(w http.ResponseWriter, r *http.Request) {
	os.Setenv("VERSION", "v0.0.1")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	fmt.Printf("os version: %s \n", version)

	if len(r.Header) > 0 {
		for k, v := range r.Header {
			fmt.Printf("Header key: %s, Header value: %s \n", k, v[0])
			w.Header().Set(k, v[0])
		}
	}

	clientip := getClientIP(r)
	log.Printf("Success! Response code: %d", 200)
	log.Printf("Success! clientip: %s", clientip)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		length := r.ContentLength
		body := make([]byte, length)
		r.Body.Read(body)
		fmt.Fprintln(w, string(body))
	})

	mux := http.NewServeMux()
	// 06. debug
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	log.SetPrefix("[boyang]-[info]-")
	log.SetFlags(log.Ldate | log.Llongfile)
	err := http.ListenAndServe("localhost:8080", nil) //DefaultServe Mux
	if err != nil {
		log.Fatal(err)
	}
}
