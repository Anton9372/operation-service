package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"operation-service/internal/config"
	"operation-service/internal/controller/http"
	"operation-service/internal/domain/service"
	"operation-service/internal/storage/postgres"
	"operation-service/pkg/logging"
	"operation-service/pkg/metric"
	"operation-service/pkg/postgresql"
	"operation-service/pkg/shutdown"
	"os"
	"syscall"
	"time"
)

func main() {
	logging.InitLogger()
	logger := logging.GetLogger()
	logger.Info("logger initialized")

	logger.Info("config initializing")
	cfg := config.GetConfig()

	logger.Info("router initializing")
	router := httprouter.New()

	metricHandler := metric.Handler{Logger: logger}
	metricHandler.Register(router)

	postgresClient, err := postgresql.NewClient(context.Background(), 5, *cfg)
	if err != nil {
		logger.Fatal(err)
	}

	categoryStorage := postgres.NewCategoryRepo(postgresClient, logger)
	categoryService := service.NewCategoryService(categoryStorage, logger)
	categoryHandler := controller.NewCategoryHandler(categoryService, logger)
	categoryHandler.Register(router)

	logger.Info("start application")
	start(router, logger, cfg)
}

func start(router http.Handler, logger *logging.Logger, cfg *config.Config) {
	logger.Infof("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	logger.Info("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
