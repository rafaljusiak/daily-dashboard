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

func timeEntriesURL(appCtx *app.AppContext, userId string) string {
	url, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	url = url.JoinPath("workspaces", appCtx.Config.WorkspaceId, "user", userId, "time-entries")

	urlQuery := url.Query()
	urlQuery.Add("page-size", "5000")
	urlQuery.Add("start", timeutils.FirstDayOfMonth(time.Now()).Format(time.RFC3339))

	url.RawQuery = urlQuery.Encode()
	return url.String()
}

func FetchUserId(appCtx *app.AppContext) (string, error) {
	client := appCtx.HTTPClient
	req, err := PrepareHTTPRequest(userURL())
	if err != nil {
		return "", err
	}
	req.Header.Add("x-api-key", appCtx.Config.ClockifyApiKey)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseData ClockifyUserData
	err = json.NewDecoder(response.Body).Decode(&responseData)

	return responseData.Id, err
}

func FetchTimeEntries(appCtx *app.AppContext) ([]ClockifyTimeEntryData, error) {
	client := appCtx.HTTPClient

	userId, err := FetchUserId(appCtx)
	if err != nil {
		return nil, err
	}

	req, err := PrepareHTTPRequest(timeEntriesURL(appCtx, userId))
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", appCtx.Config.ClockifyApiKey)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseData []ClockifyTimeEntryData
	err = json.NewDecoder(response.Body).Decode(&responseData)

	return responseData, err
}

func SumClockifyTime(timeEntries []ClockifyTimeEntryData) (int, error) {
	minutes := 0
	for _, timeEntry := range timeEntries {
		convertedDuration, err := timeutils.ConvertDurationToMinutes(
			timeEntry.TimeInterval.Duration,
		)
		if err != nil {
			return 0, err
		}
		minutes += convertedDuration
	}
	return minutes, nil
}
