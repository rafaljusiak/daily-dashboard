package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/calc"
	"github.com/rafaljusiak/daily-dashboard/external"
	"github.com/rafaljusiak/daily-dashboard/timeutils"
)

type DashboardData struct {
	AlreadyWorked   string
	CurrentIncome   float64
	ExchangeRate    float64
	MinimumHours    string
	OptimalIncome   float64
	WorkingHours    string
	WorkingTimeDiff string
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
	exchangeRateChan := make(chan float64)
	timeEntriesChan := make(chan []external.ClockifyTimeEntryData)

	errChan := make(chan error, 2)

	go func() {
		exchangeRate, err := external.FetchNBPExchangeRate(ctx.HTTPClient)
		if err != nil {
			errChan <- fmt.Errorf("error while fetching NBP exchange rate: %v", err)
			return
		}
		exchangeRateChan <- exchangeRate
	}()

	go func() {
		timeEntries, err := external.FetchTimeEntries(ctx)
		if err != nil {
			errChan <- fmt.Errorf("error while fetching Clockify time entries: %v", err)
			return
		}
		timeEntriesChan <- timeEntries
	}()

	var exchangeRate float64
	var timeEntries []external.ClockifyTimeEntryData
	for i := 0; i < 2; i++ {
		select {
		case rate := <-exchangeRateChan:
			exchangeRate = rate
		case entries := <-timeEntriesChan:
			timeEntries = entries
		case err := <-errChan:
			return nil, err
		}
	}

	alreadyWorkedMinutes, err := calc.SumDuration(timeEntries)
	if err != nil {
		return nil, fmt.Errorf("error calculating already worked minutes: %v", err)
	}

	currentIncome := calc.Income(
		timeutils.MinutesToHours(alreadyWorkedMinutes),
		ctx.Config.HourlyRate,
		exchangeRate,
	)

	workingHours := timeutils.WorkingHoursForCurrentMonth()
	optimalIncome := calc.Income(
		float64(workingHours),
		ctx.Config.HourlyRate,
		exchangeRate,
	)

	minimumHours := timeutils.WorkingHoursUntilToday()
	workingTimeDiff := -(minimumHours*60 - alreadyWorkedMinutes)

	data := &DashboardData{
		AlreadyWorked:   timeutils.MinutesToString(alreadyWorkedMinutes),
		CurrentIncome:   currentIncome,
		ExchangeRate:    exchangeRate,
		MinimumHours:    timeutils.HoursToString(minimumHours),
		OptimalIncome:   optimalIncome,
		WorkingHours:    timeutils.HoursToString(workingHours),
		WorkingTimeDiff: timeutils.MinutesToString(workingTimeDiff),
	}

	return data, nil
}
