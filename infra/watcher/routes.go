package watcher

import (
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/poixeai/proxify/infra/config"
	"github.com/poixeai/proxify/infra/logger"
)

var ConfigValue atomic.Value // global config value

func WatchJSON(file string) {
	watcher, _ := fsnotify.NewWatcher()
	go func() {
		for event := range watcher.Events {
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				cfg, err := config.LoadRoutesConfig(file)
				if err != nil {
					logger.Errorf("[routes.json] file reload failed:", err)
					continue
				}
				ConfigValue.Store(cfg)
				logger.Info("[routes.json] file reloaded successfully.")
			}
		}
	}()
	watcher.Add(file)
}

func InitRoutesWatcher() error {
	cfg, err := config.LoadRoutesConfig("routes.json")
	if err != nil {
		return err
	}
	ConfigValue.Store(cfg)

	// start watcher
	WatchJSON("routes.json")

	return nil
}

func GetRoutes() *config.RoutesConfig {
	v := ConfigValue.Load()
	if v == nil {
		return &config.RoutesConfig{}
	}
	return v.(*config.RoutesConfig)
}
