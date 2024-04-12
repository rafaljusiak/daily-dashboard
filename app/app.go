package app

import (
	"net/http"
)

type AppContext struct {
	Config 		Config
	HTTPClient 	*http.Client
}

func NewAppContext() *AppContext {
	return &AppContext{
        Config: GetConfig(),
        HTTPClient: &http.Client{},
    }
}
