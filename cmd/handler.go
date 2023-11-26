package cmd

import (
	"comrades-mc/auth/internal/handler/http"
	"comrades-mc/auth/internal/usecase"

	"github.com/comrades-mc/goutils/abstract"
	"github.com/comrades-mc/goutils/env"
	"github.com/comrades-mc/goutils/factory/server"
	"github.com/comrades-mc/goutils/factory/server/rest"
)

func httpHandler(deps abstract.Dependency, uc usecase.Usecase) server.ServiceFunc {
	return server.SetRestHandler(
		http.New(deps.GetMiddleware(), uc),
		rest.SetHTTPPort(env.GetInt("HTTP_PORT", 9090)),
	)
}
