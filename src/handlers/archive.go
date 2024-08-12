package handlers

import (
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
)

func ArchiveHandler(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext) {
	t, err := template.ParseFiles("templates/archive.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
