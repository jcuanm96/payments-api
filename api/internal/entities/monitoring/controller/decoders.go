package controller

import (
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func decodeGetMinRequiredVersionByDevice(c *fiber.Ctx) (interface{}, error) {
	var q request.GetMinRequiredVersionByDevice

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.Platform) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"platform is required.",
		)
	}
	return q, nil
}

func decodeGetDashboard(c *fiber.Ctx) (interface{}, error) {
	dashboardSecret := c.Get("Authorization", "")
	if dashboardSecret == "" {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"Missing Secret.",
		)
	} else if dashboardSecret != appconfig.Config.Auth.DashboardSecret {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"Incorrect Secret.",
		)
	}
	return struct{}{}, nil
}
