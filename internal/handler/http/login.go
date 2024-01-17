package http

import (
	"cirkel/auth/internal/domain/dto"

	"github.com/cirkel-mc/goutils/logger"
	svc "github.com/cirkel-mc/goutils/service"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/gofiber/fiber/v2"
)

func (h *httpInstance) login(c *fiber.Ctx) error {
	sc := svc.New(c, svc.Auth)
	ctx := c.UserContext()
	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPHandler:Login")
	defer trace.Finish()

	var req = new(dto.RequestLogin)
	req.RequestHeader = new(dto.RequestHeader)
	err := h.validator.BindAndValidateWithContext(ctx, c, req)
	if err != nil {
		trace.SetError(err)

		return sc.Error(ctx, err)
	}

	resp, err := h.usecase.Login(ctx, req)
	if err != nil {
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		return sc.Error(ctx, err)
	}

	return sc.OK(ctx, resp)
}
