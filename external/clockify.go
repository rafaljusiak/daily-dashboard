package external

import (
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
)

const clockifyUrl string = "https://api.clockify.me/api/v1/"

func GetClockifyData(client *http.Client, ctx *app.AppContext) {
	req, err := http.NewRequest("GET", clockifyUrl, nil)
	if err != nil {
		return
	}
	
	req.Header.Add("x-api-key", ctx.Config.ClockifyApiKey)
}
