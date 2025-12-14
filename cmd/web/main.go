package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alireza-akbarzadeh/bookings/cmd/web/routes"
	"github.com/alireza-akbarzadeh/bookings/pkg/config"
	"github.com/alireza-akbarzadeh/bookings/pkg/handlers"
	"github.com/alireza-akbarzadeh/bookings/pkg/render"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {
	// change this true
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tmpl, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalf("cannot create template cache: %v", err)
	}
	app.TemplateCache = tmpl
	app.UseCache = false
	render.NewTemplate(&app)

	port := flag.Int("port", 8080, "Port run on the sever")
	flag.Parse()

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("🚀 Server starting on http://localhost%s\n", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: routes.Setup(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
