package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafaljusiak/daily-dashboard/server"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", web.DashboardHandler).Methods("GET")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
