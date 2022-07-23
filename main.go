package main

import (
	"debug/elf"
	"fmt"
	"log"
	"net/http"
	"os"
)

//自定义hellohandler
type getEnvHandler struct{}

func (m *getEnvHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, r.Header)
	envVersion := os.Getenv("VERSION")
	if len(envVersion) != 0 {
		w.Header().set("VERSION", envVersion)
		w.Write([]byte(envVersion))
	}
	elf  {
		w.Write([]byte("VERSION is not exist"))
	}

	// envPath := os.Getenv("PATH")
	// os.Setenv("VERSION", "go version go1.18.4 windows/amd64")
	// for k, v := range os.Environ() {
	// 	fmt.Fprintln(w, k)
	// 	fmt.Fprintln(w, v)
	// }
	// w.Write([]byte(envVersion))
	// fmt.Fprintln(w, envVersion)
	// fmt.Fprintln(w, envPath)
	// r.Header.Set("VERSION", envVersion)
	// fmt.Fprint(w, r.Header.Get("VERSION"))
}

//自定义hellohandler 结束

//自定义abouthandler
type aboutHandler struct{}

func (m *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("about!!"))
}

//自定义abouthandler 结束

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("200"))
	log.Print("health check")
}

func main() {

	getenv := getEnvHandler{}
	about := aboutHandler{}

	// server := http.Server{
	// 	Addr: "localhost:8080",
	// 	// Handler: &mh,  //指定myHandler 就没有了多路路由
	// 	Handler: nil, //DefaultServe Mux
	// }
	http.Handle("/", http.FileServer(http.Dir("blog")))
	// http.Handle("/", http.FileServer(http.File("index.html")))
	http.Handle("/getversion", &getenv)
	http.Handle("/about", &about)

	// server.ListenAndServe()
	// http.ListenAndServe("localhost:8080", http.FileServer(http.Dir("blog"))) //DefaultServe Mux
	http.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is Request header information\n"))

		if len(r.Header) > 0 {
			for k, v := range r.Header {
				w.Header().Set(k, v[0])
				fmt.Fprint(w, k)
				fmt.Fprint(w, v[0])
			}

			r.ParseForm()
			if len(r.Form) > 0 {
				for k, v := range r.Form {
					log.Printf("%s=%s", k, v[0])
				}
			}
			// fmt.Fprint(w, r.Header["Accept-Encoding"])
			// fmt.Fprint(w, r.Header.Get("Accept-Encoding"))
		}
	})

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
