// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cli

import (
	cfg "github.com/clivern/badger/config"
	"github.com/clivern/badger/db"
	"github.com/clivern/badger/migration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Manage database migrations`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")

		if err := cfg.Load(configFile); err != nil {
			log.Fatal().Err(err).Msg("Failed to load configuration")
		}

		if err := cfg.SetupLogging(); err != nil {
			log.Fatal().Err(err).Msg("Failed to setup logging")
		}

		// Initialize database connection
		dbConfig := db.Config{
			Driver:          viper.GetString("database.driver"),
			Host:            viper.GetString("database.host"),
			Port:            viper.GetInt("database.port"),
			Username:        viper.GetString("database.username"),
			Password:        viper.GetString("database.password"),
			Database:        viper.GetString("database.name"),
			MaxOpenConns:    viper.GetInt("database.max_open_conns"),
			MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
			ConnMaxLifetime: viper.GetInt("database.conn_max_lifetime"),
			DataSource:      viper.GetString("database.datasource"),
		}

		conn, err := db.NewConnection(dbConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}
		defer conn.Close()

		// Create migration manager
		mgr := migration.NewManager(conn.DB, conn.Driver)

		// Register all migrations
		for _, m := range migration.GetAll() {
			mgr.Register(m)
		}

		// Run migrations
		if err := mgr.Up(); err != nil {
			log.Fatal().Err(err).Msg("Failed to run migrations")
		}

		log.Info().Msg("Migration completed successfully")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")

		if err := cfg.Load(configFile); err != nil {
			log.Fatal().Err(err).Msg("Failed to load configuration")
		}

		if err := cfg.SetupLogging(); err != nil {
			log.Fatal().Err(err).Msg("Failed to setup logging")
		}

		// Initialize database connection
		dbConfig := db.Config{
			Driver:          viper.GetString("database.driver"),
			Host:            viper.GetString("database.host"),
			Port:            viper.GetInt("database.port"),
			Username:        viper.GetString("database.username"),
			Password:        viper.GetString("database.password"),
			Database:        viper.GetString("database.name"),
			MaxOpenConns:    viper.GetInt("database.max_open_conns"),
			MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
			ConnMaxLifetime: viper.GetInt("database.conn_max_lifetime"),
			DataSource:      viper.GetString("database.datasource"),
		}

		conn, err := db.NewConnection(dbConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}
		defer conn.Close()

		// Create migration manager
		mgr := migration.NewManager(conn.DB, conn.Driver)

		// Register all migrations
		for _, m := range migration.GetAll() {
			mgr.Register(m)
		}

		// Roll back migration
		if err := mgr.Down(); err != nil {
			log.Fatal().Err(err).Msg("Failed to roll back migration")
		}

		log.Info().Msg("Rollback completed successfully")
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")

		if err := cfg.Load(configFile); err != nil {
			log.Fatal().Err(err).Msg("Failed to load configuration")
		}

		if err := cfg.SetupLogging(); err != nil {
			log.Fatal().Err(err).Msg("Failed to setup logging")
		}

		// Initialize database connection
		dbConfig := db.Config{
			Driver:          viper.GetString("database.driver"),
			Host:            viper.GetString("database.host"),
			Port:            viper.GetInt("database.port"),
			Username:        viper.GetString("database.username"),
			Password:        viper.GetString("database.password"),
			Database:        viper.GetString("database.name"),
			MaxOpenConns:    viper.GetInt("database.max_open_conns"),
			MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
			ConnMaxLifetime: viper.GetInt("database.conn_max_lifetime"),
			DataSource:      viper.GetString("database.datasource"),
		}

		conn, err := db.NewConnection(dbConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}
		defer conn.Close()

		// Create migration manager
		mgr := migration.NewManager(conn.DB, conn.Driver)

		// Register all migrations
		for _, m := range migration.GetAll() {
			mgr.Register(m)
		}

		// Show status
		if err := mgr.Status(); err != nil {
			log.Fatal().Err(err).Msg("Failed to get migration status")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)

	migrateUpCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	migrateUpCmd.MarkFlagRequired("config")
	migrateDownCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	migrateDownCmd.MarkFlagRequired("config")
	migrateStatusCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	migrateStatusCmd.MarkFlagRequired("config")
}
