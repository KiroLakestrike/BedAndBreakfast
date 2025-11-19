package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KiroLakestrike/BedAndBreakfast/internal/config"
	"github.com/KiroLakestrike/BedAndBreakfast/internal/handlers"
	"github.com/KiroLakestrike/BedAndBreakfast/internal/render"
)

func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = false
	app.PortNumber = ":8080"

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// Load Handlers
	srv := &http.Server{
		Addr:    app.PortNumber,
		Handler: routes(&app),
	}

	fmt.Println(fmt.Sprintf("Listening on http://localhost%v", app.PortNumber))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
