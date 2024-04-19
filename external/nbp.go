package external

import (
	"encoding/json"
	"net/http"
)

type NBPCurrencyData struct {
	Rates []struct {
		Mid float64 `json:"mid"`
	} `json:"rates"`
}

const nbpUrl string = "https://api.nbp.pl/api/exchangerates/rates/A/USD/?format=json"

func FetchNBPExchangeRate(client *http.Client) (float64, error) {
	req, err := http.NewRequest("GET", nbpUrl, nil)
	if err != nil {
		return 0.0, err
	}
	req.Header.Add("User-Agent", "Go/Daily-Dashboard")

	response, err := client.Do(req)
	if err != nil {
		return 0.0, err
	}
	defer response.Body.Close()

	var responseData NBPCurrencyData
	err = json.NewDecoder(response.Body).Decode(&responseData)
	return responseData.Rates[0].Mid, err
}
