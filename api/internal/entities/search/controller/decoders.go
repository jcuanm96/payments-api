package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func decodeSearchGlobal(c *fiber.Ctx) (interface{}, error) {
	var q request.SearchGlobal

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.Query) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"A non-empty query is required",
			"A non-empty query is required",
		)
	}

	q.Query = strings.ToLower(strings.TrimSpace(q.Query))

	return q, nil
}

func decodeSearchMention(c *fiber.Ctx) (interface{}, error) {
	var q request.SearchMention

	c.QueryParser(&q)

	if utils.IsEmptyValue(q.ChannelID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"channelID is required",
		)
	}

	if len(q.Query) > constants.MAX_LINK_SUFFIX_LENGTH {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Query is too long (max length %d).", constants.MAX_LINK_SUFFIX_LENGTH),
			fmt.Sprintf("Query is too long (max length %d).", constants.MAX_LINK_SUFFIX_LENGTH),
		)
	}

	q.Query = strings.ToLower(q.Query)

	return q, nil
}
