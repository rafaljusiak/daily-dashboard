package handlers

import (
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/calc"
	"github.com/rafaljusiak/daily-dashboard/external"
)

type DashboardData struct {
	CurrentIncome float64
	ExchangeRate  float64
	Minutes       string
	TimeEntries   []external.ClockifyTimeEntryData
}

func DashboardHandler(w http.ResponseWriter, r *http.Request, ctx *app.Context) {
	t, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := prepareDashboardData(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func prepareDashboardData(ctx *app.Context) (*DashboardData, error) {
	exchangeRate, err := external.FetchNBPExchangeRate(ctx.HTTPClient)
	if err != nil {
		return nil, err
	}

	timeEntries, err := external.FetchTimeEntries(ctx)
	if err != nil {
		return nil, err
	}
	minutes, _ := calc.SumDuration(timeEntries)

	currentIncome := (float64(minutes) / 60.0) * ctx.Config.HourlyRate * exchangeRate
	data := &DashboardData{
		CurrentIncome: currentIncome,
		ExchangeRate:  exchangeRate,
		Minutes:       calc.MinutesToString(minutes),
		TimeEntries:   timeEntries,
	}

	return data, nil
}
