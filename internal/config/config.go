package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// This is just a basic configuation file
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	PortNumber    string
	InfoLog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
