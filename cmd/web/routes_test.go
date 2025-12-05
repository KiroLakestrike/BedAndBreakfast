package main

import (
	"fmt"
	"testing"

	"github.com/KiroLakestrike/BedAndBreakfast/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do Nothing
	default:
		t.Error(fmt.Sprintf("routes type is not *chi.Mux, its type is: %T", v))
	}
}
