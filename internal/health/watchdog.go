package health

import (
	"context"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"time"
)

const (
	HealthCheckInterval = 10 * time.Second
)

var (
	logger = infrastructure.SpecialLogger("health", "file")
)

// StartWatchdog monitors DB & Redis health and cancels ctx when threshold exceeded
func StartWatchdog(ctx context.Context) {
	ticker := time.NewTicker(HealthCheckInterval)
	defer ticker.Stop()

	logger.Info("Health Watchdog started")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Watchdog context cancelled, stopping health checks")
			return

		case <-ticker.C:
			//logger.Info("Watchdog tick: running health checks")

			SetDBStatus(checkDBHealth(ctx))
			SetRedisStatus(checkRedisHealth(ctx))

		}
	}
}

// checkDBHealth checks database health
func checkDBHealth(ctx context.Context) bool {
	healthy := true

	// db health check
	sqlDB, err := infrastructure.DB.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		infrastructure.Log.WithError(err).Error("Database healthcheck failed")
		healthy = false
		logger.Warn("Database down")
	}

	return healthy
}

// checkRedisHealth checks redis health
func checkRedisHealth(ctx context.Context) bool {
	healthy := true

	//redis health check
	rCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := infrastructure.Redis.Ping(rCtx).Err(); err != nil {
		infrastructure.Log.WithError(err).Error("Redis healthcheck failed")
		healthy = false
		logger.Warn("Redis down")
	}

	return healthy
}
