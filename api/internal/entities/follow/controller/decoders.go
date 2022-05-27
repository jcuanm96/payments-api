package controller

import (
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func decodeFollow(c *fiber.Ctx) (interface{}, error) {
	var p request.Follow
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.UserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when trying to follow creator.",
			"userID is required",
		)
	}

	return p, nil
}

func decodeUnfollow(c *fiber.Ctx) (interface{}, error) {
	var p request.Unfollow
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.UserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when trying to unfollow creator.",
			"userID is required",
		)
	}

	return p, nil
}

func decodeIsFollowing(c *fiber.Ctx) (interface{}, error) {
	var p request.IsFollowing

	c.QueryParser(&p)

	if utils.IsEmptyValue(p.UserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"userID is required",
		)
	}

	return p, nil
}

func decodeGetFollowedGoats(c *fiber.Ctx) (interface{}, error) {
	var p request.GetFollowedGoats

	c.QueryParser(&p)

	if p.Limit == 0 {
		p.Limit = 10
	}
	return p, nil
}
