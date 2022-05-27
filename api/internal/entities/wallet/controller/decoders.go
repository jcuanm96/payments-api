package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func decodeGetTransactions(c *fiber.Ctx) (interface{}, error) {
	var p request.GetTransactions

	c.QueryParser(&p)

	if p.Limit == 0 {
		defaultLimit := int64(10)
		p.Limit = defaultLimit
	}

	return p, nil
}

func decodeGetPendingTransactions(c *fiber.Ctx) (interface{}, error) {
	var p request.GetTransactions

	if p.Limit == 0 {
		defaultLimit := int64(10)
		p.Limit = defaultLimit
	}

	return p, nil
}

func decodeDefaultPaymentMethod(c *fiber.Ctx) (interface{}, error) {
	var p request.DefaultPaymentMethod

	c.BodyParser(&p)

	numberRegex := regexp.MustCompile(constants.DEFAULT_PAYMENT_METHOD_NUMBER_REGEX)
	numberRegexMatch := numberRegex.MatchString(p.Number)
	if !numberRegexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.DEFAULT_PAYMENT_METHOD_NUMBER_REGEX_ERR,
			fmt.Sprintf("number %s did not match regex %s", p.Number, constants.DEFAULT_PAYMENT_METHOD_NUMBER_REGEX),
		)
	}

	expMonthRegex := regexp.MustCompile(constants.DEFAULT_PAYMENT_METHOD_EXP_MONTH_REGEX)
	expMonthRegexMatch := expMonthRegex.MatchString(p.ExpMonth)
	if !expMonthRegexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.DEFAULT_PAYMENT_METHOD_EXP_MONTH_REGEX_ERR,
			fmt.Sprintf("expMonth %s did not match regex %s", p.ExpMonth, constants.DEFAULT_PAYMENT_METHOD_EXP_MONTH_REGEX),
		)
	}

	expYearRegex := regexp.MustCompile(constants.DEFAULT_PAYMENT_METHOD_EXP_YEAR_REGEX)
	expYearRegexMatch := expYearRegex.MatchString(p.ExpYear)
	if !expYearRegexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.DEFAULT_PAYMENT_METHOD_EXP_YEAR_REGEX_ERR,
			fmt.Sprintf("expYear %s did not match regex %s", p.ExpYear, constants.DEFAULT_PAYMENT_METHOD_EXP_MONTH_REGEX),
		)
	}

	cvcRegex := regexp.MustCompile(constants.DEFAULT_PAYMENT_METHOD_CVC_REGEX)
	cvcRegexMatch := cvcRegex.MatchString(p.CVC)
	if !cvcRegexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.DEFAULT_PAYMENT_METHOD_CVC_REGEX_ERR,
			fmt.Sprintf("cvc %s did not match regex %s", p.CVC, constants.DEFAULT_PAYMENT_METHOD_CVC_REGEX),
		)
	}

	return p, nil
}

func decodeMakeChatPaymentIntent(c *fiber.Ctx) (interface{}, error) {
	var p request.MakeChatPaymentIntent

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.ProviderUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when making transaction.",
			"providerUserID is required",
		)
	}

	return p, nil
}

func decodeConfirmPaymentIntent(c *fiber.Ctx) (interface{}, error) {
	var p request.ConfirmPaymentIntent

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.CustomerUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when making transaction.",
			"customerUserID is required",
		)
	}

	return p, nil
}

func decodeUpsertGoatChatsPrice(c *fiber.Ctx) (interface{}, error) {
	var p request.UpsertGoatChatsPrice

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.Currency) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when updating creator chats price.",
			"currency is required.",
		)
	}

	if utils.IsEmptyValue(p.PriceInSmallestDenom) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when updating creator chats price.",
			"priceInSmallestDenom is required.",
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

	if p.PriceInSmallestDenom < constants.MIN_GOAT_CHAT_PRICE_USD_IN_SMALLEST_DENOM {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Minimum price is $%.2f.", float64(constants.MIN_GOAT_CHAT_PRICE_USD_IN_SMALLEST_DENOM)/100),
		)
	}

	if p.PriceInSmallestDenom > constants.MAX_GOAT_CHAT_PRICE_USD_IN_SMALLEST_DENOM {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Maximum price is $%.2f.", float64(constants.MAX_GOAT_CHAT_PRICE_USD_IN_SMALLEST_DENOM)/100),
		)
	}

	return p, nil
}

