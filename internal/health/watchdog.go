package health

import (
	"context"
	"github.com/CemAkan/url-shortener/internal/infrastructure/cache"
	"github.com/CemAkan/url-shortener/internal/infrastructure/db"
	"github.com/CemAkan/url-shortener/internal/infrastructure/mail"
	"github.com/CemAkan/url-shortener/internal/metrics"
	"github.com/CemAkan/url-shortener/pkg/logger"
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
			logger.Log.Infof("Watchdog context cancelled, stopping health checks")
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
	metrics.DBUp.Set(1)
	healthy := true

	// db health check
	sqlDB, err := db.DB.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		logger.Log.WithError(err).Errorf("Database healthcheck failed")
		healthy = false
		metrics.DBUp.Set(0)
	}

	return healthy
}

// checkRedisHealth checks redis health
func checkRedisHealth(ctx context.Context) bool {
	metrics.RedisUp.Set(1)
	healthy := true

	//redis health check
	rCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := cache.Redis.Ping(rCtx).Err(); err != nil {
		logger.Log.WithError(err).Errorf("Redis healthcheck failed")
		healthy = false
		metrics.RedisUp.Set(0)
	}

	return healthy
}

// checkEmailHealth checks email service health
func checkEmailHealth() bool {
	metrics.MailUp.Set(1)
	healthy := true

	//email service health check
	conn, err := mail.Mail.Dialer.Dial()
	if err != nil {
		logger.Log.WithError(err).Errorf("Mail healthcheck failed")
		healthy = false
		metrics.MailUp.Set(0)
	}
	if err := conn.Close(); err != nil {
		logger.Log.WithError(err).Errorf("Tcp socket close error during mail service healthcheck: %v", err.Error())
	}

	return healthy
}
