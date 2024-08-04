package handlers

import (
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
)

func SetupRouter(appCtx *app.AppContext) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("GET /", setHandler(DashboardHandler, appCtx, withRequiredAuth))
	router.Handle("GET /login", setHandler(GetLoginHandler, appCtx))
	router.Handle("POST /login", setHandler(PostLoginHandler, appCtx))

	router.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/favicon.ico")
	})

	return router
}

type Middleware func(http.Handler, *app.AppContext) http.Handler

func withAppContext(
	handler func(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext),
	appCtx *app.AppContext,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, appCtx)
	})
}

func withRequiredAuth(next http.Handler, appCtx *app.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.CheckAuth(w, r, appCtx)
		next.ServeHTTP(w, r)
	})
}

func setHandler(
	handler func(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext),
	appCtx *app.AppContext,
	middlewares ...Middleware,
) http.Handler {
	h := withAppContext(handler, appCtx)
	for _, middleware := range middlewares {
		h = middleware(h, appCtx)
	}
	return h
}
