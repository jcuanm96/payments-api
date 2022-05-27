package controller

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/appconfig"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func decodeSignupSMS(c *fiber.Ctx) (interface{}, error) {
	var p request.SignupSMS

	c.BodyParser(&p)

	if !p.Phone.Validate() {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid phone number.",
			"Please enter a valid phone number.",
		)
	}

	firstName, firstNameErr := user.CheckNameIsValid(p.FirstName)
	if firstNameErr != nil {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Please enter a valid first name: %v.", firstNameErr),
			fmt.Sprintf("Invalid first name: %v", firstNameErr),
		)
	}

	lastName, lastNameErr := user.CheckNameIsValid(p.LastName)
	if lastNameErr != nil {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Please enter a valid last name: %v.", lastNameErr),
			fmt.Sprintf("Invalid last name: %v", lastNameErr),
		)
	}

	if p.Email != "" {
		p.Email = strings.ToLower(p.Email)
		_, parseErr := mail.ParseAddress(p.Email)
		if parseErr != nil {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				"Please enter a valid email.",
				fmt.Sprintf("Invalid email: %v", parseErr),
			)
		}
	}

	if utils.IsEmptyValue(p.Code) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid verification code.",
			"code is required.",
		)
	}

	p.FirstName = firstName
	p.LastName = lastName

	return p, nil
}

func decodeSignupSMSGoat(c *fiber.Ctx) (interface{}, error) {
	var p request.SignupSMSGoat

	c.BodyParser(&p)

	regex := regexp.MustCompile(constants.LINKSUFFIX_REGEX)
	regexMatch := regex.MatchString(p.Username)
	if !regexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf(constants.LINKSUFFIX_REGEX_ERR, "Username"),
			fmt.Sprintf("username %s did not match regex %s", p.Username, constants.LINKSUFFIX_REGEX),
		)
	}

	if p.Email != "" {
		p.Email = strings.ToLower(p.Email)
		_, parseErr := mail.ParseAddress(p.Email)
		if parseErr != nil {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				"Please enter a valid email.",
				fmt.Sprintf("Invalid email: %v", parseErr),
			)
		}
	}

	if !p.Phone.Validate() {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid phone number.",
			"Please enter a valid phone number.",
		)
	}

	firstName, firstNameErr := user.CheckNameIsValid(p.FirstName)
	if firstNameErr != nil {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Please enter a valid first name: %v.", firstNameErr),
			fmt.Sprintf("Invalid first name: %v", firstNameErr),
		)
	}

	lastName, lastNameErr := user.CheckNameIsValid(p.LastName)
	if lastNameErr != nil {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid last name: %v.",
			fmt.Sprintf("Invalid last name: %v", lastNameErr),
		)
	}

	p.FirstName = firstName
	p.LastName = lastName

	if utils.IsEmptyValue(p.InviteCode) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid creator invite code.",
			"inviteCode is required.",
		)
	}

	if utils.IsEmptyValue(p.Code) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid verification code.",
			"code is required.",
		)
	}

	return p, nil
}

func decodeCheckPhone(c *fiber.Ctx) (interface{}, error) {
	var q request.Phone

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.CountryCode) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid country code.",
			"countryCode is required.",
		)
	}

	if utils.IsEmptyValue(q.Number) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid phone number.",
			"phoneNumber is required.",
		)
	}

	return q, nil
}

func decodeRefreshToken(c *fiber.Ctx) (interface{}, error) {
	var q request.RefreshToken

	c.BodyParser(&q)

	if utils.IsEmptyValue(q.RefreshToken) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"refreshToken is required.",
		)
	}

	return q, nil
}

func decodeSignOut(c *fiber.Ctx) (interface{}, error) {
	var q request.SignOut

	c.BodyParser(&q)

	if utils.IsEmptyValue(q.RefreshToken) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"refreshToken is required.",
		)
	}

	return q, nil
}

func decodeCheckUsername(c *fiber.Ctx) (interface{}, error) {
	var q request.Username

	c.QueryParser(&q)

	regex := regexp.MustCompile(constants.LINKSUFFIX_REGEX)
	regexMatch := regex.MatchString(q.Username)
	if !regexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf(constants.LINKSUFFIX_REGEX_ERR, "Username"),
			fmt.Sprintf("username %s did not match regex %s", q.Username, constants.LINKSUFFIX_REGEX),
		)
	}

	return q, nil
}

func decodeCheckEmail(c *fiber.Ctx) (interface{}, error) {
	var q request.Email

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.Email) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid email.",
			"email is required.",
		)
	} else {
		q.Email = strings.ToLower(q.Email)
		_, parseErr := mail.ParseAddress(q.Email)
		if parseErr != nil {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				"Please enter a valid email.",
				fmt.Sprintf("Invalid email: %v", parseErr),
			)
		}
	}

	return q, nil
}

func decodeCheckGoatInviteCode(c *fiber.Ctx) (interface{}, error) {
	var q request.GoatInviteCode

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.Code) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid invite code.",
			"code is required.",
		)
	}

	return q, nil
}

func decodeVerifySMS(c *fiber.Ctx) (interface{}, error) {
	var q request.VerifySMS

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.CountryCode) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid country code.",
			"countryCode is required.",
		)
	}

	if utils.IsEmptyValue(q.Number) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid phone number.",
			"phoneNumber is required.",
		)
	}

	return q, nil
}

func decodeSignInSMS(c *fiber.Ctx) (interface{}, error) {
	var p request.SignInSMS

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.CountryCode) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid country code.",
			"countryCode is required.",
		)
	}

	if utils.IsEmptyValue(p.Number) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid phone number.",
			"phoneNumber is required.",
		)
	}

	if utils.IsEmptyValue(p.Code) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid verification code.",
			"code is required.",
		)
	}

	return p, nil
}

func decodeGenerateGoatInviteCode(c *fiber.Ctx) (interface{}, error) {
	var q request.GenerateGoatInviteCode
	inviteSecret := c.Get("Authorization", "")
	if inviteSecret == "" {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Missing Secret.",
		)
	} else if inviteSecret != appconfig.Config.Auth.GoatInviteCodeSecret {
		return nil, httperr.New(
			401,
			http.StatusUnauthorized,
			"Incorrect Secret.",
		)
	}
	return q, nil
}
