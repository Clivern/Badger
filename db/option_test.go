// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a test database with the options table
func setupTestDB(t *testing.T) (*Connection, func()) {
	tmpFile := "/tmp/test_options_" + time.Now().Format("20060102150405") + ".db"

	config := Config{
		Driver:     "sqlite",
		DataSource: tmpFile,
	}

	conn, err := NewConnection(config)
	require.NoError(t, err, "Failed to create test database")

	// Create options table
	_, err = conn.DB.Exec(`
		CREATE TABLE options (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR(255) NOT NULL UNIQUE,
			value TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err, "Failed to create options table")

	cleanup := func() {
		conn.Close()
		os.Remove(tmpFile)
	}

	return conn, cleanup
}

func TestOptionRepository_Create(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewOptionRepository(conn.DB)

	t.Run("Create new option successfully", func(t *testing.T) {
		err := repo.Create("app_name", "Yun")
		assert.NoError(t, err)

		// Verify it was created
		opt, err := repo.Get("app_name")
		assert.NoError(t, err)
		assert.NotNil(t, opt)
		assert.Equal(t, "app_name", opt.Key)
		assert.Equal(t, "Yun", opt.Value)
	})

	t.Run("Create option with empty value", func(t *testing.T) {
		err := repo.Create("empty_key", "")
		assert.NoError(t, err)

		opt, err := repo.Get("empty_key")
		assert.NoError(t, err)
		assert.NotNil(t, opt)
		assert.Equal(t, "", opt.Value)
	})

	t.Run("Create duplicate key should fail", func(t *testing.T) {
		err := repo.Create("duplicate", "value1")
		assert.NoError(t, err)

		err = repo.Create("duplicate", "value2")
		assert.Error(t, err, "Should fail on duplicate key")
	})

	t.Run("Create option with long value", func(t *testing.T) {
		longValue := string(make([]byte, 10000))
		err := repo.Create("long_value", longValue)
		assert.NoError(t, err)

		opt, err := repo.Get("long_value")
		assert.NoError(t, err)
		assert.Equal(t, len(longValue), len(opt.Value))
	})
}

func TestOptionRepository_Get(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewOptionRepository(conn.DB)

	t.Run("Get existing option", func(t *testing.T) {
		// Create test data
		err := repo.Create("test_key", "test_value")
		require.NoError(t, err)

		// Get the option
		opt, err := repo.Get("test_key")
		assert.NoError(t, err)
		assert.NotNil(t, opt)
		assert.Equal(t, "test_key", opt.Key)
		assert.Equal(t, "test_value", opt.Value)
		assert.Greater(t, opt.ID, int64(0))
		assert.False(t, opt.CreatedAt.IsZero())
		assert.False(t, opt.UpdatedAt.IsZero())
	})

	t.Run("Get non-existent option", func(t *testing.T) {
		opt, err := repo.Get("non_existent")
		assert.NoError(t, err)
		assert.Nil(t, opt, "Should return nil for non-existent option")
	})

	t.Run("Get with special characters in key", func(t *testing.T) {
		specialKey := "key-with_special.chars"
		err := repo.Create(specialKey, "special_value")
		require.NoError(t, err)

		opt, err := repo.Get(specialKey)
		assert.NoError(t, err)
		assert.NotNil(t, opt)
		assert.Equal(t, specialKey, opt.Key)
	})
}

func TestOptionRepository_Update(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewOptionRepository(conn.DB)

	t.Run("Update existing option", func(t *testing.T) {
		// Create initial option
		err := repo.Create("version", "1.0.0")
		require.NoError(t, err)

		// Get initial timestamps
		opt1, err := repo.Get("version")
		require.NoError(t, err)
		initialUpdatedAt := opt1.UpdatedAt

		// Wait a bit to ensure timestamp changes
		time.Sleep(10 * time.Millisecond)

		// Update the option
		err = repo.Update("version", "2.0.0")
		assert.NoError(t, err)

		// Verify update
		opt2, err := repo.Get("version")
		assert.NoError(t, err)
		assert.Equal(t, "2.0.0", opt2.Value)
		assert.True(t, opt2.UpdatedAt.After(initialUpdatedAt), "UpdatedAt should be newer")
	})

	t.Run("Update non-existent option", func(t *testing.T) {
		// Update should not fail even if option doesn't exist (SQL UPDATE behavior)
		err := repo.Update("does_not_exist", "some_value")
		assert.NoError(t, err)

		// Verify it wasn't created
		opt, err := repo.Get("does_not_exist")
		assert.NoError(t, err)
		assert.Nil(t, opt)
	})

	t.Run("Update to empty value", func(t *testing.T) {
		err := repo.Create("clear_me", "initial_value")
		require.NoError(t, err)

		err = repo.Update("clear_me", "")
		assert.NoError(t, err)

		opt, err := repo.Get("clear_me")
		assert.NoError(t, err)
		assert.Equal(t, "", opt.Value)
	})

	t.Run("Update multiple times", func(t *testing.T) {
		err := repo.Create("counter", "1")
		require.NoError(t, err)

		for i := 2; i <= 5; i++ {
			err = repo.Update("counter", string(rune('0'+i)))
			assert.NoError(t, err)
		}

		opt, err := repo.Get("counter")
		assert.NoError(t, err)
		assert.Equal(t, "5", opt.Value)
	})
}

func TestOptionRepository_Delete(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewOptionRepository(conn.DB)

	t.Run("Delete existing option", func(t *testing.T) {
		// Create option
		err := repo.Create("to_delete", "value")
		require.NoError(t, err)

		// Verify it exists
		opt, err := repo.Get("to_delete")
		assert.NoError(t, err)
		assert.NotNil(t, opt)

		// Delete it
		err = repo.Delete("to_delete")
		assert.NoError(t, err)

		// Verify it's gone
		opt, err = repo.Get("to_delete")
		assert.NoError(t, err)
		assert.Nil(t, opt)
	})

	t.Run("Delete non-existent option", func(t *testing.T) {
		// Should not fail
		err := repo.Delete("does_not_exist")
		assert.NoError(t, err)
	})

	t.Run("Delete and recreate", func(t *testing.T) {
		key := "recreate_test"

		// Create
		err := repo.Create(key, "value1")
		require.NoError(t, err)

		// Delete
		err = repo.Delete(key)
		assert.NoError(t, err)

		// Recreate with different value
		err = repo.Create(key, "value2")
		assert.NoError(t, err)

		// Verify new value
		opt, err := repo.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, "value2", opt.Value)
	})
}

func TestOptionRepository_List(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewOptionRepository(conn.DB)

	t.Run("List empty table", func(t *testing.T) {
		options, err := repo.List()
		assert.NoError(t, err)
		assert.Empty(t, options)
	})

	t.Run("List multiple options", func(t *testing.T) {
		// Create test data
		testData := map[string]string{
			"app_name":    "Yun",
			"version":     "1.0.0",
			"environment": "development",
			"debug":       "true",
		}

		for key, value := range testData {
			err := repo.Create(key, value)
			require.NoError(t, err)
		}

		// List all options
		options, err := repo.List()
		assert.NoError(t, err)
		assert.Len(t, options, len(testData))

		// Verify all keys are present
		keys := make(map[string]bool)
		for _, opt := range options {
			keys[opt.Key] = true
			assert.Equal(t, testData[opt.Key], opt.Value)
			assert.Greater(t, opt.ID, int64(0))
		}

		for key := range testData {
			assert.True(t, keys[key], "Key %s should be in results", key)
		}
	})

	t.Run("List returns sorted by key", func(t *testing.T) {
		// Clear previous data
		conn.DB.Exec("DELETE FROM options")

		// Create options in random order
		keys := []string{"zebra", "apple", "mango", "banana"}
		for _, key := range keys {
			err := repo.Create(key, "value")
			require.NoError(t, err)
		}

		// List and verify order
		options, err := repo.List()
		assert.NoError(t, err)
		assert.Len(t, options, len(keys))

		// Verify alphabetical order
		assert.Equal(t, "apple", options[0].Key)
		assert.Equal(t, "banana", options[1].Key)
		assert.Equal(t, "mango", options[2].Key)
		assert.Equal(t, "zebra", options[3].Key)
	})

	t.Run("List after deletions", func(t *testing.T) {
		// Clear previous data
		conn.DB.Exec("DELETE FROM options")

		// Create 5 options
		for i := 1; i <= 5; i++ {
			err := repo.Create(string(rune('a'+i-1)), "value")
			require.NoError(t, err)
		}

		// Delete 2 options
		repo.Delete("b")
		repo.Delete("d")

		// List remaining
		options, err := repo.List()
		assert.NoError(t, err)
		assert.Len(t, options, 3)

		keys := []string{options[0].Key, options[1].Key, options[2].Key}
		assert.Contains(t, keys, "a")
		assert.Contains(t, keys, "c")
		assert.Contains(t, keys, "e")
	})
}
