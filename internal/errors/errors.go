package errors

import (
	e "errors"
)

var ErrDeviceAlreadyRegistered = e.New("device already registered")
var ErrDeviceNotFound = e.New("device not found")
