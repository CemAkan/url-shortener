package system

import (
	"context"
	"github.com/CemAkan/url-shortener/pkg/infrastructure/logger"
	"os"
	"os/signal"
	"syscall"
)

func HandleSignals(cancelFunc context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	sig := <-sigs
	logger.Log.Infof("Received signal: %s", sig)

	cancelFunc()
}
