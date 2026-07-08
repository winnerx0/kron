package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	http.HandleFunc("/delay/5", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delayed success"))
	})

	http.HandleFunc("/delay/15", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delayed success"))
	})

	http.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed"))
	})

	http.HandleFunc("/not-found", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	http.HandleFunc("/flaky", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()

		if now%2 == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("temporary failure"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("recovered"))
	})

	fmt.Println("mock server running on :9090")
	http.ListenAndServe(":9090", nil)
}