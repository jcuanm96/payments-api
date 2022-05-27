package controller

import (
	"net/http"
	"reflect"

	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/codes"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) WrapError(err error, c *fiber.Ctx) error {
	if reflect.TypeOf(err).String() == "*httperr.HTTPErr" {
		httpErr := err.(*httperr.HTTPErr)
		vlog.EndpointError(c, httpErr.Status, httpErr.Detail)

		return httpErr.Send(c)
	}

	ctr.Lr.LogError(err, c.Request())

	return ErrParseBody.SetDetail(err).Send(c)
}

var (
	ErrInternal  = httperr.New(codes.Omit, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	ErrParseBody = httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse request")
)
