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
	Errors          []string
	ExchangeRate    float64
	MinimumHours    string
	OptimalIncome   float64
	WeatherForecast template.HTML
	WorkingHours    string
	WorkingTimeDiff string
}

func DashboardHandler(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext) {
	t, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := prepareDashboardData(appCtx)
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

func prepareDashboardData(appCtx *app.AppContext) (*DashboardData, error) {
	exchangeRateChan := make(chan float64)
	timeEntriesChan := make(chan []external.ClockifyTimeEntryData)
	wttrChan := make(chan string)

	errChan := make(chan error, 3)

	go func() {
		exchangeRate, err := external.FetchNBPExchangeRate(appCtx.HTTPClient)
		if err != nil {
			errChan <- fmt.Errorf("error while fetching NBP exchange rate: %v", err)
			return
		}
		exchangeRateChan <- exchangeRate
	}()

	go func() {
		timeEntries, err := external.FetchTimeEntries(appCtx)
		if err != nil {
			errChan <- fmt.Errorf("error while fetching Clockify time entries: %v", err)
			return
		}
		timeEntriesChan <- timeEntries
	}()

	go func() {
		wttrData, err := external.FetchWttrData(appCtx)
		if err != nil {
			errChan <- fmt.Errorf("error while fetching wttr.in data: %v", err)
			return
		}
		wttrChan <- wttrData
	}()

	var exchangeRate float64
	var timeEntries []external.ClockifyTimeEntryData
	var wttrString string

	errors := []string{}

	for i := 0; i < 3; i++ {
		select {
		case rate := <-exchangeRateChan:
			exchangeRate = rate
		case entries := <-timeEntriesChan:
			timeEntries = entries
		case wttrData := <-wttrChan:
			wttrString = wttrData
		case err := <-errChan:
			errors = append(errors, err.Error())
		}
	}

	alreadyWorkedMinutes, err := calc.SumDuration(timeEntries)
	if err != nil {
		errors = append(errors, fmt.Sprintf("error calculating already worked minutes: %v", err))
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
		Errors:          errors,
		ExchangeRate:    exchangeRate,
		MinimumHours:    timeutils.HoursToString(minimumHours),
		OptimalIncome:   optimalIncome,
		WeatherForecast: template.HTML(wttrString),
		WorkingHours:    timeutils.HoursToString(workingHours),
		WorkingTimeDiff: timeutils.MinutesToString(workingTimeDiff),
	}

	return data, nil
}
