package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/handlers"
)

func main() {
	log.Println("============================")
	log.Println("=     Daily Dashboard      =")
	log.Println("============================")

	ctx := context.Background()
	appCtx := app.NewAppContext(ctx)
	router := handlers.SetupRouter(appCtx)

	port := appCtx.Config.Port
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
