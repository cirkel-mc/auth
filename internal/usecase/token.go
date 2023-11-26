package usecase

import (
	"comrades-mc/auth/internal/domain/dto"
	"comrades-mc/auth/internal/domain/model"
	"context"
	"fmt"

	"github.com/comrades-mc/goutils/convert"
	"github.com/comrades-mc/goutils/errs"
	"github.com/comrades-mc/goutils/tracer"
	"github.com/comrades-mc/goutils/types"
)

func (u *usecaseInstance) generateTokens(ctx context.Context, channel, deviceId string, user *model.User) (resp *dto.Token, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:GenerateTokens")
	defer trace.Finish()

	fmt.Println("sebelum typesss")
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

	fmt.Println("sebelum set nih")
	resp, err = u.cache.SetAccessToken(ctx, tc)
	if err != nil {
		trace.SetError(err)
		fmt.Println("yah gagakl set access token", err)

		return nil, errs.NewErrorWithCodeErr(err, errs.INSERT_REDIS_FAILED)
	}

	return
}
