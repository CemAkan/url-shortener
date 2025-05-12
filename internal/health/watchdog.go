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

// StartWatchdog monitors DB & Redis health and cancels ctx when threshold exceeded
func StartWatchdog(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(HealthCheckInterval)
	defer ticker.Stop()

	failures := 0

	config.Log.Info("Health Watchdog started")

	for {
		select {
		case <-ctx.Done():
			config.Log.Info("Watchdog context cancelled, stopping health checks")
			return

		case <-ticker.C:
			if !checkHealth(ctx) {
				failures++
				config.Log.Warnf("Healthcheck failed (%d/%d)", failures, MaxFailures)
				if failures >= MaxFailures {
					config.Log.Error("Failure threshold reached. Triggering shutdown via context cancel")
					cancel()
					return
				}
			} else {
				failures = 0
				config.Log.Info("Healthcheck passed: DB & Redis OK")
			}
		}
	}
}

func checkHealth(ctx context.Context) bool {
	healthy := true

	// db health check
	sqlDB, err := config.DB.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		config.Log.WithError(err).Error("Database healthcheck failed")
		healthy = false
	}

	//redis health check
	rCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := config.Redis.Ping(rCtx).Err(); err != nil {
		config.Log.WithError(err).Error("Redis healthcheck failed")
		healthy = false
	}

	return healthy
}
