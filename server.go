package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func StartServer(config *Config) error {
	// Using in-memory storage for now
	// storage := initStorage()

	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("Hako Database %s", Version),
	})

	app.Get("/:database", func(c *fiber.Ctx) error {
		fmt.Printf("%s\n", c.Params("database"))
		return nil
	})

	app.Get("/:database/:key", func(c *fiber.Ctx) error {
		fmt.Printf("%s\n", c.Params("database"))
		fmt.Printf("%s\n", c.Params("key"))
		return nil
	})

	app.Listen(":3000")

	return nil
}
