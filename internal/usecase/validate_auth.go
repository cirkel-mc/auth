package usecase

import (
	"cirkel/auth/internal/domain/constant"
	"cirkel/auth/internal/domain/dto"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cirkel-mc/goutils/errs"
	"github.com/cirkel-mc/goutils/tracer"
)

func (u *usecaseInstance) ValidateAuth(ctx context.Context, req *dto.RequestHeader) (resp *dto.Auth, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:ValidateAuth")
	defer trace.Finish()

	auths := strings.Split(req.Authorization, " ")
	if len(auths) != 2 {
		err = fmt.Errorf("authorization invalid")

		return nil, errs.NewErrorWithCodeErr(err, errs.INVALID_AUTH)
	}

	switch auths[0] {
	case constant.Basic:
		resp, err = u.basicAuth(ctx, auths[1])
	case constant.Bearer:
	default:
		err = errs.NewError(fmt.Errorf("token type invalid"), http.StatusUnauthorized, 100, errs.UNAUTHORIZED.Message())
	}

	if err != nil {
		return
	}

	return
}
