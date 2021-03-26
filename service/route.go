package service

import (
	"fmt"
	"net/http"

	"github.com/factly/dega-server/config"
	"github.com/factly/x/healthx"

	_ "github.com/factly/dega-server/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/factly/dega-server/service/core"
	"github.com/factly/dega-server/service/core/action/meta"
	"github.com/factly/dega-server/service/core/action/request/organisation"
	"github.com/factly/dega-server/service/core/action/request/space"
	factCheck "github.com/factly/dega-server/service/fact-check"
	"github.com/factly/dega-server/service/podcast"
	"github.com/factly/dega-server/util"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(loggerx.Init())
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	// r.Use(middlewarex.GormRequestID(&config.DB))

	if viper.IsSet("mode") && viper.GetString("mode") == "development" {
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		fmt.Println("Swagger @ http://localhost:7789/swagger/index.html")
	}

	if viper.IsSet("iframely_url") {
		r.Mount("/meta", meta.Router())
	}

	sqlDB, _ := config.DB.DB()

	healthx.RegisterRoutes(r, healthx.ReadyCheckers{
		"database":    sqlDB.Ping,
		"keto":        util.KetoChecker,
		"kavach":      util.KavachChecker,
		"kratos":      util.KratosChecker,
		"meilisearch": util.MeiliChecker,
	})

	r.With(middlewarex.CheckUser, middlewarex.CheckSpace(1), util.GenerateOrganisation, middlewarex.CheckAccess("dega", 1, util.GetOrganisation)).Group(func(r chi.Router) {
		r.Mount("/core", core.Router())
		r.With(util.FactCheckPermission).Mount("/fact-check", factCheck.Router())
		r.With(util.PodcastPermission).Mount("/podcast", podcast.Router())
	})

	r.With(middlewarex.CheckUser).Group(func(r chi.Router) {
		r.Post("/core/requests/organisations", organisation.Create)
		r.With(middlewarex.CheckSpace(1)).Post("/core/requests/spaces", space.Create)
	})

	return r
}
