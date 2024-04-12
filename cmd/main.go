package main

import (
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/handlers"
)

func main() {
	ctx := app.NewAppContext()
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.DashboardHandler(w, r, ctx)
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}
