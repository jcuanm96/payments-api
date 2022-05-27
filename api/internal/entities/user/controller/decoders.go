package controller

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/user"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/codes"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func decodeUpdateCurrentUser(c *fiber.Ctx) (interface{}, error) {
	var p request.UpdateUser

	c.BodyParser(&p)

	if p.Username != nil {
		trimmed := strings.TrimSpace(*p.Username)
		p.Username = &trimmed
		regex := regexp.MustCompile(constants.LINKSUFFIX_REGEX)
		regexMatch := regex.MatchString(*p.Username)
		if !regexMatch {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				fmt.Sprintf(constants.LINKSUFFIX_REGEX_ERR, "Username"),
				fmt.Sprintf("username %s did not match regex %s", *p.Username, constants.LINKSUFFIX_REGEX),
			)
		}
	}

	if p.FirstName != nil {
		firstName, firstNameErr := user.CheckNameIsValid(*p.FirstName)
		if firstNameErr != nil {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				fmt.Sprintf("Please enter a valid first name: %v.", firstNameErr),
				fmt.Sprintf("Invalid first name: %v", firstNameErr),
			)
		}
		p.FirstName = &firstName
	}
	if p.LastName != nil {
		lastName, lastNameErr := user.CheckNameIsValid(*p.LastName)
		if lastNameErr != nil {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				fmt.Sprintf("Please enter a valid last name: %v", lastNameErr),
				fmt.Sprintf("Invalid last name: %v", lastNameErr),
			)
		}
		p.LastName = &lastName
	}

	p.Email = sanitizeEmail(p.Email)
	if p.Email != nil {
		_, parseErr := mail.ParseAddress(*p.Email)
		if parseErr != nil {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				"Please enter a valid email",
				fmt.Sprintf("Invalid email: %v", parseErr),
			)
		}
	}

	return p, nil
}

func sanitizeEmail(s *string) *string {
	if s == nil {
		return nil
	}
	res := strings.ToLower(*s)
	return &res
}

func decodeListUsers(c *fiber.Ctx) (interface{}, error) {
	var p request.ListUsers
	// We hardcode these values to always return the first
	// page of search results to the frontend.
	p.Page = 1
	p.Size = 10
	c.QueryParser(&p)

	const maxQueryLength = 50
	if len(p.Query) > maxQueryLength {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Query must be less than %d characters.", maxQueryLength),
			fmt.Sprintf("query must be less than %d characters.", maxQueryLength),
		)
	}

	const maxUserTypesLength = 50
	if len(p.UserTypes) > maxUserTypesLength {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("userTypes must be less than %d characters.", maxUserTypesLength),
		)
	}

	return p, nil
}

func decodeUploadProfileAvatar(c *fiber.Ctx) (interface{}, error) {
	var p request.UploadProfileAvatar
	file, err := c.FormFile("profileAvatar")
	if err != nil || file == nil {
		return nil, httperr.New(
			codes.Omit,
			http.StatusBadRequest,
			"Sorry, something went wrong on our end",
			"profileAvatar is required.",
		)
	}

	p.ProfileAvatar = file

	return p, nil
}

func decodeUpsertBio(c *fiber.Ctx) (interface{}, error) {
	var p request.UpsertBio

	c.BodyParser(&p)

	if len(p.TextContent) > constants.MAX_USER_BIO_LENGTH {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Bio must be less than %d characters.", constants.MAX_USER_BIO_LENGTH),
			fmt.Sprintf("Bio must be less than %d characters.", constants.MAX_USER_BIO_LENGTH),
		)
	}

	return p, nil
}

func decodeUpdatePhone(c *fiber.Ctx) (interface{}, error) {
	var p request.UpdatePhone

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
		if utils.IsEmptyValue(p.CountryCode) {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				"Please enter a valid verification code.",
				"code is required.",
			)
		}
	}

	return p, nil
}

func decodeGetUserContactByID(c *fiber.Ctx) (interface{}, error) {
	var p request.GetUserContactByID

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

func decodeGetGoatProfile(c *fiber.Ctx) (interface{}, error) {
	var p request.GetGoatProfile

	c.QueryParser(&p)

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

func decodeGetGoatUsers(c *fiber.Ctx) (interface{}, error) {
	var p request.GetGoatUsers

	c.QueryParser(&p)

	if p.Limit == 0 {
		p.Limit = 10
	}

	return p, nil
}

func decodeBlockUser(c *fiber.Ctx) (interface{}, error) {
	var p request.BlockUser
	c.BodyParser(&p)

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

func decodeUnblockUser(c *fiber.Ctx) (interface{}, error) {
	var p request.UnblockUser

	c.BodyParser(&p)

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

func decodeReportUser(c *fiber.Ctx) (interface{}, error) {
	var p request.ReportUser
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.ReportedUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"reportedUserID is required",
		)
	}

	if len(p.Description) > constants.MAX_REPORT_USER_DESCRIPTION_LENGTH {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Report description must be less than %d characters.", constants.MAX_REPORT_USER_DESCRIPTION_LENGTH),
			fmt.Sprintf("description must be less than %d characters.", constants.MAX_REPORT_USER_DESCRIPTION_LENGTH),
		)
	}

	return p, nil
}

func decodeUpgradeUserToGoat(c *fiber.Ctx) (interface{}, error) {
	var p request.UpgradeUserToGoat
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.InviteCode) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please enter a valid invite code.",
			"inviteCode is required",
		)
	}

	return p, nil
}