func decodeUpsertBank(c *fiber.Ctx) (interface{}, error) {
	var p request.UpsertBank

	c.BodyParser(&p)

	if v := validate.Struct(p); !v.Validate() {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			v.Errors.One(),
		)
	}

	p.BankName = strings.ToUpper(strings.TrimSpace(p.BankName))
	p.AccountType = strings.ToUpper(strings.TrimSpace(p.AccountType))
	p.AccountNumber = strings.ReplaceAll(p.AccountNumber, " ", "")
	p.RoutingNumber = strings.ReplaceAll(p.RoutingNumber, " ", "")
	p.AccountHolderName = strings.TrimSpace(p.AccountHolderName)
	p.AccountHolderType = strings.ToUpper(strings.TrimSpace(p.AccountHolderType))
	if p.AccountHolderType != constants.PERSONAL && p.AccountHolderType != constants.BUSINESS {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Account holder type must be one of %s, %s", constants.PERSONAL, constants.BUSINESS),
		)
	}

	p.BillingAddress.Street1 = strings.TrimSpace(p.BillingAddress.Street1)
	p.BillingAddress.Street2 = strings.TrimSpace(p.BillingAddress.Street2)
	p.BillingAddress.City = strings.TrimSpace(p.BillingAddress.City)
	p.BillingAddress.State = strings.TrimSpace(p.BillingAddress.State)
	p.BillingAddress.PostalCode = strings.TrimSpace(p.BillingAddress.PostalCode)
	p.BillingAddress.Country = strings.TrimSpace(p.BillingAddress.Country)

	return p, nil
}

func decodeGetBalance(c *fiber.Ctx) (interface{}, error) {
	var q request.GetBalance

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.Currency) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when getting balance. Please try again.",
			"currency is required.",
		)
	}

	// Only supporting USD for now
	if !strings.EqualFold(q.Currency, constants.DEFAULT_CURRENCY) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Only USD currency supported at the moment.",
			"Only USD currency supported at the moment.",
		)
	}

	return q, nil
}

func decodeGetGoatChatPrice(c *fiber.Ctx) (interface{}, error) {
	var q request.GetGoatChatPrice

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.GoatUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when getting creator chat price.",
			"goatUserID is required.",
		)
	}

	return q, nil
}

func decodeMarkProviderAsPaid(c *fiber.Ctx) (interface{}, error) {
	var p request.MarkProviderAsPaid
	c.BodyParser(&p)

	payoutSecret := c.Get("Authorization", "")
	if payoutSecret == "" {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Missing Secret.",
		)
	} else if payoutSecret != appconfig.Config.Auth.PayoutSecret {
		return nil, httperr.New(
			401,
			http.StatusUnauthorized,
			"Incorrect Secret.",
		)
	}

	if utils.IsEmptyValue(p.ProviderID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when marking provider as paid.",
			"providerID is required.",
		)
	}

	if utils.IsEmptyValue(p.AmountPaid) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when marking provider as paid.",
			"amountPaid is required.",
		)
	}

	if utils.IsEmptyValue(p.PayPeriodID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when marking provider as paid.",
			"payPeriodID is required.",
		)
	}

	return p, nil
}

func decodeListUnpaidProviders(c *fiber.Ctx) (interface{}, error) {
	var q request.ListUnpaidProviders

	c.QueryParser(&q)

	payoutSecret := c.Get("Authorization", "")
	if payoutSecret == "" {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Missing Secret.",
		)
	} else if payoutSecret != appconfig.Config.Auth.PayoutSecret {
		return nil, httperr.New(
			401,
			http.StatusUnauthorized,
			"Incorrect Secret.",
		)
	}

	return q, nil
}

func decodeGetPayoutPeriods(c *fiber.Ctx) (interface{}, error) {
	var q request.GetPayoutPeriods

	c.QueryParser(&q)

	payoutSecret := c.Get("Authorization", "")
	if payoutSecret == "" {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Missing Secret.",
		)
	} else if payoutSecret != appconfig.Config.Auth.PayoutSecret {
		return nil, httperr.New(
			401,
			http.StatusUnauthorized,
			"Incorrect Secret.",
		)
	}

	if q.Limit == 0 {
		q.Limit = 10
	}

	return q, nil
}

func decodeListPayoutHistory(c *fiber.Ctx) (interface{}, error) {
	var q request.ListPayoutHistory

	c.QueryParser(&q)

	payoutSecret := c.Get("Authorization", "")
	if payoutSecret == "" {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Missing Secret.",
		)
	} else if payoutSecret != appconfig.Config.Auth.PayoutSecret {
		return nil, httperr.New(
			401,
			http.StatusUnauthorized,
			"Incorrect Secret.",
		)
	}

	if q.Limit == 0 {
		q.Limit = 10
	}

	return q, nil
}

func decodeGetPayPeriod(c *fiber.Ctx) (interface{}, error) {
	var q request.GetPayPeriod

	c.QueryParser(&q)

	payoutSecret := c.Get("Authorization", "")
	if payoutSecret == "" {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Missing Secret.",
		)
	} else if payoutSecret != appconfig.Config.Auth.PayoutSecret {
		return nil, httperr.New(
			401,
			http.StatusUnauthorized,
			"Incorrect Secret.",
		)
	}

	if utils.IsEmptyValue(q.Timestamp) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when getting pay period.",
			"timestamp is required.",
		)
	}

	return q, nil
}
