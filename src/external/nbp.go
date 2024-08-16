package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaljusiak/daily-dashboard/timeutils"
)

type NBPCurrencyData struct {
	Rates []struct {
		Mid float64 `json:"mid"`
	} `json:"rates"`
}

func FetchNBPExchangeRate(client *http.Client) (float64, error) {
	url := "https://api.nbp.pl/api/exchangerates/rates/A/USD/?format=json"
	req, err := PrepareHTTPRequest(url)
	if err != nil {
		return 0.0, err
	}

	response, err := client.Do(req)
	if err != nil {
		return 0.0, err
	}
	defer response.Body.Close()

	var responseData NBPCurrencyData
	err = json.NewDecoder(response.Body).Decode(&responseData)
	return responseData.Rates[0].Mid, err
}

func FetchArchivalNBPExchangeRate(client *http.Client, date time.Time) (float64, error) {
	dateUntil := timeutils.GetMonthLaterOrToday(date)

	url := fmt.Sprintf(
		"https://api.nbp.pl/api/exchangerates/rates/a/usd/%s/%s/?format=json",
		date.Format("2006-01-02"),
		dateUntil.Format("2006-01-02"),
	)

	req, err := PrepareHTTPRequest(url)
	if err != nil {
		return 0.0, err
	}

	response, err := client.Do(req)
	if err != nil {
		return 0.0, err
	}
	defer response.Body.Close()

	var responseData NBPCurrencyData
	err = json.NewDecoder(response.Body).Decode(&responseData)
	return responseData.Rates[0].Mid, err
}
