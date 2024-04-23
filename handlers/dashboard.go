package handlers

import (
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/calc"
	"github.com/rafaljusiak/daily-dashboard/external"
)

type DashboardData struct {
	AlreadyWorked string
	CurrentIncome float64
	ExchangeRate  float64
	TimeEntries   []external.ClockifyTimeEntryData
	WorkingHours  string
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

	alreadyWorkedMinutes, _ := calc.SumDuration(timeEntries)
	currentIncome := calc.Income(
		calc.MinutesToHours(alreadyWorkedMinutes), 
		ctx.Config.HourlyRate, 
		exchangeRate,
	)
	workingHours := calc.WorkingHoursForCurrentMonth()

	data := &DashboardData{
		AlreadyWorked: calc.MinutesToString(alreadyWorkedMinutes),
		CurrentIncome: currentIncome,
		ExchangeRate:  exchangeRate,
		TimeEntries:   timeEntries,
		WorkingHours:  calc.MinutesToString(workingHours),
	}

	return data, nil
}
