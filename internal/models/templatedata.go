package models

import "github.com/KiroLakestrike/BedAndBreakfast/internal/forms"

// TemplateData holds data sent from handlers to Templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float64
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
