package controller

import (
	"context"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
)

func UpdateCurrentUser(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.UpdateUser)

	res, err := svc.UpdateCurrentUser(ctx, req)

	return res, err
}

func ListUsers(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.ListUsers)

	res, err := svc.ListUsers(ctx, req)

	return res, err
}

func UploadProfileAvatar(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.UploadProfileAvatar)

	res, err := svc.UploadProfileAvatar(ctx, req)

	return res, err
}

func UpsertBio(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.UpsertBio)

	res, err := svc.UpsertBio(ctx, req)

	return res, err
}

func UpdatePhone(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.UpdatePhone)

	res, err := svc.UpdatePhone(ctx, req)

	return res, err
}

func GetCurrentUser(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)

	res, err := svc.GetCurrentUser(ctx)
	return res, err
}

func GetMyProfile(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)

	res, err := svc.GetMyProfile(ctx)

	return res, err
}

func GetUserContactByID(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.GetUserContactByID)

	var wg sync.WaitGroup

	var user *response.User
	var getUserByIDErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		user, getUserByIDErr = svc.GetUserByID(ctx, req.UserID)
	}()

	var isContactRes *response.IsContact
	var isContactErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		isContactRes, isContactErr = svc.IsContact(ctx, req.UserID)
	}()

	var isBlocked bool
	var isBlockedErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		isBlocked, isBlockedErr = svc.IsBlocked(ctx, req.UserID)
	}()
	wg.Wait()

	if getUserByIDErr != nil {
		return nil, getUserByIDErr
	}

	if isContactErr != nil {
		return nil, isContactErr
	}

	if isBlockedErr != nil {
		return nil, isBlockedErr
	}

	res := response.GetUserContactByID{
		User:      user,
		IsContact: isContactRes.IsContact,
		IsBlocked: isBlocked,
	}

	return res, nil
}

func GetGoatProfile(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.GetGoatProfile)

	res, err := svc.GetGoatProfile(ctx, req)
	return res, err
}

func GetGoatUsers(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.GetGoatUsers)

	res, err := svc.GetGoatUsers(ctx, req)

	return res, err
}

func BlockUser(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.BlockUser)

	err := svc.BlockUser(ctx, req.UserID)
	return nil, err
}

func UnblockUser(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.UnblockUser)

	err := svc.UnblockUser(ctx, req.UserID)
	return nil, err
}

func ReportUser(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	req := incomeRequest.(request.ReportUser)

	err := svc.ReportUser(ctx, req)

	return nil, err
}

func HardDeleteUser(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)

	err := svc.HardDeleteUser(ctx)

	return nil, err
}

func UpgradeUserToGoat(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)
	currentUser, getCurrUserErr := svc.GetCurrentUser(ctx)
	if getCurrUserErr != nil {
		vlog.Errorf(ctx, "Error occurred when retrieving user from db. Err: %s", getCurrUserErr.Error())
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			"Something went wrong when trying to upgrade to creator",
		)
	} else if currentUser == nil {
		return nil, httperr.NewCtx(
			ctx,
			404,
			http.StatusInternalServerError,
			"Current user is nil for UpgradeUserToGoat",
		)
	} else if currentUser.Type == "GOAT" {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"You are already a creator.",
		)
	} else if currentUser.Email == nil {
		return nil, httperr.NewCtx(
			ctx,
			400,
			http.StatusBadRequest,
			"Please set your email in order to upgrade to a creator.",
		)
	}

	upgradeReq := incomeRequest.(request.UpgradeUserToGoat)

	signupReq := request.SignupSMSGoat{
		FirstName:  currentUser.FirstName,
		LastName:   currentUser.LastName,
		Username:   currentUser.Username,
		Email:      *currentUser.Email,
		InviteCode: upgradeReq.InviteCode,
		Phone: request.Phone{
			Number:      currentUser.Phonenumber,
			CountryCode: currentUser.CountryCode,
		},
	}

	_, upgradeErr := svc.UpgradeUserToGoat(ctx, currentUser.Phonenumber, signupReq)

	return nil, upgradeErr
}

func GetInviteCodeStatuses(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(user.Usecase)

	res, err := svc.GetInviteCodeStatuses(ctx)

	return res, err
}
