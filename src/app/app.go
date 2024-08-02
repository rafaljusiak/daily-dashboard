package app

import (
	"net/http"
	"time"
)

type AppContext struct {
	Config     Config
	HTTPClient *http.Client
}

func NewAppContext() *AppContext {
	return &AppContext{
		Config: LoadConfig(),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
