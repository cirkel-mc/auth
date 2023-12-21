package http

import (
	"cirkel/auth/internal/domain/dto"
	"net/http"

	"github.com/cirkel-mc/goutils/logger"
	"github.com/cirkel-mc/goutils/response"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/gofiber/fiber/v2"
)

func (h *httpInstance) register(c *fiber.Ctx) error {
	ctx := c.UserContext()
	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPHandler:Register")
	defer trace.Finish()

	var req = new(dto.RequestRegister)
	req.RequestHeader = new(dto.RequestHeader)
	err := h.validator.BindAndValidateWithContext(ctx, c, req)
	if err != nil {
		trace.SetError(err)

		return response.Error(ctx, err).JSON(c)
	}

	resp, err := h.usecase.Register(ctx, req)
	if err != nil {
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		return response.Error(ctx, err).JSON(c)
	}

	return response.Success(ctx, http.StatusCreated, resp).JSON(c)
}
