package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

const generateGoatCodeErr = "Something went wrong when trying to generate an invite code."

func (svc *usecase) GenerateGoatInviteCode(ctx context.Context, insertCode bool) (*response.GenerateGoatInviteCode, error) {
	goatInviteCodeLen := 6
	goatInviteCode := ""
	for {
		goatInviteCodeb := make([]byte, goatInviteCodeLen/2)
		if _, genCodeErr := rand.Read(goatInviteCodeb); genCodeErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				generateGoatCodeErr,
				fmt.Sprintf("Could not generate creator invite code. Err: %v", genCodeErr),
			)
		}

		goatInviteCode = fmt.Sprintf("%X", goatInviteCodeb)

		_, wasFound, getCodeStatusErr := svc.user.GetInviteCodeStatus(ctx, goatInviteCode)

		if getCodeStatusErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				generateGoatCodeErr,
				fmt.Sprintf("Error occurred connecting to db when checking for creator invite code. Err: %v", getCodeStatusErr),
			)
		}

		if !wasFound {
			break
		}
	}

	if insertCode {
		_, insertInviteCodeErr := svc.user.CreateGoatInviteCode(ctx, goatInviteCode)

		if insertInviteCodeErr != nil {
			return nil, httperr.NewCtx(
				ctx,
				500,
				http.StatusInternalServerError,
				generateGoatCodeErr,
				fmt.Sprintf("Error occurred connecting to db when creating creator invite code. Err: %v", insertInviteCodeErr),
			)
		}
	}

	res := response.GenerateGoatInviteCode{
		GoatInviteCode: goatInviteCode,
	}

	return &res, nil
}
