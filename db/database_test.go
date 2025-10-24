// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLiteConnection(t *testing.T) {
	// Create a temporary SQLite database
	tmpFile := "/tmp/test_yun.db"
	defer os.Remove(tmpFile)

	config := Config{
		Driver:     "sqlite",
		DataSource: tmpFile,
	}

	conn, err := NewConnection(config)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// Test ping
	err = conn.Ping()
	assert.NoError(t, err)

	// Test close
	err = conn.Close()
	assert.NoError(t, err)
}

func TestUnsupportedDriver(t *testing.T) {
	config := Config{
		Driver: "postgresql",
	}

	conn, err := NewConnection(config)
	assert.Error(t, err)
	assert.Nil(t, conn)
	assert.Contains(t, err.Error(), "unsupported database driver")
}
