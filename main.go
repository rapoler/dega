package main

import (
	"net/http"
	"os"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core"
	"github.com/factly/dega-server/service/factcheck"
	"github.com/factly/dega-server/util"
	"github.com/factly/x/loggerx"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	_ "github.com/factly/dega-server/docs" // docs is generated by Swag CLI, you have to import it.
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Dega API
// @version 1.0
// @description Dega server API

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8820
// @BasePath /
func main() {

	godotenv.Load()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8820"
	}

	port = ":" + port

	// db setup
	config.SetupDB()

	// open log file
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	if err == nil {
		r.Use(loggerx.NewLogger(file))
	}
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	/* disable swagger in production */
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.With(util.CheckUser, util.CheckSpace).Group(func(r chi.Router) {
		r.Mount("/factcheck", factcheck.Router())
		r.Mount("/core", core.Router())
	})

	http.ListenAndServe(port, r)
}
