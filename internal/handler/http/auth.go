package http

import (
	"cirkel/auth/internal/domain/dto"
	"net/http"

	"github.com/cirkel-mc/goutils/constants"
	"github.com/cirkel-mc/goutils/convert"
	"github.com/cirkel-mc/goutils/response"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/gofiber/fiber/v2"
)

func (h *httpInstance) validateAuth(c *fiber.Ctx) error {
	ctx := c.UserContext()
	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPHandler:ValidateAuth")
	defer trace.Finish()

	var req = new(dto.RequestHeader)
	err := h.validator.BindAndValidateWithContext(ctx, c, req)
	if err != nil {
		trace.SetError(err)

		return response.Error(ctx, err).JSON(c)
	}

	auth, err := h.usecase.ValidateAuth(ctx, req)
	if err != nil {
		trace.SetError(err)

		return response.Error(ctx, err).JSON(c)
	}

	c.Set(constants.UserId, convert.IntToString(int64(auth.UserId)))
	c.Set(constants.Username, auth.Username)
	c.Set(constants.Email, auth.Email)
	c.Set(constants.RoleId, convert.IntToString(int64(auth.RoleId)))
	c.Set(constants.RoleKey, auth.RoleKey)
	c.Set(constants.DeviceId, auth.DeviceId)
	c.Set(constants.Channel, auth.Channel)
	c.Set(constants.PublicKey, auth.PublicKey)

	return response.Success(ctx, http.StatusOK, "ok").JSON(c)
}
