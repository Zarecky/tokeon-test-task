package middleware

import (
	websocket_pkg "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const WebsocketAllowed = "ws_allowed"

func (m *Middleware) Websocket() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket_pkg.IsWebSocketUpgrade(ctx) {
			ctx.Locals(WebsocketAllowed, true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
