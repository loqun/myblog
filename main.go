package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/loqun/fiber-server/config"
	"github.com/loqun/fiber-server/handlers"
	"github.com/loqun/fiber-server/helper"
	"github.com/loqun/fiber-server/middleware"
	"github.com/loqun/fiber-server/routes"
	_ "github.com/mattn/go-sqlite3"

	"log"
)

var ctx = context.Background()

func main() {

	// Load .env file (only for local development)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (normal in production)")
	}

	//connect to sqlite database
	database, err := sql.Open("sqlite3", "./myapp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	//execute the schema sql to create tables if not exists
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	// Test connection
	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}

	//redis env data (optional)
	redisHost := helper.GetEnv("REDIS_HOST", "")
	redisPort := helper.GetEnv("REDIS_PORT", "6379")
	redisPassword := helper.GetEnv("REDIS_PASSWORD", "")

	// Clean up redis host if it contains protocol or port
	if redisHost != "" {
		// Remove redis:// prefix if present
		redisHost = strings.TrimPrefix(redisHost, "redis://")
		// Remove port if included in host
		if idx := strings.LastIndex(redisHost, ":"); idx != -1 {
			redisHost = redisHost[:idx]
		}
	}

	log.Printf("Redis config - Host: %s, Port: %s, Password set: %v", redisHost, redisPort, redisPassword != "")

	var rdb *redis.Client
	if redisHost != "" {
		log.Println("Redis configuration found, connecting...")
		addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: redisPassword,
		})

		_, err = rdb.Ping(ctx).Result()
		if err != nil {
			log.Printf("Could not connect to Redis: %v", err)
			log.Println("Continuing without Redis...")
		} else {
			fmt.Println("Successfully connected to Redis!")
		}
	} else {
		log.Println("No REDIS_HOST set, skipping Redis connection")
	}

	app := config.Setup()
	app.Use(cors.New())
	app.Use(csrf.New())

	//serve images
	app.Static("/", "./static/images")

	//session store for the app
	store := session.New()

	//wesocket endpoint + middleware for websocket
	wsGroup := app.Group("/ws", middleware.WebSocketMiddleware())
	wsGroup.Get("/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println("WebSocket connection established")
		log.Println("WebSocket query received")
		log.Println("WebSocket session info") // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)

		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	//initialize handlers with session store
	h := handlers.New(store, database)

	//middleware stack for web endpoints
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// wire routes, passing session store for auth middleware
	routes.Setup(app, h, store)

	if rdb != nil {
		defer rdb.Close()
	}
	log.Fatal(app.Listen(":8000"))
}
