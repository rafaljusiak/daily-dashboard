package handlers

import (
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/calc"
	"github.com/rafaljusiak/daily-dashboard/external"
)

type DashboardData struct {
	ExchangeRate float64
	Minutes      string
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

	minutes, _ := calc.SumDuration(timeEntries)
	data := DashboardData{
		ExchangeRate: exchangeRate,
		Minutes:      calc.MinutesToString(minutes),
		TimeEntries:  timeEntries,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
