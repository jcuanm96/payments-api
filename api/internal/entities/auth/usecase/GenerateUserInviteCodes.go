package service

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/utils"
)

func (svc *usecase) GenerateUserInviteCodes(ctx context.Context, exec utils.Executable, userID int) error {
	const numInvitesPerUser = 3
	codes := []string{}
	for i := 0; i < numInvitesPerUser; i++ {
		shouldInsertCode := false
		codeRes, generateCodeErr := svc.GenerateGoatInviteCode(ctx, shouldInsertCode)
		if generateCodeErr != nil {
			return generateCodeErr
		}
		codes = append(codes, codeRes.GoatInviteCode)
	}
	return svc.repo.InsertUserInviteCodes(ctx, exec, userID, codes)
}
