package http

import (
	"cirkel/auth/internal/domain/dto"
	"fmt"
	"time"

	"github.com/cirkel-mc/goutils/constants"
	"github.com/cirkel-mc/goutils/convert"
	"github.com/cirkel-mc/goutils/env"
	"github.com/cirkel-mc/goutils/errs"
	"github.com/cirkel-mc/goutils/logger"
	svc "github.com/cirkel-mc/goutils/service"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/cirkel-mc/goutils/zone"
	"github.com/gofiber/fiber/v2"
)

func (h *httpInstance) validateAuth(c *fiber.Ctx) error {
	sc := svc.New(c, svc.Auth)
	ctx := c.UserContext()
	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPHandler:ValidateAuth")
	defer trace.Finish()

	logger.Log.Printf(ctx, "all request header: %v", c.GetReqHeaders())

	var req = new(dto.RequestHeader)
	err := h.validator.BindAndValidateWithContext(ctx, c, req)
	if err != nil {
		trace.SetError(err)

		return sc.Error(ctx, err)
	}

	token := c.Cookies(constants.CookieAccessToken)
	if token != "" {
		req.Authorization = token
	}

	auth, err := h.usecase.ValidateAuth(ctx, req)
	if err != nil {
		trace.SetError(err)

		return sc.Error(ctx, err)
	}

	ccsrfToken := c.Cookies(constants.CookieCsrfToken)
	if ccsrfToken == "" {
		err = fmt.Errorf("csrf_token is empty")
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		return sc.Error(ctx, errs.NewErrorWithCodeErr(err, errs.InvalidCsrfToken))
	}
	// validate csrf token and if valid will generate new csrf token
	csrfToken, err := h.usecase.ValidateCsrfToken(ctx, ccsrfToken)
	if err != nil {
		trace.SetError(err)

		return sc.Error(ctx, err)
	}

	// set new csrf token into cookie
	expired := time.Now().In(zone.TzJakarta()).Add(env.GetDuration("HTTP_CSRF_TOKEN_DURATION", 30*time.Minute))
	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieCsrfToken,
		Value:    csrfToken,
		Path:     "/",
		Domain:   env.GetString("CIRKEL_URL"),
		MaxAge:   int(env.GetDuration("HTTP_CSRF_TOKEN_DURATION", 30*time.Minute)),
		Expires:  expired,
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	// set all data users into headers
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
