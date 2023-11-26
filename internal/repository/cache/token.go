package cache

import (
	"comrades-mc/auth/internal/domain/dto"
	"context"
	"fmt"
	"time"

	"github.com/comrades-mc/goutils/constants"
	"github.com/comrades-mc/goutils/convert"
	"github.com/comrades-mc/goutils/env"
	"github.com/comrades-mc/goutils/helpers"
	"github.com/comrades-mc/goutils/logger"
	"github.com/comrades-mc/goutils/tracer"
	"github.com/comrades-mc/goutils/types"
)

func (c *cacheRepository) GetAccessToken(ctx context.Context, accessToken string) (resp *types.TokenClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "CacheRepository:GetAccessToken")
	defer trace.Finish()

	resp = new(types.TokenClaim)
	key := fmt.Sprintf(prefixAccessToken, accessToken)
	err = c.client.GetStruct(ctx, resp, key)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed get access_token: %s", err)

		return
	}

	return
}

func (c *cacheRepository) GetRefreshToken(ctx context.Context, refreshToken string) (resp *types.TokenClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "CacheRepository:GetRefreshToken")
	defer trace.Finish()

	resp = new(types.TokenClaim)
	key := fmt.Sprintf(prefixRefreshToken, refreshToken)
	err = c.client.GetStruct(ctx, resp, key)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed get refresh_token: %s", err)

		return
	}

	return
}

func (c *cacheRepository) SetAccessToken(ctx context.Context, tc *types.TokenClaim) (resp *dto.Token, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "CacheRepository:SetAccessToken")
	defer trace.Finish()

	value, _ := convert.InterfaceToString(tc)

	accessToken, err := helpers.GenerateRandomString(64)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed when generate random_string: %s", err)

		return
	}
	key := fmt.Sprintf(prefixAccessToken, accessToken)
	accessTokenExpired := env.GetDuration("HTTP_ACCESS_TOKEN_DURATION", 24*time.Hour)
	// set into redis
	err = c.client.Set(ctx, key, value, accessTokenExpired)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed set access_token into redis: %s", err)

		return
	}

	refreshToken, err := helpers.GenerateRandomString(64)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed when generaee random_string for refresh_token: %s", err)

		return
	}
	key = fmt.Sprintf(prefixRefreshToken, refreshToken)
	refreshTokenExpired := env.GetDuration("HTTP_REFRESH_TOKEN_DURATION", 7*24*time.Hour)
	// set into redis
	err = c.client.Set(ctx, key, value, refreshTokenExpired)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed set refresh_token into redis: %s", err)

		return
	}

	resp = &dto.Token{
		TokenType:    constants.Bearer,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTokenExpired.Seconds()),
	}
	return
}
