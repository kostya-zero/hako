package server

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kostya-zero/hako/internal/config"
	"github.com/kostya-zero/hako/internal/store"
	"github.com/kostya-zero/hako/internal/utils"
	"github.com/shirou/gopsutil/v4/host"
)

func LoadSnapshot(s *store.Storage, cfg *config.Config) error {
	f, err := os.Open(cfg.SnapshotFile)
	if err != nil {
		return errors.New("cannot open snapshot file")
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)

	var snapshot map[string]map[string]string

	if err = decoder.Decode(&snapshot); err != nil {
		return errors.New("cannot decode snapshot")
	}

	s.Load(snapshot)

	return nil
}

func RunSnapshotService(c *config.Config, s *store.Storage, ctx context.Context) {
	ticker := time.NewTicker(time.Second * 30)
	select {
	case <-ticker.C:
		PerformSnapshot(c, s)
	case <-ctx.Done():
		utils.L.Info("Performing final snapshot save...")
		PerformSnapshot(c, s)
	}
}

func PerformSnapshot(c *config.Config, s *store.Storage) {
	if !s.IsDirty() {
		utils.L.Warn("No need to make a snapshot. Storage is not modified.")
		return
	}

	snapshot := s.MakeSnapshot()
	file, err := os.Create(c.SnapshotFile)
	if err != nil {
		utils.L.Errorf("Failed to create file: %s", err.Error())
		return
	}
	defer file.Close()

	encode := gob.NewEncoder(file)

	if err = encode.Encode(snapshot); err != nil {
		utils.L.Errorf("Failed to encode a snapshot: %s", err.Error())
		return
	}

	utils.L.Info("Successfully made a snapshot.")
}

func StartServer(cfg *config.Config) error {
	var wg sync.WaitGroup
	storage := store.NewStorage()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if cfg.SnapshotEnabled && cfg.SnapshotFile != "" {
		utils.L.Info("Snapshot enabled. Loading a snapshot.")
		err := LoadSnapshot(&storage, cfg)
		if err != nil {
			utils.L.Warnf("Failed to load snapshot: %s. Using empty database instead.", err.Error())
		}

		// Start snapshot service.
		wg.Go(func() {
			RunSnapshotService(cfg, &storage, ctx)
		})
	}

	utils.L.Infof("Starting Hako Server %s", utils.Version)

	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("Hako Database %s", utils.Version),
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

		utils.L.Info("New database created.", "name", name)

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

		utils.L.Info("Database deleted.", "name", name)

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

		utils.L.Info("New key added.", "database", DBName, "key", keyName)

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

		utils.L.Info("Key deleted.", "database", DBName, "key", keyName)

		return c.Status(200).JSON(fiber.Map{
			"ok": true,
		})
	})

	app.Get("/system/storage", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"snapshot_enabled": cfg.SnapshotEnabled,
			"count_dbs":        storage.CountDB(),
		})
	})

	app.Get("/system/software", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"version": utils.Version,
			"arch":    runtime.GOARCH,
			"os":      runtime.GOOS,
		})
	})

	app.Get("/system/machine", func(c *fiber.Ctx) error {
		system, _ := host.Info()

		return c.Status(200).JSON(fiber.Map{
			"os":               system.OS,
			"version":          system.KernelVersion,
			"numcpu":           runtime.NumCPU(),
			"arch":             system.KernelArch,
			"platform":         system.Platform,
			"platform_version": system.PlatformVersion,
		})
	})

	wg.Go(func() {
		if err := app.Listen(cfg.Address); err != nil {
			utils.L.Errorf("Server listen stopped: %s", err.Error())
		}
	})

	<-ctx.Done()
	utils.L.Info("Performing shutdown...")
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		utils.L.Errorf("Failed to shutdown server: %s", err.Error())
	}

	wg.Wait()
	return nil
}
