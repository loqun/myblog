package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// RequireAuth ensures that a user is logged in before allowing access to
// subsequent handlers. It checks the session store for a "user" value and
// either proceeds or returns an error/redirect if not present.
func RequireAuth(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			// session retrieval failure
			return c.Status(fiber.StatusInternalServerError).SendString("session error")
		}

		user := sess.Get("user")
		if user == nil {
			// not authenticated; redirect to login page
			// if the request expects JSON we could return a 401, but for the
			// simple blog create flow redirecting to login makes sense.
			return c.Redirect("/login")
		}

		// make the user available for later handlers if needed
		c.Locals("user", user)
		return c.Next()
	}
}
