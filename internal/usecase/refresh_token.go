package usecase

import (
	"cirkel/auth/internal/domain/dto"
	"context"
	"database/sql"
	"errors"

	"github.com/cirkel-mc/goutils/errs"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/redis/go-redis/v9"
)

func (u *usecaseInstance) RefreshToken(ctx context.Context, req *dto.Token) (*dto.Token, error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:RefreshToken")
	defer trace.Finish()

	// find the old refreshToken to redis
	tc, err := u.cache.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		trace.SetError(err)
		if !errors.Is(err, redis.Nil) {
			return nil, errs.NewErrorWithCodeErr(err, errs.RedisError)
		}

		return nil, errs.NewErrorWithCodeErr(err, errs.DataNotFound)
	}

	// find user based on token
	user, err := u.psql.FindUserById(ctx, tc.GetUserId())
	if err != nil {
		trace.SetError(err)

		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewErrorWithCodeErr(err, errs.SQLError)
		}

		return nil, errs.NewErrorWithCodeErr(err, errs.UserNotFound)
	}

	resp, err := u.generateTokens(ctx, &dto.RequestLogin{}, user)
	if err != nil {
		trace.SetError(err)

		return nil, err
	}

	return resp, nil
}
