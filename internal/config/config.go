package config

import "html/template"

// This is just a basic configuation file
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	PortNumber    string
}
