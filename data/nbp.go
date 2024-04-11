package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const nbpUrl string = "https://api.nbp.pl/api/exchangerates/rates/A/USD/?format=json"

func GetNBPExchangeRate() (float32, error) {
	response, err := http.Get(nbpUrl)

	if err != nil {
		return 0.0, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0.0, err
	}

	var nbpData map[string]interface{}
	if err := json.Unmarshal(body, &nbpData); err != nil {
		fmt.Println("Invalid JSON:", err)
		return 0.0, err
	}
	fmt.Println(nbpData)

	return 0.0, nil
}
