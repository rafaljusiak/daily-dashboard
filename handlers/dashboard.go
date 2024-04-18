package handlers

import (
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/external"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request, ctx *app.Context) {
	t, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exchangeRate, err := external.FetchNBPExchangeRate(ctx.HTTPClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, exchangeRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
