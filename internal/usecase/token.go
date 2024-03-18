package usecase

import (
	"cirkel/auth/internal/domain/dto"
	"cirkel/auth/internal/domain/model"
	"context"

	"github.com/cirkel-mc/goutils/errs"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/cirkel-mc/goutils/types"
)

func (u *usecaseInstance) generateTokens(ctx context.Context, req interface{}, user *model.User) (resp *dto.Token, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:GenerateTokens")
	defer trace.Finish()

	var publicKey, deviceId, channel string
	switch r := req.(type) {
	case *dto.RequestLogin:
		publicKey = r.PublicKey
		deviceId = r.DeviceId
		channel = r.Channel
	case *dto.RequestRegister:
		publicKey = r.PublicKey
		deviceId = r.DeviceId
		channel = r.Channel
	}

	tc := &types.TokenClaim{
		UserId:     user.Id,
		Username:   user.Username,
		Email:      user.Email,
		RoleId:     user.Role.Id,
		RoleKey:    user.Role.Key,
		Channel:    channel,
		DeviceId:   deviceId,
		UserStatus: user.Status.String(),
		IsPartner:  user.IsPartner,
		PublicKey:  publicKey,
	}

	resp, err = u.cache.SetAccessToken(ctx, tc)
	if err != nil {
		trace.SetError(err)

		return nil, errs.NewErrorWithCodeErr(err, errs.RedisError)
	}

	csrfToken, err := u.generateCsrfToken(ctx)
	if err != nil {
		trace.SetError(err)

		return nil, err
	}
	resp.CsrfToken = csrfToken

	return
}
