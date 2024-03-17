package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/cirkel-mc/goutils/env"
	"github.com/cirkel-mc/goutils/helpers"
	"github.com/cirkel-mc/goutils/logger"
	"github.com/cirkel-mc/goutils/tracer"
)

func (c *cacheRepository) GetCsrfToken(ctx context.Context, csrfToken string) error {
	trace, ctx := tracer.StartTraceWithContext(ctx, "CacheRepository:GetCsrfToken")
	defer trace.Finish()

	key := fmt.Sprintf(prefixAccessToken, csrfToken)
	_, err := c.client.Get(ctx, key)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed get csrf_token: %s", err)

		return err
	}

	return nil
}

func (c *cacheRepository) SetCsrfToken(ctx context.Context) (string, error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "CacheRepository:SetCsrfToken")
	defer trace.Finish()

	csrfToken, err := helpers.GenerateRandomString(16)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "error when generate random string: %s", err)

		return "", err
	}
	now := time.Now().Unix()
	unix := fmt.Sprintf("%d", now)
	csrfToken += unix[len(unix)-8:]
	key := fmt.Sprintf(prefixCsrfToken, csrfToken)

	csrfTokenExpired := env.GetDuration("HTTP_CSRF_TOKEN_DURATION", 30*time.Minute)
	err = c.client.Set(ctx, key, "1", csrfTokenExpired)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed to set csrf_token: %s", err)

		return "", err
	}

	return csrfToken, nil
}

func (c *cacheRepository) DeleteCsrfToken(ctx context.Context, csrfToken string) error {
	trace, ctx := tracer.StartTraceWithContext(ctx, "CacheRepository:DeleteCsrfToken")
	defer trace.Finish()

	key := fmt.Sprintf(prefixCsrfToken, csrfToken)
	err := c.client.Del(ctx, key)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed to delete csrf_token: %s", err)

		return err
	}

	return nil
}
