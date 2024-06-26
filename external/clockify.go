package external

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/timeutils"
)

const apiUrl string = "https://api.clockify.me/api/v1/"

type ClockifyUserData struct {
	Id string `json:"id"`
}

type ClockifyTimeEntryData struct {
	TimeInterval struct {
		Duration string `json:"duration"`
	} `json:"timeInterval"`
}

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
	urlQuery.Add("start", timeutils.FirstDayOfMonth(time.Now()).Format(time.RFC3339))

	url.RawQuery = urlQuery.Encode()
	return url.String()
}

func FetchUserId(ctx *app.Context) (string, error) {
	client := ctx.HTTPClient
	req, err := PrepareHTTPRequest(userURL())
	if err != nil {
		return "", err
	}
	req.Header.Add("x-api-key", ctx.Config.ClockifyApiKey)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseData ClockifyUserData
	err = json.NewDecoder(response.Body).Decode(&responseData)

	return responseData.Id, err
}

func FetchTimeEntries(ctx *app.Context) ([]ClockifyTimeEntryData, error) {
	client := ctx.HTTPClient

	userId, err := FetchUserId(ctx)
	if err != nil {
		return nil, err
	}

	req, err := PrepareHTTPRequest(timeEntriesURL(ctx, userId))
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", ctx.Config.ClockifyApiKey)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseData []ClockifyTimeEntryData
	err = json.NewDecoder(response.Body).Decode(&responseData)

	return responseData, err
}
