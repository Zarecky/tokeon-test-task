package controllers

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SenderService interface {
	SendMessage(ctx context.Context, deviceID *uuid.UUID, text string) error
}

type Sender struct {
	validator     *validator.Validate
	senderService SenderService
}

func NewSender(validator *validator.Validate, senderService SenderService) *Sender {
	return &Sender{
		validator,
		senderService,
	}
}

type SendBodyDto struct {
	DeviceID *uuid.UUID `json:"device_id"`
	Text     string     `json:"text" validate:"required"`
}

// Send godoc
//
//	@Summary		send message to the devices
//	@Description	send message to the device with id in body or to lthe all devices if id is not provided in body
//	@Tags			sender
//	@Accept			json
//	@Param			body			body		SendBodyDto	true	"Data"
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/send [post]
func (ctl *Sender) Send() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(SendBodyDto)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if err := ctl.validator.Struct(*body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		innterCtx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
		defer cancel()

		if err := ctl.senderService.SendMessage(innterCtx, body.DeviceID, body.Text); err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
