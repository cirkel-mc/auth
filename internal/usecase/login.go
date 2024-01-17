package usecase

import (
	"cirkel/auth/internal/domain/dto"
	"context"

	"github.com/cirkel-mc/goutils/config/database/dbc"
	"github.com/cirkel-mc/goutils/errs"
	"github.com/cirkel-mc/goutils/logger"
	"github.com/cirkel-mc/goutils/tracer"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecaseInstance) Login(ctx context.Context, req *dto.RequestLogin) (resp *dto.Token, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:Login")
	defer trace.Finish()

	user, err := u.psql.FindUserByUsername(ctx, req.Username)
	if err != nil {
		trace.SetError(err)

		return nil, errs.NewErrorWithCodeErr(err, errs.UserNotFound)
	}

	// match the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		return nil, errs.NewErrorWithCodeErr(err, errs.PasswordNotMatch)
	}

	// start the transaction
	err = u.psql.StartTransaction(ctx, func(ctx context.Context, sd dbc.SqlDbc) error {
		// repo := psql.New(sd, sd)

		// generate access token
		token, err := u.generateTokens(ctx, req.Channel, req.DeviceId, user)
		if err != nil {
			trace.SetError(err)

			return err
		}
		resp = token

		// insert user_device
		// err = repo.CreateUserDevice(ctx, &model.UserDevice{
		// 	UserId:    user.Id,
		// 	DeviceId:  req.DeviceId,
		// 	Channel:   req.Channel,
		// 	UserAgent: req.UserAgent,
		// 	FcmToken:  null.NewString(req.FcmToken),
		// })
		// if err != nil {
		// 	trace.SetError(err)

		// 	return errs.NewErrorWithCodeErr(err, errs.INSERT_DB_FAIL)
		// }

		return nil
	})

	if err != nil {
		return
	}

	return
}
