package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/bookings/pkg/config"
	"github.com/alireza-akbarzadeh/bookings/pkg/models"
	"github.com/alireza-akbarzadeh/bookings/pkg/render"
)

// Repo the repository used by the handlers.
var Repo *Repository

// Repository is the repository type.
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new Repository.
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers.
func NewHandlers(r *Repository) {
	Repo = r
}

// Home handles requests to the root endpoint.
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, "home.page.tmpl", &models.TemplateData{})
}

// About handles requests to the /about endpoint.
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)

	stringMap["test"] = "Hello, again."

	// Get the IP stored in session (set in Home handler)
	remoteIP := repo.App.Session.GetString(r.Context(), "remote_ip")

	// Add to template data
	stringMap["remote_ip"] = remoteIP

	render.Template(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
