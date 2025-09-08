package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subscriptions-api/pkg/utils/logger"
)

type Stoppable interface {
	Shutdown(ctx context.Context) error
	Stop()
}

// WaitForSignals ожидает SIGINT/SIGTERM и корректно завершает сервисы
func WaitForSignals(timeout time.Duration, log *logger.Logger, services ...Stoppable) {
	// Контекст, который отменяется при сигнале
	ctx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	// Ждём сигнал
	<-ctx.Done()

	// Логируем причину отмены через ctx.Err()
	if err := ctx.Err(); err != nil {
		log.Info("Shutdown triggered by signal", "reason", err)
	} else {
		log.Info("Shutdown triggered")
	}

	// Контекст для graceful shutdown с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	// Останавливаем все сервисы
	for _, srv := range services {
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("Graceful shutdown failed, forcing stop", "error", err)
			srv.Stop()
		} else {
			log.Info("Service stopped gracefully")
		}
	}
}
