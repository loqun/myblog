package middleware

import (
	"github.com/gofiber/contrib/websocket"
	"log"
	"github.com/gofiber/fiber/v2"
)
 
func WebSocketMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.

		log.Println("WebSocket Endpoint Hit")
		log.Println("WebSocket connection request received")
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)	
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

