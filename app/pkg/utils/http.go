package utils

import (
	"io"
	"operation-service/pkg/logging"
)

func CloseRequestBody(logger *logging.Logger, body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logger.Fatalf("Error closing request body: %v", err)
	}
}
