package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func StartServer(config *Config) error {
	// Using in-memory storage for now
	storage := NewStorage()

	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("Hako Database %s", Version),
	})

	app.Get("/db", func(c *fiber.Ctx) error {
		dbs := storage.GetDBNames()
		return c.Status(200).JSON(fiber.Map{
			"dbs":   dbs,
			"count": len(dbs),
		})
	})

	app.Post("/db/:database", func(c *fiber.Ctx) error {
		name := c.Params("database")
		err := storage.CreateDatabase(name)
		if err != nil {
			return c.Status(409).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"name": name,
		})
	})

	app.Delete("/db/:database", func(c *fiber.Ctx) error {
		err := storage.DeleteDatabase(c.Params("database"))
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(204).JSON(fiber.Map{
			"ok": true,
		})
	})

	app.Get("/db/:database/keys", func(c *fiber.Ctx) error {
		databaseName := c.Params("database")
		db, err := storage.GetDatabase(databaseName)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		keys := db.GetAllKeys()

		return c.JSON(fiber.Map{
			"db":    databaseName,
			"keys":  keys,
			"count": len(*keys),
		})
	})

	app.Get("/db/:database/kv/:key", func(c *fiber.Ctx) error {
		DBName := c.Params("database")
		keyName := c.Params("key")

		db, err := storage.GetDatabase(DBName)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		key := db.Get(keyName)
		if key == nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "key not found",
			})
		}

		return c.Status(200).SendString(*key)
	})

	app.Post("/db/:database/kv/:key", func(c *fiber.Ctx) error {
		DBName := c.Params("database")
		keyName := c.Params("key")

		body := string(c.BodyRaw())

		if body == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "body is empty",
			})
		}

		db, err := storage.GetDatabase(DBName)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = db.Set(keyName, body)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"ok": true,
		})
	})

	app.Delete("/db/:database/kv/:key", func(c *fiber.Ctx) error {
		DBName := c.Params("database")
		keyName := c.Params("key")

		db, err := storage.GetDatabase(DBName)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = db.Delete(keyName)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"ok": true,
		})
	})

	app.Listen(":3000")

	return nil
}
