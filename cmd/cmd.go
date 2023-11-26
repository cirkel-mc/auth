package cmd

import (
	"comrades-mc/auth/internal/repository/cache"
	"comrades-mc/auth/internal/repository/psql"
	"comrades-mc/auth/internal/usecase"

	"github.com/comrades-mc/goutils/abstract"
	"github.com/comrades-mc/goutils/config"
	"github.com/comrades-mc/goutils/factory/server"
	"github.com/comrades-mc/goutils/middleware"
)

func Serve(cfg config.Config) *server.Server {
	deps := injectDependencies(cfg)

	// middleware
	mdl := middleware.New(nil)
	deps.SetMiddleware(mdl)

	// repositories
	psqlRepo := psql.New(deps.GetSQLDatabase(abstract.Master).Database(), deps.GetSQLDatabase(abstract.Slave).Database())
	cacheRepo := cache.New(deps.GetRedisDatabase().Client())

	// usecase
	uc := usecase.New(psqlRepo, cacheRepo)

	// initiates services
	svc := server.NewApplicationService(
		server.SetConfiguration(cfg),
		httpHandler(deps, uc),
	)

	// initiates server
	srv := server.New(svc)
	return srv
}
