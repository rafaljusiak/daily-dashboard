package external

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/data"
)

const apiUrl string = "https://api.clockify.me/api/v1/"

func userURL() string {
	url, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	url.JoinPath("user")
	return url.String()
}

func timeEntriesURL(ctx *app.AppContext, userId int) string {
	url, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	url.JoinPath("workspaces", ctx.Config.WorkspaceId, "user", strconv.Itoa(userId), "time-entries")

	urlQuery := url.Query()
	urlQuery.Add("page-size", "5000")
	urlQuery.Add("start", "2024-04-01T00:00:00Z") // TODO get start of the current month

	url.RawQuery = urlQuery.Encode()
	return url.String()
}

func FetchUserId(client *http.Client, ctx *app.AppContext) (int, error) {
	req, err := http.NewRequest("GET", userURL(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("x-api-key", ctx.Config.ClockifyApiKey)

	response, err := client.Do(req)

	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	userData, err := data.ReadJson(body)
	if err != nil {
		return 0, err
	}
	userId := userData["id"].(int)

	return userId, nil
}

func FetchTimeEntries(client *http.Client, ctx *app.AppContext) ([]interface{}, error) {
	userId, err := FetchUserId(client, ctx)
	if err != nil {
		panic(err)
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
