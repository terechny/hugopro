package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API"))
	})

	rp := NewReverseProxy("hugo", "1313")
	r.Handle("/*", rp.ReverseProxy(nextHandler))

	http.ListenAndServe(":8080", r)
}

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {

			next.ServeHTTP(w, r)

		} else {

			u, _ := url.Parse("http://hugo:1313")
			r.URL.Scheme = u.Scheme
			r.URL.Host = u.Host
			r.URL.Path = u.Path
			proxy := httputil.NewSingleHostReverseProxy(u)
			proxy.ServeHTTP(w, r)

		}

	})

}

const content = ``

func WorkerTest() {
	t := time.NewTicker(5 * time.Second)
	var b byte = 0
	for {
		select {
		case <-t.C:
			err := os.WriteFile("/app/static/_index.md", []byte(fmt.Sprintf(content, b)), 0644)
			if err != nil {
				log.Println(err)
			}
			b++
		}
	}
}

/*
func Worker() {
	t := time.NewTicker(5 * time.Second)
	var counter int
	for {
		select {
		case <-t.C:

			currentTime := time.Now().Format("2006-01-02 15:04:05")

			fmt.Println(currentTime, counter)

			err := os.WriteFile("/app/static/tasks/_index.md", []byte(fmt.Sprintf("Текущее время: %s\n\nСчетчик: %d", currentTime, counter)), 0644)

			if err != nil {
				log.Println(err)
			}

			counter++
		}
	}
}

*/
