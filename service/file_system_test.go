// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"testing"

	"github.com/clivern/badger/pkg"
	"github.com/stretchr/testify/assert"
)

// TestUnitFileSystem tests the FileSystem service
func TestUnitFileSystem(t *testing.T) {
	fs := NewFileSystem()

	t.Run("FileExists should detect existing files", func(t *testing.T) {
		assert.True(t, fs.FileExists(fmt.Sprintf("%s/.gitignore", pkg.GetBaseDir("cache"))))
		assert.False(t, fs.FileExists(fmt.Sprintf("%s/not_found.md", pkg.GetBaseDir("cache"))))
	})

	t.Run("DirExists should detect existing directories", func(t *testing.T) {
		assert.True(t, fs.DirExists(pkg.GetBaseDir("cache")))
		assert.False(t, fs.DirExists(fmt.Sprintf("%s/not_found", pkg.GetBaseDir("cache"))))
	})

	t.Run("EnsureDir should create and delete directories", func(t *testing.T) {
		newDir := fmt.Sprintf("%s/new", pkg.GetBaseDir("cache"))

		assert.NoError(t, fs.EnsureDir(newDir, 0775))
		assert.True(t, fs.DirExists(newDir))
		assert.NoError(t, fs.DeleteDir(newDir))
	})
}
