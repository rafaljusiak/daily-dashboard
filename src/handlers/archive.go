package handlers

import (
	"html/template"
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
	"github.com/rafaljusiak/daily-dashboard/income"
)

type ArchiveData struct {
	IncomeDocs []income.IncomeDocument
	Errors     []string
}

func ArchiveHandler(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext) {
	t, err := template.ParseFiles("templates/archive.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	documents, err := income.GetDocumentList(appCtx, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := ArchiveData{
		IncomeDocs: documents,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
