package usecase

import (
	"comrades-mc/auth/internal/domain/dto"
	"context"
	"net/http"

	"github.com/comrades-mc/goutils/config/database/dbc"
	"github.com/comrades-mc/goutils/errs"
	"github.com/comrades-mc/goutils/logger"
	"github.com/comrades-mc/goutils/tracer"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecaseInstance) Login(ctx context.Context, req *dto.RequestLogin) (resp *dto.Token, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:Login")
	defer trace.Finish()

	user, err := u.psql.FindUserByUsername(ctx, req.Username)
	if err != nil {
		trace.SetError(err)

		return nil, errs.NewError(err, http.StatusNotFound, 2101, "User tidak ditemukan")
	}

	// match the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		return nil, errs.NewError(err, http.StatusBadRequest, 2102, "Username atau password tidak sesuai")
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
