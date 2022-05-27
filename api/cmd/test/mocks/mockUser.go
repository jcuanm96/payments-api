package mocks

import (
	"context"
	"net/http"
	"strconv"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/jackc/pgx/v4"
)

type MockUserUsecase struct {
	repo   user.Repository
	useruc user.Usecase
}

func NewMockUser(repo user.Repository, useruc user.Usecase) user.Usecase {
	return &MockUserUsecase{repo: repo, useruc: useruc}
}

func (muc *MockUserUsecase) GetCurrentUser(ctx context.Context) (*response.User, error) {
	id, atoiErr := strconv.Atoi(ctx.Value("CURRENT_USER_UUID").(string))
	if atoiErr != nil {
		vlog.Fatalf(ctx, "Unable to convert Id in mock GetCurrent User: %s", atoiErr)
	}

	return muc.useruc.GetUserByID(ctx, id)
}

func (muc *MockUserUsecase) GetUserByID(ctx context.Context, id int) (*response.User, error) {
	return muc.useruc.GetUserByID(ctx, id)
}

func (muc *MockUserUsecase) IsContact(ctx context.Context, contactID int) (*response.IsContact, error) {
	id, atoiErr := strconv.Atoi(ctx.Value("CURRENT_USER_UUID").(string))
	if atoiErr != nil {
		vlog.Fatalf(ctx, "Unable to convert Id in mock IsContact User: %s", atoiErr)
	}

	isContact, isContactErr := muc.repo.IsContact(ctx, id, contactID)
	if isContactErr != nil {
		vlog.Errorf(ctx, "Error checking IsContact: %s", isContactErr)
		return nil, httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			constants.ErrSomethingWentWrong,
		)
	}

	return &response.IsContact{IsContact: isContact}, nil
}

// UNIMPLEMENTED MOCK METHODS
func (muc *MockUserUsecase) GetMyProfile(ctx context.Context) (*response.MyProfile, error) {
	return nil, nil
}

func (muc *MockUserUsecase) BlockUser(ctx context.Context, blockUserID int) error {
	return nil
}

func (muc *MockUserUsecase) UnblockUser(ctx context.Context, unblockUserID int) error {
	return nil
}

func (muc *MockUserUsecase) IsBlocked(ctx context.Context, userID int) (bool, error) {
	return false, nil
}

func (muc *MockUserUsecase) GetGoatProfile(ctx context.Context, req request.GetGoatProfile) (*response.GetGoatProfile, error) {
	return nil, nil
}

func (muc *MockUserUsecase) ListUsers(ctx context.Context, req request.ListUsers) (*response.ListUsers, error) {
	return nil, nil
}

func (muc *MockUserUsecase) GetUserByUUID(ctx context.Context, runnable utils.Runnable, uuid string) (*response.User, error) {
	return nil, nil
}
func (muc *MockUserUsecase) GetUserByEmail(ctx context.Context, runnable utils.Runnable, email string) (*response.User, error) {

	return nil, nil
}
func (muc *MockUserUsecase) GetUserByPhone(ctx context.Context, runnable utils.Runnable, phone string) (*response.User, error) {

	return nil, nil
}
func (muc *MockUserUsecase) GetUserByUsername(ctx context.Context, runnable utils.Runnable, username string) (*response.User, error) {

	return nil, nil
}
func (muc *MockUserUsecase) SoftDeleteUser(ctx context.Context, userID int) error {
	return nil
}
func (muc *MockUserUsecase) HardDeleteUser(ctx context.Context) error {
	return nil
}
func (muc *MockUserUsecase) UndeleteUser(ctx context.Context, userID int) error {
	return nil
}

func (muc *MockUserUsecase) CheckPhoneAlreadyExists(ctx context.Context, value string) (bool, error) {

	return false, nil
}
func (muc *MockUserUsecase) CheckEmailAlreadyExists(ctx context.Context, runnable utils.Runnable, value string) (*response.User, error) {
	return nil, nil
}
func (muc *MockUserUsecase) CheckUsernameAlreadyExists(ctx context.Context, runnable utils.Runnable, username string) (*response.User, error) {
	return nil, nil
}

func (muc *MockUserUsecase) CreateGoatInviteCode(ctx context.Context, code string) (string, error) {
	return "", nil
}
func (muc *MockUserUsecase) GetInviteCodeStatus(ctx context.Context, code string) (*int, bool, error) {

	return nil, false, nil
}
func (muc *MockUserUsecase) GetNextAvailableUsername(ctx context.Context, firstName string, lastName string) (*string, error) {

	return nil, nil
}
func (muc *MockUserUsecase) UseInviteCode(ctx context.Context, tx pgx.Tx, code string, userID int) error {
	return nil
}
func (muc *MockUserUsecase) UpsertStripeAccountIDByEmail(ctx context.Context, email string, stripeAccountID string) error {
	return nil
}
func (muc *MockUserUsecase) UpsertUser(ctx context.Context, tx pgx.Tx, user *response.User) (*response.User, error) {

	return nil, nil
}
func (muc *MockUserUsecase) UpdateUserProfileAvatar(ctx context.Context, userID int, profileAvatarFilename string) error {
	return nil
}
func (muc *MockUserUsecase) UpdateCurrentUser(ctx context.Context, req request.UpdateUser) (*response.User, error) {
	userID := 1
	email := "fake@email96.com"
	updateUserParams := request.UpdateUser{
		Email: &email,
	}
	updateErr := muc.repo.UpdateUser(ctx, muc.repo.MasterNode(), userID, updateUserParams)
	updatedUser := response.User{
		ID:    userID,
		Email: &email,
	}
	return &updatedUser, updateErr
}
func (muc *MockUserUsecase) UpsertGoatChatsPrice(ctx context.Context, tx pgx.Tx, priceInSmallestDenom int64, currency string, goatUserID int) error {
	return nil
}
func (muc *MockUserUsecase) UploadProfileAvatar(ctx context.Context, req request.UploadProfileAvatar) (*response.UploadProfileAvatar, error) {

	return nil, nil
}
func (muc *MockUserUsecase) UpsertBio(ctx context.Context, req request.UpsertBio) (*response.UpsertBio, error) {

	return nil, nil
}

func (muc *MockUserUsecase) UpdatePhone(ctx context.Context, req request.UpdatePhone) (*response.UpdatePhone, error) {
	return nil, nil
}

func (muc *MockUserUsecase) SoftDeleteCurrentUser(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func (muc *MockUserUsecase) GetGoatUsers(ctx context.Context, req request.GetGoatUsers) (*response.GetGoatUsers, error) {
	return nil, nil
}

func (muc *MockUserUsecase) GetUserByStripeAccountID(ctx context.Context, runnable utils.Runnable, stripeAccountID string) (*response.User, error) {
	return nil, nil
}

func (muc *MockUserUsecase) GetUserByStripeID(ctx context.Context, runnable utils.Runnable, stripeID string) (*response.User, error) {
	return nil, nil
}

func (muc *MockUserUsecase) ReportUser(ctx context.Context, req request.ReportUser) error {
	return nil
}

func (muc *MockUserUsecase) UpgradeUserToGoat(ctx context.Context, formattedNumber string, req request.SignupSMSGoat) (*response.AuthSuccess, error) {
	return nil, nil
}

func (muc *MockUserUsecase) GetInviteCodeStatuses(ctx context.Context) (*response.GetInviteCodeStatuses, error) {
	return nil, nil
}
