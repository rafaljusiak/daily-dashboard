package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/handlers"
)

func main() {
	log.Println("===================================")
	log.Println("=     Daily Dashboard by R.J.     =")
	log.Println("===================================")

	ctx := app.NewContext()
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		handlers.DashboardHandler(w, r, ctx)
	})
	router.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/favicon.ico")
	})

	port := ctx.Config.Port
	log.Printf("Server is running on port %v", port)

	s := &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		return
	}
}
