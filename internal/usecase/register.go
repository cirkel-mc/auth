package usecase

import (
	"comrades-mc/auth/internal/domain/constant"
	"comrades-mc/auth/internal/domain/dto"
	"comrades-mc/auth/internal/domain/model"
	"comrades-mc/auth/internal/repository/psql"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/comrades-mc/goutils/config/database/dbc"
	"github.com/comrades-mc/goutils/env"
	"github.com/comrades-mc/goutils/errs"
	"github.com/comrades-mc/goutils/logger"
	"github.com/comrades-mc/goutils/tracer"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecaseInstance) Register(ctx context.Context, req *dto.RequestRegister) (resp *dto.Token, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:Register")
	defer trace.Finish()

	// get customer role
	role, err := u.psql.FindRoleById(ctx, env.GetInt("CUSTOMER_ROLE_ID"))
	if err != nil {
		trace.SetError(err)

		return nil, errs.NewError(err, http.StatusNotFound, 2009, errs.NotFound)
	}

	// use variable errUser just incase
	_, errUser := u.psql.FindUserByEmail(ctx, req.Email)
	if errUser == nil {
		errUser = fmt.Errorf("email already exists")
		trace.SetError(errUser)

		return nil, errs.NewError(errUser, http.StatusConflict, 2001, "E-mail sudah digunakan")
	}

	_, errUser = u.psql.FindUserByUsername(ctx, req.Username)
	if errUser == nil {
		errUser = fmt.Errorf("username already exists")
		trace.SetError(errUser)

		return nil, errs.NewError(errUser, http.StatusConflict, 2002, "Username sudah digunakan")
	}

	userSeq, err := u.psql.GetUserNextVal(ctx)
	if err != nil {
		trace.SetError(err)

		return nil, errs.NewError(err, http.StatusNotFound, 2003, errs.NotFound)
	}

	// generate password with bcrypt
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "generate password error: %s", err)

		return nil, errs.NewErrorWithCodeErr(err, errs.GENERAL_ERROR)
	}

	now := time.Now().In(u.tz)
	user := &model.User{
		BaseModel: model.BaseModel{
			Id:        userSeq,
			CreatedAt: now,
		},
		Username: req.Username,
		Email:    req.Email,
		Password: string(password),
		Status:   constant.NotYetVerified,
		Role:     role,
	}

	// start the transaction
	err = u.psql.StartTransaction(ctx, func(ctx context.Context, sd dbc.SqlDbc) error {
		repo := psql.New(sd, sd)

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

		// insert users
		err = repo.CreateUser(ctx, user)
		if err != nil {
			trace.SetError(err)

			return errs.NewErrorWithCodeErr(err, errs.INSERT_DB_FAIL)
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}
