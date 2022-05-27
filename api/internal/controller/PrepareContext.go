package controller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

func (ctr *Ctr) PrepareContext(c *fiber.Ctx) (*context.Context, error) {
	ctx := context.Background()
	uid, ok := c.Locals("uid").(string)
	if ok && uid != "" {
		ctx = context.WithValue(ctx, "CURRENT_USER_UUID", uid)
	}

	uuid, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		return nil, uuidErr
	}

	requestID := uuid.String()
	ctx = context.WithValue(ctx, "REQUEST_ID", requestID)
	c.Locals("REQUEST_ID", requestID)

	return &ctx, nil
}
