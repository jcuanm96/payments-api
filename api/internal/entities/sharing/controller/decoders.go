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
)

func decodeGetBioData(c *fiber.Ctx) (interface{}, error) {
	var p request.GetBioData

	c.QueryParser(&p)

	regex := regexp.MustCompile(constants.LINKSUFFIX_REGEX)
	regexMatch := regex.MatchString(p.Username)
	if !regexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf(constants.LINKSUFFIX_REGEX_ERR, "Username"),
		)
	}

	return p, nil
}

func decodeGetLink(c *fiber.Ctx) (interface{}, error) {
	var p request.GetLink

	c.QueryParser(&p)

	regex := regexp.MustCompile(constants.LINKSUFFIX_REGEX)
	regexMatch := regex.MatchString(p.LinkSuffix)
	if !regexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf(constants.LINKSUFFIX_REGEX_ERR, "Link path"),
		)
	}

	return p, nil
}

func decodeNewMessageLink(c *fiber.Ctx) (interface{}, error) {
	var p request.NewMessageLink

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.ChannelID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"channelID is required.",
		)
	}

	if utils.IsEmptyValue(p.MessageID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"messageID is required.",
		)
	}

	return p, nil
}

func decodeGetMessageByLinkSuffix(c *fiber.Ctx) (interface{}, error) {
	var q request.GetMessageByLink

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.LinkSuffix) {
		return nil, httperr.New(
			404,
			http.StatusNotFound,
			"Sorry, that message doesn't seem to exist.",
			"linkSuffix is required.",
		)
	}

	return q, nil
}

func decodePublicGetMessageByLink(c *fiber.Ctx) (interface{}, error) {
	var q request.GetMessageByLink

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.LinkSuffix) {
		return nil, httperr.New(
			404,
			http.StatusNotFound,
			"Sorry, that message doesn't seem to exist.",
			"linkSuffix is required.",
		)
	}

	return q, nil
}

func decodeUpsertBioLinks(c *fiber.Ctx) (interface{}, error) {
	var p request.UpsertBioLinks

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.ThemeID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please select a theme.",
			"themeID is required.",
		)
	}

	if len(p.TextContents) != len(p.Links) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"There are not as many items in textContents as in links.",
		)
	}

	if len(p.TextContents) > constants.MAX_LINKS {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("There are too many items in textContents. Max: %d.", constants.MAX_LINKS),
		)
	}

	for i, url := range p.Links {
		if !utils.IsUrl(url) {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("%s is not a valid url.", url),
			)
		}
		if strings.Contains(p.TextContents[i], "|") || strings.Contains(url, "|") {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				"\"|\" is a forbidden character.",
				"\"|\" is a forbidden character.",
			)
		}
	}

	return p, nil
}

func decodeUpsertTheme(c *fiber.Ctx) (interface{}, error) {
	themeSecret := c.Get("Authorization", "")
	if themeSecret == "" {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
		)
	} else if themeSecret != appconfig.Config.Sharing.ThemeSecret {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
		)
	}
	var p request.UpsertTheme

	c.BodyParser(&p)

	return p, nil
}

func decodeGetThemes(c *fiber.Ctx) (interface{}, error) {
	var p request.GetThemes
	c.QueryParser(&p)

	if p.Limit == 0 {
		p.Limit = 10
	}

	return p, nil
}
