package utils

import (
	"io"
	"operation-service/pkg/logging"
	"time"
)

func DoWithAttempts(fn func() error, attempts int, delay time.Duration) error {
	var err error
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return err
}

func CloseBody(logger *logging.Logger, body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logger.Fatalf("Error closing body: %v", err)
	}
}
