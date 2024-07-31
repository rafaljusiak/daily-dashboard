package external

import (
	"fmt"
	"io"

	"github.com/rafaljusiak/daily-dashboard/app"
)

func FetchWttrData(ctx *app.Context) (string, error) {
	req, err := PrepareHTTPRequest(fmt.Sprintf("https://wttr.in/%s?T", ctx.Config.City))
	if err != nil {
		return "", err
	}

	response, err := ctx.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
