package handlers

import (
	"net/http"

	"github.com/KiroLakestrike/BedAndBreakfast/pkg/config"
	"github.com/KiroLakestrike/BedAndBreakfast/pkg/models"
	"github.com/KiroLakestrike/BedAndBreakfast/pkg/render"
)

var Repo *Repository

// Repository type is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home Handler for our Home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About Handler for our About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// Perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	// send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
