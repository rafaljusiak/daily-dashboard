package app

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

const SessionIdCookieName = "sessionid"

func HashPassword(password string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(password))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func CheckAuth(w http.ResponseWriter, r *http.Request, appCtx *AppContext) {
	cookie, err := r.Cookie(SessionIdCookieName)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if cookie.Value != HashPassword(appCtx.Config.Password) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

