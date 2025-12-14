package routes

import (
	"net/http"

	"github.com/alireza-akbarzadeh/bookings/pkg/config"
	"github.com/alireza-akbarzadeh/bookings/pkg/handlers"
	"github.com/alireza-akbarzadeh/bookings/pkg/middleware"
	"github.com/go-chi/chi"
	chiMid "github.com/go-chi/chi/middleware"
)

func Setup(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(chiMid.Recoverer)
	mux.Use(middleware.Logger())
	mux.Use(middleware.NoSurf(app))
	mux.Use(middleware.SessionLoad(app))

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	return mux
}
