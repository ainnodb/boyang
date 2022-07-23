package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

//自定义abouthandler
type aboutHandler struct{}

func (m *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("about!!"))
}

//自定义abouthandler 结束

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
			// w.Header().Set(k, v)
			// fmt.Fprint(w, k)
			// fmt.Fprint(w, v[0])
		}
	}

	clientip := getClientIP(r)
	//fmt.Println(r.RemoteAddr)
	log.Printf("Success! Response code: %d", 200)
	log.Printf("Success! clientip: %s", clientip)
}

func main() {
	about := aboutHandler{}
	// server := http.Server{
	// 	Addr: "localhost:8080",
	// 	// Handler: &mh,  //指定myHandler 就没有了多路路由
	// 	Handler: nil, //DefaultServe Mux
	// }
	http.HandleFunc("/", index)
	// http.Handle("/", http.FileServer(http.File("index.html")))
	http.Handle("/about", &about)

	// server.ListenAndServe()
	// http.ListenAndServe("localhost:8080", http.FileServer(http.Dir("blog"))) //DefaultServe Mux
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		length := r.ContentLength
		body := make([]byte, length)
		r.Body.Read(body)
		fmt.Fprintln(w, string(body))
	})

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		query := url.Query()

		id := query["id"]
		log.Panicln(id)

		name := query.Get("name")
		log.Panicln(name)
	})

	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Fprintln(w, r.Form)
		fmt.Fprintln(w, r.PostForm)
		fmt.Fprintln(w, r.FormValue("first_name"))
	})

	http.HandleFunc("/healthz", healthzHandler)
	log.SetPrefix("[boyang]-[info]-")
	log.SetFlags(log.Ldate | log.Llongfile)
	err := http.ListenAndServe("localhost:8080", nil) //DefaultServe Mux
	if err != nil {
		log.Fatal(err)
	}
}
