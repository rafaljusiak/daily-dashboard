package main

import (
	"log"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/handlers"
)

func main() {
	log.Println("===================================")
	log.Println("=     Daily Dashboard by R.J.     =")
	log.Println("===================================")

	ctx := app.NewContext()
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.DashboardHandler(w, r, ctx)
	})
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/favicon.ico")
	})

	port := ctx.Config.Port
	log.Printf("Server is running on port %v", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		return
	}
}
