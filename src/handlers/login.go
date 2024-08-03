package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/rafaljusiak/daily-dashboard/app"
)

type LoginData struct {
	Errors []string
}

func GetLoginHandler(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext) {
	data := LoginData{Errors: []string{}}

	if errCookie, err := r.Cookie("error"); err == nil {
		if errCookie.Value != "" {
			data.Errors = append(data.Errors, errCookie.Value)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    "",
			HttpOnly: true,
		})
	}

	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request, appCtx *app.AppContext) {
	if r.FormValue("password") == appCtx.Config.Password {
		http.SetCookie(w, &http.Cookie{
			Name:     "sessionid",
			Value:    app.HashPassword(appCtx.Config.Password),
			HttpOnly: true,
			Expires:  time.Now().Add(365 * 24 * time.Hour),
		})

		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    "Wrong password",
			HttpOnly: true,
		})
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}
