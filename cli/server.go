// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	conf "github.com/clivern/badger/config"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the badger backend server",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		if err := conf.Load(config); err != nil {
			panic(err.Error())
		}

		// Setup logging
		if err := conf.SetupLogging(); err != nil {
			panic(err.Error())
		}

		// Setup and configure the HTTP server
		r := conf.Setup()

		// Run the server
		if err := conf.Run(r); err != nil {
			panic(fmt.Sprintf("Server error: %s", err.Error()))
		}
	},
}

func init() {
	serverCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	serverCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(serverCmd)
}
