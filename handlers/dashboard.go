package handlers

import (
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/external"
)

type DashboardData struct {
	ExchangeRate float64
	TimeEntries  []external.ClockifyTimeEntryData
}

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

	timeEntries, err := external.FetchTimeEntries(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		ExchangeRate: exchangeRate,
		TimeEntries:  timeEntries,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
