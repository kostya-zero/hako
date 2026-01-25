package main

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoadSnapshot(storage *Storage, config *Config) error {
	f, err := os.Open(config.SnapshotFile)
	if err != nil {
		return errors.New("cannot open snapshot file")
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)

	var snapshot map[string]map[string]string

	if err = decoder.Decode(&snapshot); err != nil {
		return errors.New("cannot decode snapshot")
	}

	storage.Load(snapshot)

	return nil
}

func RunSnapshotService(config *Config, storage *Storage, ctx context.Context) {
	ticker := time.NewTicker(time.Second * 30)
	select {
	case <-ticker.C:
		PerformSnapshot(config, storage)
	case <-ctx.Done():
		l.Info("Performing final snapshot save...")
		PerformSnapshot(config, storage)
	}
}

func PerformSnapshot(config *Config, storage *Storage) {
	if !storage.IsDirty() {
		return
	}

	snapshot := storage.MakeSnapshot()
	file, err := os.Create(config.SnapshotFile)
	if err != nil {
		l.Errorf("Failed to create file: %s", err.Error())
		return
	}
	defer file.Close()

	encode := gob.NewEncoder(file)

	if err = encode.Encode(snapshot); err != nil {
		l.Errorf("Failed to encode a snapshot: %s", err.Error())
		return
	}

	l.Info("Successfully made a snapshot.")
}

func StartServer(config *Config) error {
	storage := NewStorage()
	ctx := context.Background()

	if config.SnapshotsEnabled && config.SnapshotFile != "" {
		err := LoadSnapshot(&storage, config)
		if err != nil {
			l.Warnf("Failed to load snapshot: %s. Using empty database instead.", err.Error())
		}

		// Start snapshot service.
		go RunSnapshotService(config, &storage, ctx)
	}

	l.Info("Starting Hako Server %s", Version)

	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("Hako Database %s", Version),
	})

	// Required because of how fiber works with params
	safeCopy := func(value string) string {
		return string(append([]byte(nil), value...))
	}

	app.Get("/db", func(c *fiber.Ctx) error {
		dbs := storage.GetDBNames()
		return c.Status(200).JSON(fiber.Map{
			"dbs":   dbs,
			"count": len(dbs),
		})
	})

	app.Post("/db/:database", func(c *fiber.Ctx) error {
		name := safeCopy(c.Params("database"))
		err := storage.CreateDatabase(name)
		if err != nil {
			return c.Status(409).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		l.Info("New database created.", "name", name)

		return c.Status(201).JSON(fiber.Map{
			"name": name,
		})
	})

	app.Delete("/db/:database", func(c *fiber.Ctx) error {
		name := safeCopy(c.Params("database"))
		err := storage.DeleteDatabase(name)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		l.Info("Database deleted.", "name", name)

		return c.Status(204).JSON(fiber.Map{
			"ok": true,
		})
	})

	app.Get("/db/:database/keys", func(c *fiber.Ctx) error {
		databaseName := safeCopy(c.Params("database"))
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
			"count": len(keys),
		})
	})

	app.Get("/db/:database/kv/:key", func(c *fiber.Ctx) error {
		DBName := safeCopy(c.Params("database"))
		keyName := safeCopy(c.Params("key"))

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
		DBName := safeCopy(c.Params("database"))
		keyName := safeCopy(c.Params("key"))

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

		l.Info("New key added.", "database", DBName, "key", keyName)

		return c.Status(200).JSON(fiber.Map{
			"ok": true,
		})
	})

	app.Delete("/db/:database/kv/:key", func(c *fiber.Ctx) error {
		DBName := safeCopy(c.Params("database"))
		keyName := safeCopy(c.Params("key"))

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

		l.Info("Key deleted.", "database", DBName, "key", keyName)

		return c.Status(200).JSON(fiber.Map{
			"ok": true,
		})
	})

	app.Listen(":3000")

	return nil
}
