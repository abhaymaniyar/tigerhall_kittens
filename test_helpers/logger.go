package test_helpers

import (
	"sync"
	"tigerhall_kittens/internal/logger"
)

var loggerInitOnce sync.Once

func InitializeLogger() {
	loggerInitOnce.Do(func() {
		logger.Init(logger.DEBUG)
	})
}
