package http

import (
	"cirkel/auth/internal/domain/dto"

	"github.com/cirkel-mc/goutils/constants"
	"github.com/cirkel-mc/goutils/convert"
	svc "github.com/cirkel-mc/goutils/service"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/gofiber/fiber/v2"
)

func (h *httpInstance) validateAuth(c *fiber.Ctx) error {
	sc := svc.New(c, svc.Auth)
	ctx := c.UserContext()
	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPHandler:ValidateAuth")
	defer trace.Finish()

	var req = new(dto.RequestHeader)
	err := h.validator.BindAndValidateWithContext(ctx, c, req)
	if err != nil {
		trace.SetError(err)

		return sc.Error(ctx, err)
	}

	auth, err := h.usecase.ValidateAuth(ctx, req)
	if err != nil {
		trace.SetError(err)

		return sc.Error(ctx, err)
	}

	c.Set(constants.UserId, convert.IntToString(int64(auth.UserId)))
	c.Set(constants.Username, auth.Username)
	c.Set(constants.Email, auth.Email)
	c.Set(constants.RoleId, convert.IntToString(int64(auth.RoleId)))
	c.Set(constants.RoleKey, auth.RoleKey)
	c.Set(constants.DeviceId, auth.DeviceId)
	c.Set(constants.Channel, auth.Channel)
	c.Set(constants.PublicKey, auth.PublicKey)

	return sc.OK(ctx, "ok")
}
