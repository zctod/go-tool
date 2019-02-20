package util_server

import (
	"net/http"
	"testing"
	"time"
)

var handlersMap = make(map[string]HandlersFunc)

type HandlersFunc func(http.ResponseWriter, *http.Request)

type Handles struct{}

func (*Handles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := handlersMap[r.URL.String()]; ok {
		h(w, r)
	}
}

func TestGracefulExitWeb(t *testing.T) {

	handlersMap["/hello"] = func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello, World!"))
	}

	var server = &http.Server{
		Addr:           ":8080",
		Handler:        &Handles{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = server.ListenAndServe()
	GracefulExitWeb(server)
}
