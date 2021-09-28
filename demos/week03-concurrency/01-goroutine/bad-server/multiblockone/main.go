package main

import (
	"fmt"
	"net/http"
)

func serveApp()  {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "你好 world")
	})
	_ = http.ListenAndServe("0.0.0.0:8080", mux)
}

func serveDebug()  {
	_ = http.ListenAndServe("127.0.0.1:8001", http.DefaultServeMux)
}

func main() {
	go serveDebug()
	serveApp()
}
