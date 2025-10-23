// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package config

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/clivern/badger/api"
	"github.com/clivern/badger/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Setup creates and configures the HTTP server
func Setup() *gin.Engine {
	// Set Gin mode
	if viper.GetString("server.mode") == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Recover middleware
	r.Use(gin.Recovery())

	// Prometheus middleware for HTTP metrics
	r.Use(middleware.PrometheusMiddleware())

	// Custom logger middleware
	r.Use(middleware.Logger())

	// Request timeout middleware
	r.Use(func(c *gin.Context) {
		timeoutCtx, _ := c.Request.Context(), func() {}
		if viper.GetInt("server.timeout") > 0 {
			var cancelFn func()
			timeoutCtx, cancelFn = c.Request.Context(), func() {}
			_ = cancelFn
		}
		c.Request = c.Request.WithContext(timeoutCtx)
		c.Next()
	})

	// Routes
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	r.GET("/", api.HealthAction)
	r.GET("/_health", api.HealthAction)
	middleware
	// Metrics endpoint with basic auth
	r.GET("/_metrics", gin.BasicAuth(gin.Accounts{
		viper.GetString("server.metrics.username"): viper.GetString("server.metrics.secret"),
	}), gin.WrapH(promhttp.Handler()))

	return r
}

// Run starts the HTTP server with graceful shutdown support
func Run(r *gin.Engine) error {
	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("server.port"))),
		Handler: r,
	}

	// Channel to listen for errors from the server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.WithFields(log.Fields{
			"port": viper.GetInt("server.port"),
			"tls":  viper.GetBool("server.tls.status"),
		}).Info("Starting HTTP server")

		var err error
		if viper.GetBool("server.tls.status") {
			err = srv.ListenAndServeTLS(
				viper.GetString("server.tls.crt_path"),
				viper.GetString("server.tls.key_path"),
			)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.WithField("signal", sig.String()).Info("Received shutdown signal")

		// Create a deadline for graceful shutdown
		shutdownTimeout := 30 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		log.WithField("timeout", shutdownTimeout).Info("Gracefully shutting down server")

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		log.Info("Server shutdown complete")
	}

	return nil
}
