package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/external"
	"github.com/rafaljusiak/daily-dashboard/income"
	"github.com/rafaljusiak/daily-dashboard/timeutils"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNoClockifyData = errors.New("no clockify data")

func SyncHandler(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext) {
	processedDate := time.Now()

	for {
		processedDate = processedDate.AddDate(0, 0, -1)
		err := processMonth(appCtx, r, processedDate)

		if err == ErrNoClockifyData {
			break
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/archive", http.StatusFound)
}

func processMonth(
	appCtx *app.AppContext,
	r *http.Request,
	processedMonth time.Time,
) error {
	startDate := timeutils.FirstDayOfMonth(processedMonth)
	endDate := timeutils.LastDayOfMonth(processedMonth)
	timeEntries, err := external.FetchClockifyTimeEntries(appCtx, startDate, endDate)
	if err != nil {
		return err
	}

	totalTime, err := external.SumClockifyTime(timeEntries)
	if err != nil {
		return err
	}

	if totalTime == 0 {
		return ErrNoClockifyData
	}

	exchangeRate, err := external.FetchNBPExchangeRate(appCtx.HTTPClient)
	if err != nil {
		return err
	}

	totalIncome := income.Calculate(
		timeutils.MinutesToHours(totalTime),
		appCtx.Config.HourlyRate,
		exchangeRate,
	)

	formattedDate := endDate.Format("2006-01-02")
	incomeDoc, err := income.GetDocumentByDate(
		appCtx,
		r.Context(),
		formattedDate,
	)

	if err == mongo.ErrNoDocuments {
		log.Printf("No document for the %s month, inserting new one", formattedDate)
		incomeDoc = income.IncomeDocument{
			Date:         formattedDate,
			TimeMinutes:  totalTime,
			HourlyRate:   appCtx.Config.HourlyRate,
			ExchangeRate: exchangeRate,
			TotalIncome:  totalIncome,
		}

		err = income.InsertDocument(appCtx, r.Context(), incomeDoc)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		log.Println("Updating existing document")
		incomeDoc.TimeMinutes = totalTime
		incomeDoc.ExchangeRate = exchangeRate
		incomeDoc.TotalIncome = totalIncome
		err = income.UpdateDocument(appCtx, r.Context(), incomeDoc)
		if err != nil {
			return err
		}
	}

	return nil
}
