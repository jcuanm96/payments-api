package controller

import (
	"net/http"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func decodeCreateContact(c *fiber.Ctx) (interface{}, error) {
	var p request.CreateContact
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.ContactID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"contactID is required.",
		)
	}

	return p, nil
}

func decodeGetContacts(c *fiber.Ctx) (interface{}, error) {
	var p request.GetContacts

	c.QueryParser(&p)

	return p, nil
}

func decodeDeleteContact(c *fiber.Ctx) (interface{}, error) {
	var p request.DeleteContact

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.ContactID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"contactID is required.",
		)
	}

	return p, nil
}

func decodeIsContact(c *fiber.Ctx) (interface{}, error) {
	var p request.IsContact

	c.QueryParser(&p)

	if utils.IsEmptyValue(p.ContactID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"contactID is required.",
		)
	}

	return p, nil
}

func decodeUploadContacts(c *fiber.Ctx) (interface{}, error) {
	var p request.UploadContacts

	c.BodyParser(&p)

	return p, nil
}

func decodeBatchAddUsersToContacts(c *fiber.Ctx) (interface{}, error) {
	var p cloudtasks.AddUserContactsTask
	c.BodyParser(&p)

	return p, nil
}
