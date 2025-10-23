// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clivern/badger/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SetupLogging configures the logging system based on viper configuration
func SetupLogging() error {
	fys := service.NewFileSystem()

	// Handle log file creation if not stdout
	if viper.GetString("server.log.output") != "stdout" {
		dir, _ := filepath.Split(viper.GetString("server.log.output"))

		// Create dir
		if !fys.DirExists(dir) {
			if err := fys.EnsureDir(dir, 0775); err != nil {
				return fmt.Errorf("directory [%s] creation failed: %w", dir, err)
			}
		}

		// Create log file if not exists
		if !fys.FileExists(viper.GetString("server.log.output")) {
			f, err := os.Create(viper.GetString("server.log.output"))
			if err != nil {
				return fmt.Errorf("error while creating log file [%s]: %w", viper.GetString("server.log.output"), err)
			}
			f.Close()
		}
	}

	// Set output destination
	if viper.GetString("server.log.output") == "stdout" {
		log.SetOutput(os.Stdout)
		gin.DefaultWriter = os.Stdout
	} else {
		f, err := os.OpenFile(
			viper.GetString("server.log.output"),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0775,
		)
		if err != nil {
			return fmt.Errorf("error opening log file: %w", err)
		}
		log.SetOutput(f)
		gin.DefaultWriter = f
	}

	// Set log level
	lvl := strings.ToLower(viper.GetString("server.log.level"))
	level, err := log.ParseLevel(lvl)

	if err != nil {
		level = log.InfoLevel
	}

	log.SetLevel(level)

	// Set log format
	if viper.GetString("server.log.format") == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	return nil
}
