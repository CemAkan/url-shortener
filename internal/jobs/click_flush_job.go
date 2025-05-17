package job

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/pkg/infrastructure/logger"
	"time"
)

func StartClickFlushJob(flusher *app.ClickFlusherService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Log.Infof("Click Flusher Job started with interval: %s", interval)

	for {
		<-ticker.C
		flusher.FlushClicks()
	}
}
