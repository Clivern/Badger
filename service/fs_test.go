// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"testing"

	"github.com/clivern/yun/pkg"
	"github.com/stretchr/testify/assert"
)

// TestUnitFileExists tests the FileExists function
func TestUnitFileExists(t *testing.T) {
	t.Run("should return true for existing files", func(t *testing.T) {
		assert.True(t, FileExists(fmt.Sprintf("%s/.gitignore", pkg.GetBaseDir("cache"))))
	})

	t.Run("should return false for non-existing files", func(t *testing.T) {
		assert.False(t, FileExists(fmt.Sprintf("%s/not_found.md", pkg.GetBaseDir("cache"))))
	})

	t.Run("should return false for directories", func(t *testing.T) {
		assert.False(t, FileExists(pkg.GetBaseDir("cache")))
	})
}

// TestUnitDirExists tests the DirExists function
func TestUnitDirExists(t *testing.T) {
	t.Run("should return true for existing directories", func(t *testing.T) {
		assert.True(t, DirExists(pkg.GetBaseDir("cache")))
	})

	t.Run("should return false for non-existing directories", func(t *testing.T) {
		assert.False(t, DirExists(fmt.Sprintf("%s/not_found", pkg.GetBaseDir("cache"))))
	})

	t.Run("should return false for files", func(t *testing.T) {
		assert.False(t, DirExists(fmt.Sprintf("%s/.gitignore", pkg.GetBaseDir("cache"))))
	})
}

// TestUnitEnsureDir tests the EnsureDir function
func TestUnitEnsureDir(t *testing.T) {
	t.Run("should create new directory", func(t *testing.T) {
		newDir := fmt.Sprintf("%s/test_new_dir", pkg.GetBaseDir("cache"))

		// Clean up if exists from previous test
		_ = DeleteDir(newDir)

		assert.NoError(t, EnsureDir(newDir, 0755))
		assert.True(t, DirExists(newDir))

		// Clean up
		assert.NoError(t, DeleteDir(newDir))
	})

	t.Run("should not error if directory already exists", func(t *testing.T) {
		existingDir := pkg.GetBaseDir("cache")
		assert.NoError(t, EnsureDir(existingDir, 0755))
		assert.True(t, DirExists(existingDir))
	})
}

// TestUnitDeleteDir tests the DeleteDir function
func TestUnitDeleteDir(t *testing.T) {
	t.Run("should delete existing directory", func(t *testing.T) {
		testDir := fmt.Sprintf("%s/test_delete_dir", pkg.GetBaseDir("cache"))

		// Create the directory first
		assert.NoError(t, EnsureDir(testDir, 0755))
		assert.True(t, DirExists(testDir))

		// Delete it
		assert.NoError(t, DeleteDir(testDir))
		assert.False(t, DirExists(testDir))
	})

	t.Run("should not error when deleting non-existing directory", func(t *testing.T) {
		nonExistingDir := fmt.Sprintf("%s/non_existing_dir", pkg.GetBaseDir("cache"))
		assert.NoError(t, DeleteDir(nonExistingDir))
	})
}
