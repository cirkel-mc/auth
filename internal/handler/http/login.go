package http

import (
	"cirkel/auth/internal/domain/dto"
	"time"

	"github.com/cirkel-mc/goutils/constants"
	"github.com/cirkel-mc/goutils/env"
	"github.com/cirkel-mc/goutils/logger"
	svc "github.com/cirkel-mc/goutils/service"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/cirkel-mc/goutils/zone"
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

	// set cookie access token
	expired := time.Now().In(zone.TzJakarta()).Add(time.Duration(resp.ExpiresIn) * time.Second)
	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieAccessToken,
		Value:    resp.AccessToken,
		Path:     "/",
		Domain:   env.GetString("CIRKEL_URL"),
		MaxAge:   int(resp.ExpiresIn),
		Expires:  expired,
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
	// set cookie for refresh token
	expired = time.Now().In(zone.TzJakarta()).Add(time.Duration(resp.RefreshTokenExpiresIn) * time.Second)
	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieRefreshToken,
		Value:    resp.RefreshToken,
		Path:     "/",
		Domain:   env.GetString("CIRKEL_URL"),
		MaxAge:   int(resp.RefreshTokenExpiresIn),
		Expires:  expired,
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
	// set cookie for csrf token
	expired = time.Now().In(zone.TzJakarta()).Add(env.GetDuration("HTTP_CSRF_TOKEN_DURATION", 30*time.Minute))
	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieCsrfToken,
		Value:    resp.CsrfToken,
		Path:     "/",
		Domain:   env.GetString("CIRKEL_URL"),
		MaxAge:   int(env.GetDuration("HTTP_CSRF_TOKEN_DURATION", 30*time.Minute)),
		Expires:  expired,
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return sc.OK(ctx, "Login berhasil")
}
