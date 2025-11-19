package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KiroLakestrike/BedAndBreakfast/internal/config"
	"github.com/KiroLakestrike/BedAndBreakfast/internal/handlers"
	"github.com/KiroLakestrike/BedAndBreakfast/internal/models"
	"github.com/KiroLakestrike/BedAndBreakfast/internal/render"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var app config.AppConfig

func main() {

	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

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
