package middleware

import (
	"net/http"

	"github.com/alireza-akbarzadeh/bookings/pkg/config"
	"github.com/justinas/nosurf"
)

// NuSurf adds CSRF protection to all POST requests
func NoSurf(app *config.AppConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		csrfHandler := nosurf.New(next)
		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   app.InProduction,
			SameSite: http.SameSiteLaxMode,
		})
		return csrfHandler
	}
}

// SessionLoad load and save session on every request
func SessionLoad(app *config.AppConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return app.Session.LoadAndSave(next)
	}
}
