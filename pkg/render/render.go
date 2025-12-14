package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/alireza-akbarzadeh/bookings/pkg/config"
	"github.com/alireza-akbarzadeh/bookings/pkg/models"
)

var app *config.AppConfig

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// NewTemplate sets the config for the template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

func Template(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// create a template cache
	var tc map[string]*template.Template
	if app.UseCache {
		// ge requested template form cache
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatalf("template %q not found", tmpl)
	}
	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files name .*page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	// ranee throw all the files started with *page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
