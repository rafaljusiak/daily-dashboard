package web

import (
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/external"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	t, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exchangeRate, err := external.GetNBPExchangeRate(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, exchangeRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
