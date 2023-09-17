package services

import (
	"tokeon-test-task/internal/config"
	"tokeon-test-task/internal/services/device"
	"tokeon-test-task/pkg/log"
)

type Services struct {
	deviceService *device.Service
}

func New(logger log.Logger, config *config.Config) (*Services, error) {
	return &Services{
		deviceService: device.New(),
	}, nil
}

func (s *Services) Device() *device.Service {
	return s.deviceService
}
