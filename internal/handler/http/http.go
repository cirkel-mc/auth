package http

import (
	"cirkel/auth/internal/usecase"

	"github.com/cirkel-mc/goutils/abstract"
	"github.com/cirkel-mc/goutils/validation"
	"github.com/gofiber/fiber/v2"
)

type httpInstance struct {
	validator  validation.Validation
	middleware abstract.Middleware
	usecase    usecase.Usecase
}

func New(m abstract.Middleware, u usecase.Usecase) abstract.RESTHandler {
	return &httpInstance{
		validator:  validation.New(),
		middleware: m,
		usecase:    u,
	}
}

func (h *httpInstance) Router(r fiber.Router) {
	r.Get("", h.validateAuth)

	{
		// for nginx middleware to validate authorization headers
		r.Get("/validate", h.validateAuth)
	}

	// v1 := r.Group("/v1", h.middleware.HTTPSignatureValidate)
	v1 := r.Group("/v1")
	v1.Post("/register", h.register)
	v1.Post("/login", h.login)
}
