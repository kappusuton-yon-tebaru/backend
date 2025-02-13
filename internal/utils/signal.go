package utils

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"go.uber.org/zap"
)

func WaitForTermination(logger *logger.Logger, cleanup func()) <-chan bool {
	done := make(chan bool, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		logger.Info("cleaning up")

		defer func() {
			if r := recover(); r != nil {
				logger.Fatal("panic occured while cleaning up", zap.Any("recovered", r), zap.Stack("stacktrace"))
			}
		}()

		cleanup()

		logger.Info("cleanup successfully")

		done <- true
	}()

	return done
}
