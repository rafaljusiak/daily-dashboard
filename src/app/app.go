package app

import (
	"net/http"
)

type Context struct {
	Config     Config
	HTTPClient *http.Client
}

func NewContext() *Context {
	return &Context{
		Config:     LoadConfig(),
		HTTPClient: &http.Client{},
	}
}
