package controller

import (
	"fmt"
	"net/http"
	"strings"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func decodeGetSubscriptionPrices(c *fiber.Ctx) (interface{}, error) {
	var p request.GetSubscriptionPrices

	c.QueryParser(&p)

	if utils.IsEmptyValue(p.GoatUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when trying to get subscription prices.",
			"goatUserID is required",
		)
	}

	return p, nil
}

func decodeSubscribeCurrUserToGoat(c *fiber.Ctx) (interface{}, error) {
	var p request.SubscribeCurrUserToGoat

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.GoatUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when trying to subscribe.",
			"goatUserID is required",
		)
	}

	return p, nil
}

func decodeUnsubscribeCurrUserFromGoat(c *fiber.Ctx) (interface{}, error) {
	var p request.UnsubscribeCurrUserFromGoat

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.GoatUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when trying to get unsubscribe.",
			"goatUserID is required",
		)
	}

	return p, nil
}

func decodeGetUserSubscriptions(c *fiber.Ctx) (interface{}, error) {
	var q request.GetUserSubscriptions

	c.QueryParser(&q)

	if q.Limit == nil {
		defaultLimit := uint64(10)
		q.Limit = &defaultLimit
	}

	if q.CursorID == nil {
		defaultCursorID := 0
		q.CursorID = &defaultCursorID
	}

	return q, nil
}

func decodeUpdateSubscriptionPrice(c *fiber.Ctx) (interface{}, error) {
	var p request.UpdateSubscriptionPrice

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.Currency) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when updating creator subscription price.",
			"currency is required.",
		)
	}

	if !strings.EqualFold(p.Currency, constants.DEFAULT_CURRENCY) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Only USD currency supported at the moment.",
			"Only USD currency supported at the moment.",
		)
	}

	if utils.IsEmptyValue(p.PriceInSmallestDenom) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when updating creator subscription price.",
			"priceInSmallestDenom is required.",
		)
	}

	if p.PriceInSmallestDenom < constants.MIN_USER_SUBSCRIPTION_PRICE_USD_IN_SMALLEST_DENOM {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("The price you set is too small. Must be greater than or equal to %d", constants.MIN_USER_SUBSCRIPTION_PRICE_USD_IN_SMALLEST_DENOM),
			fmt.Sprintf("The price you set is too small. Must be greater than or equal to %d", constants.MIN_USER_SUBSCRIPTION_PRICE_USD_IN_SMALLEST_DENOM),
		)
	}

	if p.PriceInSmallestDenom > constants.MAX_USER_SUBSCRIPTION_PRICE_USD_IN_SMALLEST_DENOM {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("The price you set is too large. Must be less than or equal to %d", constants.MAX_USER_SUBSCRIPTION_PRICE_USD_IN_SMALLEST_DENOM),
			fmt.Sprintf("The price you set is too large. Must be less than or equal to %d", constants.MAX_USER_SUBSCRIPTION_PRICE_USD_IN_SMALLEST_DENOM),
		)
	}

	return p, nil
}

func decodeCheckUserSubscribedToGoat(c *fiber.Ctx) (interface{}, error) {
	var q request.CheckUserSubscribedToGoat

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.ProviderUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"providerUserID is required.",
		)
	}

	return q, nil
}

func decodeGetMyPaidGroupSubscriptions(c *fiber.Ctx) (interface{}, error) {
	var q request.GetMyPaidGroupSubscriptions

	c.QueryParser(&q)

	if q.Limit == 0 {
		q.Limit = 10
	}

	return q, nil
}

func decodeGetMyPaidGroupSubscription(c *fiber.Ctx) (interface{}, error) {
	var q request.GetMyPaidGroupSubscription

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.ChannelID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"channelID is required",
		)
	}

	return q, nil
}

func decodeBatchUnsubscribePaidGroup(c *fiber.Ctx) (interface{}, error) {
	var p cloudtasks.StripePaidGroupUnsubscribeTask

	c.BodyParser(&p)

	return p, nil
}
