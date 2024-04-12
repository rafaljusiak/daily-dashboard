package external

import (
	"io"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/data"
)

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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0.0, err
	}

	nbpData, err := data.ReadJson(body)
	if err != nil {
		return 0.0, err
	}

	rates := nbpData["rates"].([]interface{})
	mid := rates[0].(map[string]interface{})["mid"].(float64)

	return mid, nil
}
