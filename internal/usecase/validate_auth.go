package usecase

import (
	"cirkel/auth/internal/domain/constant"
	"cirkel/auth/internal/domain/dto"
	"context"
	"fmt"
	"strings"

	"github.com/cirkel-mc/goutils/errs"
	"github.com/cirkel-mc/goutils/logger"
	"github.com/cirkel-mc/goutils/tracer"
)

func (u *usecaseInstance) ValidateAuth(ctx context.Context, req *dto.RequestHeader) (resp *dto.Auth, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:ValidateAuth")
	defer trace.Finish()

	auths := strings.Split(req.Authorization, " ")
	if len(auths) > 1 {
		if auths[0] != constant.Basic {
			err = fmt.Errorf("authorization invalid")

			return nil, errs.NewErrorWithCodeErr(err, errs.InvalidAuth)
		}

		resp, err = u.basicAuth(ctx, auths[1])
		if err != nil {
			return
		}

		return
	}

	resp, err = u.auth(ctx, req.Authorization)
	if err != nil {
		return
	}

	return
}

func (u *usecaseInstance) auth(ctx context.Context, token string) (*dto.Auth, error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:Auth")
	defer trace.Finish()

	tc, err := u.cache.GetAccessToken(ctx, token)
	if err != nil {
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		return nil, errs.NewErrorWithCodeErr(err, errs.SessionExpired)
	}

	return &dto.Auth{
		UserId:    tc.GetUserId(),
		Username:  tc.GetUsername(),
		Email:     tc.GetEmail(),
		RoleId:    tc.GetRoleId(),
		RoleKey:   tc.GetRoleKey(),
		DeviceId:  tc.GetDeviceId(),
		Channel:   tc.GetChannel(),
		PublicKey: tc.GetPublicKey(),
	}, nil
}
