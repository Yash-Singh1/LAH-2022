package main

import (
	"context"
	"lah-2022/backend/auth"
	"lah-2022/backend/predictions"
	"os"

	"lah-2022/backend/ent"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	client, err := ent.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		e.Logger.Fatalf("failed creating schema resources: %v", err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://mail.google.com"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Routes
	auth.RegisterRoutes(client, e.Group("/auth"))
	predictions.RegisterRoutes(client, e.Group("/predict"))

	// Start server
	e.Logger.Fatal(e.Start(":3500"))
}
