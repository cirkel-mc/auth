package usecase

import (
	"cirkel/auth/internal/domain/dto"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/cirkel-mc/goutils/constants"
	"github.com/cirkel-mc/goutils/errs"
	"github.com/cirkel-mc/goutils/logger"
	"github.com/cirkel-mc/goutils/tracer"
)

func (u *usecaseInstance) basicAuth(ctx context.Context, req *dto.RequestHeader, token string) (resp *dto.Auth, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:BasicAuth")
	defer trace.Finish()

	t, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed to decode token: %s", err)

		return nil, errs.NewErrorWithCodeErr(err, errs.InvalidAuth)
	}

	tokens := strings.Split(string(t), ":")
	if len(tokens) != 2 {
		err = fmt.Errorf("invalid username and password")
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		return nil, errs.NewErrorWithCodeErr(err, errs.InvalidAuth)
	}

	username := tokens[0]
	password := tokens[1]

	client, err := u.psql.FindClientByClientId(ctx, username)
	if err != nil {
		trace.SetError(err)
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewErrorWithCodeErr(err, errs.GeneralError)
		}

		return nil, errs.NewErrorWithCodeErr(err, errs.UserNotFound)
	}

	s := sha256.New()
	s.Write([]byte(password))
	hashedPassword := fmt.Sprintf("%x", s.Sum(nil))
	logger.Log.Debugf(ctx, "hashedPassword %s and client.ClientSecret %s", hashedPassword, client.ClientSecret)

	if client.ClientSecret != hashedPassword {
		err = fmt.Errorf("password incorrect")

		trace.SetError(err)
		logger.Log.Error(ctx, err)
		return nil, errs.NewErrorWithCodeErr(err, errs.BadRequest)
	}

	return &dto.Auth{
		Channel:   client.Channel,
		PublicKey: client.PublicKey,
		DeviceId:  req.DeviceId,
		RoleKey:   constants.ACLPublic,
	}, nil
}
