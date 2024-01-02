package usecase

import (
	"cirkel/auth/internal/domain/dto"
	"cirkel/auth/internal/domain/model"
	"context"

	"github.com/cirkel-mc/goutils/convert"
	"github.com/cirkel-mc/goutils/errs"
	"github.com/cirkel-mc/goutils/tracer"
	"github.com/cirkel-mc/goutils/types"
)

func (u *usecaseInstance) generateTokens(ctx context.Context, channel, deviceId string, user *model.User) (resp *dto.Token, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:GenerateTokens")
	defer trace.Finish()

	tc := &types.TokenClaim{
		UserId:       convert.IntToString(int64(user.Id)),
		Username:     user.Username,
		Email:        user.Email,
		RoleId:       convert.IntToString(int64(user.Role.Id)),
		RoleKey:      user.Role.Key,
		Channel:      channel,
		DeviceId:     deviceId,
		StatusVerify: user.Status.String(),
	}

	resp, err = u.cache.SetAccessToken(ctx, tc)
	if err != nil {
		trace.SetError(err)

		return nil, errs.NewErrorWithCodeErr(err, errs.INSERT_REDIS_FAILED)
	}

	return
}
