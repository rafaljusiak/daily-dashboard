package handlers

import (
	"net/http"

	"github.com/rafaljusiak/daily-dashboard/app"
)

func SyncHandler(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext) {
	// TODO: implement
	http.Redirect(w, r, "/", http.StatusFound)
}
