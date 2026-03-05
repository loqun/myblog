package handlers

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/loqun/fiber-server/game"
)

type Handler struct {
	store *session.Store
	db    *sql.DB
}

func New(store *session.Store, db *sql.DB) *Handler {
	return &Handler{
		store: store,
		db:    db,
	}
}

func (h *Handler) Home(c *fiber.Ctx) error {
	type Project struct {
		Name         string
		Description  string
		Technologies []string
		URL          string
	}

	type ContactLink struct {
		Name string
		URL  string
	}

	data := fiber.Map{
		"Title":  "Arif Muftalib - Developer ",
		"Name":   "arif muftalib",
		"Bio":    "Full-stack developer passionate about building modern web applications and sharing knowledge through writing.",
		"About":  "I'm a software developer with a love for solving real world  problems using code. With experience in both frontend and backend development, I enjoy creating seamless user experiences and efficient server-side logic. Most of my work using Javascript/Typescript and PHP/Laravel for backend. Ocassionally I code in golang for learning purpose.",
		"Skills": []string{"Go", "JavaScript", "TypeScript", "React", "Node.js", "PostgreSQL", "Docker", "AWS"},
		"Projects": []Project{
			{
				Name:         "Fiber Blog Platform",
				Description:  "A modern blog platform built with Go Fiber, featuring real-time capabilities and clean design.",
				Technologies: []string{"Go", "Fiber", "SQLite", "HTML/CSS"},
				URL:          "/blog",
			},
			{
				Name:         "WebSocket Chat App",
				Description:  "Real-time chat application with WebSocket support for instant messaging.",
				Technologies: []string{"Go", "WebSockets", "Redis", "JavaScript"},
				URL:          "/ws",
			},
		},
		"ContactMessage": "Feel free to reach out if you'd like to collaborate or just say hello!",
		"ContactLinks": []ContactLink{
			{Name: "Email", URL: "mailto:arif@example.com"},
			{Name: "GitHub", URL: "https://github.com/arifkiddocare"},
			{Name: "Blog", URL: "/blog"},
		},
	}

	return c.Render("home", data)
}

func (h *Handler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *Handler) Html(c *fiber.Ctx) error {
	log.Println("Serving HTML file")
	return c.SendFile("static/index.html")
}

func (h *Handler) Game(c *fiber.Ctx) error {

	//check if there is game manager , if not create the manager and create the lobby
	gameManager := game.NewGameManager()
	availableGame := gameManager.GetAvailableGame()
	if availableGame == nil {
		// No available game, create a new one
		newGameID := "game1" // Generate a unique game ID as needed
		availableGame = gameManager.CreateNewGame(newGameID)
		log.Println("Created new game with ID:", newGameID)
	} else {
		log.Println("Found available game")
	}

	return c.JSON(availableGame)

}

func (h *Handler) ListAllGames(c *fiber.Ctx) error {

	//check if there is game manager , if not create the manager and create the lobby
	gameManager := game.NewGameManager()
	availableGame := gameManager.GetAvailableGame()

	// This is a placeholder implementation
	return c.JSON(availableGame)
}

func (h *Handler) Session(c *fiber.Ctx) error {

	//show the session id
	session, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not get session")
	}

	return c.JSON(fiber.Map{"session": session.ID()})

}

func (h *Handler) TestCsrf(c *fiber.Ctx) error {

	return c.SendString("CSRF token is valid!")

}

func (h *Handler) LoginForm(c *fiber.Ctx) error {
	// render generic login page; default action posts to /login
	return c.Render("login", fiber.Map{
		"Title":  "Login",
		"Action": "/login",
		"csrf":   c.Locals("csrf"),
	})
}

// AdminLoginForm renders the same login template but configures the form to
// submit back to /backoffice/login. This allows an "admin" login URL to exist
// while reusing the existing authentication handler.
func (h *Handler) AdminLoginForm(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title":  "Admin Login",
		"Action": "/backoffice/login",
		"csrf":   c.Locals("csrf"),
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) ProcessLogin(c *fiber.Ctx) error {

	// Authentication logic using environment variables
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	log.Printf("Admin credentials - Username: %s, Password: %s\n", adminUser, adminPass)
	
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if adminUser == "" {
		adminUser = "admin"
	}
	if adminPass == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Admin credentials not configured"})
	}

	if req.Username == adminUser && req.Password == adminPass {
		session, err := h.store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create session"})
		}
		session.Set("user", req.Username)
		if err := session.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save session"})
		}
		return c.JSON(fiber.Map{"message": "Login successful"})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})

}
