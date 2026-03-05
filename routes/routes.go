package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/loqun/fiber-server/handlers"
	"github.com/loqun/fiber-server/middleware"
)

// Setup configures all application routes. We pass the session store
// so authentication middleware can validate user sessions.
func Setup(app *fiber.App, h *handlers.Handler, store *session.Store) {
	app.Get("/", h.Home)
	app.Get("/api/health", h.Health)
	app.Get("/api/html", h.Html)
	app.Get("/api/game", h.Game)
	app.Get("/api/list-games", h.ListAllGames)
	app.Get("/api/session", h.Session)

	// blog routes
	app.Get("/blog", h.Blog)

	// create endpoints require login
	auth := middleware.RequireAuth(store)
	app.Get("/blog/form", auth, h.BlogForm)
	app.Post("/store-blog", auth, h.StoreBlog)

	app.Get("/blog/:id", h.GetBlog)
	app.Get("/api/blogs", h.GetAllBlog)

	// regular login
	app.Get("/login", h.LoginForm)
	app.Post("/login", h.ProcessLogin)

	// alternate path for administrators
	app.Get("/backoffice/login", h.AdminLoginForm)
	app.Post("/backoffice/login", h.ProcessLogin)
}
