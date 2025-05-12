package health

import (
	"context"
	"time"

	"github.com/CemAkan/url-shortener/config"
)

const (
	HealthCheckInterval = 10 * time.Second
	MaxFailures         = 3
)

func StartWatchdog(shutdownFunc func()) {
	config.Log.Info("Health watchdog started")

	failures := 0
	ticker := time.NewTicker(HealthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		healthy := true

		// DB Health
		sqlDB, err := config.DB.DB()
		if err != nil {
			config.Log.WithError(err).Error("DB handle retrieval failed")
			healthy = false
		} else if err = sqlDB.Ping(); err != nil {
			config.Log.WithError(err).Error("Database ping failed")
			healthy = false
		}

		// Redis Health
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		if err := config.Redis.Ping(ctx).Err(); err != nil {
			config.Log.WithError(err).Error("Redis ping failed")
			healthy = false
		}
		cancel()

		if healthy {
			failures = 0
			config.Log.Info("Healthcheck passed: DB & Redis OK")
			continue
		}

		failures++
		config.Log.Warnf("Healthcheck failed (%d/%d)", failures, MaxFailures)

		if failures >= MaxFailures {
			config.Log.Error("Healthcheck failure threshold reached. Triggering shutdown.")
			shutdownFunc()
			break
		}
	}
}
