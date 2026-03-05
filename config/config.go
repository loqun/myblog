package config

import (
	"html/template"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func Setup() *fiber.App {
	engine := html.New("./static", ".html")
	engine.AddFunc("truncate", func(s string, length int) string {
		if len(s) <= length {
			return s
		}
		return s[:length]
	})
	engine.AddFunc("safeHTML", func(s string) template.HTML {
		// Only use for trusted, pre-sanitized content
		return template.HTML(s)
	})
	engine.AddFunc("currentYear", func() int {
		return time.Now().Year()
	})
	
	return fiber.New(fiber.Config{
		Views:         engine,
		Prefork:       false, // Disabled for cloud deployments
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber-Server",
		AppName:       "Dev App v1.0.1",
	})
}