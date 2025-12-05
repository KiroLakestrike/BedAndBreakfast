package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"html/template"

	"github.com/KiroLakestrike/BedAndBreakfast/internal/config"
	"github.com/KiroLakestrike/BedAndBreakfast/internal/models"
	"github.com/KiroLakestrike/BedAndBreakfast/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true
	app.PortNumber = ":8080"

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewTemplates(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals", Repo.Generals)
	mux.Get("/presidential", Repo.Presidential)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/contact", Repo.Contact)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds csrf protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache compiles templates from files and caches them in a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// Initialize a map to hold compiled templates with template file names as keys
	myCache := map[string]*template.Template{}

	// Search the templates directory for page template files with .page.tmpl extension
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
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
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err // Return error if layout glob fails
		}

		// If layout templates are found, parse and associate them with the current page template
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
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
