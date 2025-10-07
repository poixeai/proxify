package watcher

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/poixeai/proxify/infra/config"
	"github.com/poixeai/proxify/infra/logger"
)

var ConfigValue atomic.Value // global config value

func WatchJSON(file string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Errorf("failed to create fsnotify watcher: %v", err)
		return
	}

	go func() {
		for event := range watcher.Events {
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				cfg, err := config.LoadRoutesConfig(file)
				if err != nil {
					logger.Errorf("[routes.json] file reload failed:", err)
					continue
				}

				if err := validateRoutes(cfg); err != nil {
					logger.Errorf("[routes.json] validation failed: %v", err)
					continue
				}

				ConfigValue.Store(cfg)
				logger.Info("[routes.json] file reloaded successfully.")
			}
		}
	}()

	if err := watcher.Add(file); err != nil {
		logger.Warnf("watcher: file [%s] not found, skip watching", file)
	}
}

func InitRoutesWatcher() error {
	const file = "routes.json"

	cfg, err := config.LoadRoutesConfig(file)
	if err != nil {
		if os.IsNotExist(err) {
			// if file not found, load default config
			logger.Warnf("[routes.json] not found, loading default config.")
			cfg = &config.RoutesConfig{
				Routes: []config.Route{
					// default routes, add more
					{Path: "/openai", Target: "https://api.openai.com"},
				},
			}
		} else {
			logger.Errorf("failed to load routes config: %v", err)
			return err
		}
	} else {
		logger.Infof("[routes.json] loaded successfully (%d routes)", len(cfg.Routes))
	}

	// validate routes
	if err := validateRoutes(cfg); err != nil {
		logger.Errorf("route validation failed: %v", err)
		return err
	}

	ConfigValue.Store(cfg)
	WatchJSON(file)

	return nil
}

func GetRoutes() *config.RoutesConfig {
	v := ConfigValue.Load()
	if v == nil {
		return &config.RoutesConfig{}
	}
	return v.(*config.RoutesConfig)
}

func validateRoutes(cfg *config.RoutesConfig) error {
	seen := make(map[string]bool)
	for _, r := range cfg.Routes {
		path := r.Path

		// 1. check empty
		if path == "" {
			return errors.New("invalid route: empty path is not allowed")
		}

		// 2. check reserved
		if config.ReservedTopRoutes[path] {
			return fmt.Errorf("invalid route: path '%s' is reserved by system", path)
		}

		// 3. check duplicate
		if seen[path] {
			return fmt.Errorf("invalid route: duplicate path '%s'", path)
		}
		seen[path] = true
	}
	return nil
}
