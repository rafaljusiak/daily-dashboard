package web

import (
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/data"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exchangeRate, err := data.GetNBPExchangeRate()
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
