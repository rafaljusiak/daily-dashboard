package handlers

import (
	"net/http"
	"time"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/external"
	"github.com/rafaljusiak/daily-dashboard/income"
	"github.com/rafaljusiak/daily-dashboard/timeutils"
	"go.mongodb.org/mongo-driver/mongo"
)

func SyncHandler(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext) {
	previousMonth := time.Now().AddDate(0, -1, 0)

	startDate := timeutils.FirstDayOfMonth(previousMonth)
	endDate := timeutils.LastDayOfMonth(previousMonth)
	timeEntries, err := external.FetchClockifyTimeEntries(appCtx, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalTime, err := external.SumClockifyTime(timeEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exchangeRate, err := external.FetchNBPExchangeRate(appCtx.HTTPClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalIncome := income.Calculate(
		timeutils.MinutesToHours(totalTime),
		appCtx.Config.HourlyRate,
		exchangeRate,
	)

	formattedDate := previousMonth.Format("2006-01-02")
	incomeDoc, err := income.GetDocumentByDate(
		appCtx,
		r.Context(),
		formattedDate,
	)
	if err == mongo.ErrNoDocuments {
		incomeDoc = income.IncomeDocument{
			Date:         formattedDate,
			TimeHours:    timeutils.MinutesToHours(totalTime),
			HourlyRate:   appCtx.Config.HourlyRate,
			ExchangeRate: exchangeRate,
			TotalIncome:  totalIncome,
		}

		err = income.InsertDocument(appCtx, r.Context(), incomeDoc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		incomeDoc.TimeHours = timeutils.MinutesToHours(totalTime)
		incomeDoc.ExchangeRate = exchangeRate
		incomeDoc.TotalIncome = totalIncome
		err = income.UpdateDocument(appCtx, r.Context(), incomeDoc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
