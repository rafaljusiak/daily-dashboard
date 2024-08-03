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

func EnforceAuthMiddleware(next http.Handler, appCtx *AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionIdCookieName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if cookie.Value != HashPassword(appCtx.Config.Password) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
