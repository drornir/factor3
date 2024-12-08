package example

import (
	"sync"

	"github.com/drornir/factor3/pkg/config"
)

var (
	globalConfig       Config
	globalConfigLoader *config.Loader
	globalConfigLock   sync.RWMutex
)

func Global() Config {
	globalConfigLock.RLock()
	defer globalConfigLock.RUnlock()
	return globalConfig
}

func SetGlobal(c Config) {
	globalConfigLock.Lock()
	defer globalConfigLock.Unlock()
	globalConfig = c
}
