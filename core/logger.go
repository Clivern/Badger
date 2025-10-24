// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/clivern/yun/service"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// SetupLogging configures the logging system based on viper configuration
func SetupLogging() error {
	var writer io.Writer

	// Handle log file creation if not stdout
	if viper.GetString("app.log.output") != "stdout" {
		dir, _ := filepath.Split(viper.GetString("app.log.output"))

		// Create dir
		if !service.DirExists(dir) {
			if err := service.EnsureDir(dir, 0775); err != nil {
				return fmt.Errorf("directory [%s] creation failed: %w", dir, err)
			}
		}

		// Create log file if not exists
		if !service.FileExists(viper.GetString("app.log.output")) {
			f, err := os.Create(viper.GetString("app.log.output"))
			if err != nil {
				return fmt.Errorf("error while creating log file [%s]: %w", viper.GetString("app.log.output"), err)
			}
			f.Close()
		}

		f, err := os.OpenFile(
			viper.GetString("app.log.output"),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0775,
		)
		if err != nil {
			return fmt.Errorf("error opening log file: %w", err)
		}
		writer = f
	} else {
		writer = os.Stdout
	}

	// Set log format
	if viper.GetString("app.log.format") == "json" {
		log.Logger = zerolog.New(writer).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: writer}).With().Timestamp().Logger()
	}

	// Set log level
	level := strings.ToLower(viper.GetString("app.log.level"))

	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return nil
}
