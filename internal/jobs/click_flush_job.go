package job

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"time"
)

func StartClickFlushJob(flusher *app.ClickFlusherService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	infrastructure.Log.Infof("Click Flusher Job started with interval: %s", interval)

	for {
		<-ticker.C
		flusher.FlushClicks()
	}
}
