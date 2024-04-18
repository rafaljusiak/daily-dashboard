package external

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/data"
)

const apiUrl string = "https://api.clockify.me/api/v1/"

func userURL() string {
	url, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	url = url.JoinPath("user")
	return url.String()
}

func timeEntriesURL(ctx *app.Context, userId string) string {
	url, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	url = url.JoinPath("workspaces", ctx.Config.WorkspaceId, "user", userId, "time-entries")

	urlQuery := url.Query()
	urlQuery.Add("page-size", "5000")
	urlQuery.Add("start", "2024-04-01T00:00:00Z") // TODO get start of the current month

	url.RawQuery = urlQuery.Encode()
	return url.String()
}

func FetchUserId(ctx *app.Context) (string, error) {
	client := ctx.HTTPClient

	url := userURL()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("user-agent", "Go/Daily-Dashboard")
	req.Header.Add("x-api-key", ctx.Config.ClockifyApiKey)

	log.Printf("Sending request to %s", url)
	response, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	userData, err := data.ReadJson(body)
	if err != nil {
		return "", err
	}
	return userData["id"].(string), nil
}

func FetchTimeEntries(ctx *app.Context) ([]interface{}, error) {
	client := ctx.HTTPClient

	userId, err := FetchUserId(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", timeEntriesURL(ctx, userId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", ctx.Config.ClockifyApiKey)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	arr := []interface{}{}
	if json.Unmarshal(body, &arr) != nil {
		return nil, err
	}

	return arr, nil
}
