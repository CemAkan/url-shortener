package health

import (
	"context"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"time"
)

const (
	HealthCheckInterval = 10 * time.Second
)

// StartWatchdog monitors DB & Redis health and cancels ctx when threshold exceeded
func StartWatchdog(ctx context.Context) {
	ticker := time.NewTicker(HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			infrastructure.Log.Infof("Watchdog context cancelled, stopping health checks")
			return

		case <-ticker.C:
			//logger.Info("Watchdog tick: running health checks")

			SetDBStatus(checkDBHealth(ctx))
			SetRedisStatus(checkRedisHealth(ctx))
			SetEmailStatus(checkEmailHealth())

		}
	}
}

// checkDBHealth checks database health
func checkDBHealth(ctx context.Context) bool {
	healthy := true

	// db health check
	sqlDB, err := infrastructure.DB.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		infrastructure.Log.WithError(err).Errorf("Database healthcheck failed")
		healthy = false
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
		infrastructure.Log.WithError(err).Errorf("Redis healthcheck failed")
		healthy = false
	}

	return healthy
}

// checkEmailHealth checks email service health
func checkEmailHealth() bool {
	healthy := true

	//email service health check
	conn, err := infrastructure.Mail.Dialer.Dial()
	if err != nil {
		infrastructure.Log.WithError(err).Errorf("Mail healthcheck failed")
		healthy = false
	}
	if err := conn.Close(); err != nil {
		infrastructure.Log.WithError(err).Errorf("Tcp socket close error during mail service healthcheck: %v", err.Error())
	}

	return healthy
}
