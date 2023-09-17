package controllers

import (
	"tokeon-test-task/pkg/log"

	"github.com/go-playground/validator"
)

type Controllers struct {
	common *Common
	device *Device
	sender *Sender
}

func New(log log.Logger, validator *validator.Validate, deviceService DeviceService, senderService SenderService) *Controllers {
	return &Controllers{
		common: NewCommon(),
		device: NewDevice(log, deviceService),
		sender: NewSender(validator, senderService),
	}
}

func (c *Controllers) Common() *Common {
	return c.common
}

func (c *Controllers) Device() *Device {
	return c.device
}

func (c *Controllers) Sender() *Sender {
	return c.sender
}
