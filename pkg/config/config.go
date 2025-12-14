package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	Infolog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
