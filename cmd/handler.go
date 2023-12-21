package cmd

import (
	"cirkel/auth/internal/handler/http"
	"cirkel/auth/internal/usecase"

	"github.com/cirkel-mc/goutils/abstract"
	"github.com/cirkel-mc/goutils/env"
	"github.com/cirkel-mc/goutils/factory/server"
	"github.com/cirkel-mc/goutils/factory/server/rest"
)

func httpHandler(deps abstract.Dependency, uc usecase.Usecase) server.ServiceFunc {
	return server.SetRestHandler(
		http.New(deps.GetMiddleware(), uc),
		rest.SetHTTPPort(env.GetInt("HTTP_PORT", 9090)),
	)
}
