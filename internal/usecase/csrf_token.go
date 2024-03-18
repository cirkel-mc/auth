package usecase

import (
	"context"
	"errors"

	"github.com/cirkel-mc/goutils/errs"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/redis/go-redis/v9"
)

func (u *usecaseInstance) ValidateCsrfToken(ctx context.Context, csrfToken string) (string, error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:ValidateCsrfToken")
	defer trace.Finish()

	err := u.validateCsrfToken(ctx, csrfToken)
	if err != nil {
		trace.SetError(err)

		return "", err
	}

	// generate new csrfToken when csrfToken before is valid
	csrfToken, err = u.generateCsrfToken(ctx)
	if err != nil {
		trace.SetError(err)

		return "", err
	}

	return csrfToken, nil
}

func (u *usecaseInstance) generateCsrfToken(ctx context.Context) (string, error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:PGenerateCsrfToken")
	defer trace.Finish()

	csrfToken, err := u.cache.SetCsrfToken(ctx)
	if err != nil {
		trace.SetError(err)

		return "", errs.NewErrorWithCodeErr(err, errs.RedisInsertFailed)
	}

	return csrfToken, nil
}

func (u *usecaseInstance) validateCsrfToken(ctx context.Context, csrfToken string) error {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:PValidateCsrfToken")
	defer trace.Finish()

	err := u.cache.GetCsrfToken(ctx, csrfToken)
	if err != nil {
		trace.SetError(err)

		if !errors.Is(err, redis.Nil) {
			return errs.NewErrorWithCodeErr(err, errs.RedisError)
		}

		return errs.NewErrorWithCodeErr(err, errs.RedisNil)
	}

	// delete old csrfToken
	err = u.cache.DeleteCsrfToken(ctx, csrfToken)
	if err != nil {
		trace.SetError(err)

		return errs.NewErrorWithCodeErr(err, errs.RedisDeleteFailed)
	}

	return nil
}
