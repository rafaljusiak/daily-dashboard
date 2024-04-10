package web

import (
	"bufio"
	"html/template"
	"net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := http.Get("https://rafaljusiak.pl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	data := ""
	for i := 0; scanner.Scan(); i++ {
		data += scanner.Text()
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
