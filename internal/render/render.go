package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/KiroLakestrike/BedAndBreakfast/internal/config"
	"github.com/KiroLakestrike/BedAndBreakfast/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds default data to the template data struct before rendering
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html/template and writes output to http.ResponseWriter
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	// Get the template cache from AppConfig depending on whether caching is enabled
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		// Create a template cache on the fly if caching is disabled
		tc, _ = CreateTemplateCache()
	}
	// Retrieve the specified template from the cache by name
	t, ok := tc[tmpl]
	if !ok {
		// Log fatal error and stop execution if template not found in cache
		log.Fatal("Could not get template from template cache")
	}

	// Create a buffer to temporarily hold the executed template output
	buf := new(bytes.Buffer)

	// Add any default data to the template data
	td = AddDefaultData(td, r)
	// Execute the template with the provided data, writing output to the buffer
	err := t.Execute(buf, td)

	if err != nil {
		// Log any error encountered during template execution
		log.Println(err)
	}

	// Write the rendered content from buffer to the HTTP response
	_, err = buf.WriteTo(w)
	if err != nil {
		// Log any error encountered during writing to response
		log.Println(err)
	}
}

// CreateTemplateCache compiles templates from files and caches them in a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	// Initialize a map to hold compiled templates with template file names as keys
	myCache := map[string]*template.Template{}

	// Search the templates directory for page template files with .page.tmpl extension
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err // Return current cache and error if file glob fails
	}

	// Iterate over each page template file found
	for _, page := range pages {
		// Extract the base filename (without directory) to use as template name
		name := filepath.Base(page)

		// Create a new template with the extracted name and parse the page file into it
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err // Return error if page parsing fails
		}

		// Search for layout templates with .layout.tmpl extension used for page wrapping
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err // Return error if layout glob fails
		}

		// If layout templates are found, parse and associate them with the current page template
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err // Return error if layout parsing fails
			}
		}

		// Save the fully parsed and associated template set in the cache by filename
		myCache[name] = ts
	}

	// Return the completed template cache map with no error
	return myCache, nil
}
