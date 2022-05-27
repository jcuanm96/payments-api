package controller

import (
	"context"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/follow"
)

func Follow(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(follow.Usecase)
	req := incomeRequest.(request.Follow)

	res, err := svc.Follow(ctx, req.UserID)

	return res, err
}

func Unfollow(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(follow.Usecase)
	req := incomeRequest.(request.Unfollow)

	err := svc.Unfollow(ctx, req.UserID)

	return nil, err
}

func IsFollowing(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(follow.Usecase)
	req := incomeRequest.(request.IsFollowing)

	res, err := svc.IsFollowing(ctx, req.UserID)

	return res, err
}

func GetFollowedGoats(uc interface{}, ctx context.Context, incomeRequest interface{}) (interface{}, error) {
	svc := uc.(follow.Usecase)
	req := incomeRequest.(request.GetFollowedGoats)

	res, err := svc.GetFollowedGoats(ctx, req)

	return res, err
}
