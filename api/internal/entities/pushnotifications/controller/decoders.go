package controller

import (
	"net/http"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	push "github.com/VamaSingapore/vama-api/internal/entities/pushnotifications"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func decodeUpdateFcmToken(c *fiber.Ctx) (interface{}, error) {
	var p request.UpdateFcmToken
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.Token) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"fcmToken is required",
		)
	}

	return p, nil
}

func decodeSetGoatPostNotifications(c *fiber.Ctx) (interface{}, error) {
	var p request.SetGoatPostNotifications
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.GoatID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"goatID is required",
		)
	}

	return p, nil
}

func decodeUpdateSetting(c *fiber.Ctx) (interface{}, error) {
	var p request.UpdatePushSetting
	c.BodyParser(&p)

	if _, ok := push.ValidSettingIDs[p.ID]; !ok {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"invalid id.",
		)
	}

	if !(p.NewSetting == constants.PUSH_NOTIFICATION_ON ||
		p.NewSetting == constants.PUSH_NOTIFICATION_OFF ||
		p.NewSetting == constants.PUSH_NOTIFICATION_UNSET) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"invalid newSetting value.",
		)
	}

	return p, nil
}
