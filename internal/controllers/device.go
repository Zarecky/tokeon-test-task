package controllers

import (
	"context"
	"tokeon-test-task/pkg/log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DeviceService interface {
	Register(id uuid.UUID) error
	Get(id uuid.UUID) (<-chan string, error)
	Close(id uuid.UUID) error
}

type Device struct {
	log           log.Logger
	deviceService DeviceService
}

func NewDevice(log log.Logger, deviceService DeviceService) *Device {
	return &Device{
		log,
		deviceService,
	}
}

func (d *Device) websocketCfg() *websocket.Config {
	return &websocket.Config{
		RecoverHandler: func(conn *websocket.Conn) {
			if err := recover(); err != nil {
				conn.WriteJSON(fiber.Map{"error": "Internal Server Error"})
			}
		},
	}
}

// Connect godoc
//
//	@Summary		open connect via websocket
//	@Description	open connect via websocket
//	@Param			id			path		string		true	"Unique id of the connecting device"
//	@Tags			device
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/ws/{id} [get]
func (d *Device) Connect(ctx context.Context) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		mt := websocket.TextMessage

		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			if err := c.WriteMessage(mt, []byte("id is not valid uuid")); err != nil {
				d.log.Errorf("write: %v", err)
			}

			if err := c.Close(); err != nil {
				d.log.Errorf("close: %v", err)
			}

			return
		}

		if err := d.deviceService.Register(id); err != nil {
			if err := c.WriteMessage(mt, []byte(err.Error())); err != nil {
				d.log.Errorf("write: %v", err)
			}

			if err := c.Close(); err != nil {
				d.log.Errorf("close: %v", err)
			}

			return
		}

		received := make(chan []byte)

		go func(ctx context.Context, message chan []byte) {
			var (
				msg []byte
				err error
			)
			for {
				if mt, msg, err = c.ReadMessage(); err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						close(message)
					} else {
						d.log.Errorf("read: %v", err)
					}
					break
				}

				message <- msg
			}
		}(ctx, received)

		ch, err := d.deviceService.Get(id)
		if err != nil {
			d.log.Errorf("failed to get device channel: %v", err)
			return
		}

		for {
			select {
			case msg, ok := <-received:
				d.log.Infof("revieved message from device %s: %s", id, msg)
				if !ok {
					if err := d.deviceService.Close(id); err != nil {
						d.log.Errorf("write: %v", err)
					}
					return
				}
			case msg := <-ch:
				if err = c.WriteMessage(mt, []byte(msg)); err != nil {
					d.log.Errorf("write: %v", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}, *d.websocketCfg())
}
