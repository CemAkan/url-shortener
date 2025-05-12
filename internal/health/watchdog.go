package health

import (
	"context"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	infrastructure2 "github.com/CemAkan/url-shortener/pkg/infrastructure"
	"time"
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

	infrastructure2.Log.Info("Health Watchdog started")

	for {
		select {
		case <-ctx.Done():
			infrastructure2.Log.Info("Watchdog context cancelled, stopping health checks")
			return

		case <-ticker.C:
			if !checkHealth(ctx) {
				failures++
				infrastructure2.Log.Warnf("Healthcheck failed (%d/%d)", failures, MaxFailures)
				if failures >= MaxFailures {
					infrastructure2.Log.Error("Failure threshold reached. Triggering shutdown via context cancel")
					cancel()
					return
				}
			} else {
				failures = 0
				infrastructure2.Log.Info("Healthcheck passed: DB & Redis OK")
			}
		}
	}
}

func checkHealth(ctx context.Context) bool {
	healthy := true

	// db health check
	sqlDB, err := infrastructure.DB.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		infrastructure2.Log.WithError(err).Error("Database healthcheck failed")
		healthy = false
	}

	//redis health check
	rCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := infrastructure2.Redis.Ping(rCtx).Err(); err != nil {
		infrastructure2.Log.WithError(err).Error("Redis healthcheck failed")
		healthy = false
	}

	return healthy
}
