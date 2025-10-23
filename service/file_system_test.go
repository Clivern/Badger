// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"testing"

	"github.com/clivern/badger/pkg"

	"github.com/franela/goblin"
)

// TestUnitFileSystem
func TestUnitFileSystem(t *testing.T) {
	g := goblin.Goblin(t)

	fs := NewFileSystem()

	g.Describe("#FileSystem", func() {
		g.It("It should satisfy test cases", func() {
			g.Assert(fs.FileExists(fmt.Sprintf("%s/.gitignore", pkg.GetBaseDir("cache")))).Equal(true)
			g.Assert(fs.FileExists(fmt.Sprintf("%s/not_found.md", pkg.GetBaseDir("cache")))).Equal(false)
			g.Assert(fs.DirExists(pkg.GetBaseDir("cache"))).Equal(true)
			g.Assert(fs.DirExists(fmt.Sprintf("%s/not_found", pkg.GetBaseDir("cache")))).Equal(false)
			g.Assert(fs.EnsureDir(fmt.Sprintf("%s/new", pkg.GetBaseDir("cache")), 775)).Equal(nil)
			g.Assert(fs.DirExists(fmt.Sprintf("%s/new", pkg.GetBaseDir("cache")))).Equal(true)
			g.Assert(fs.DeleteDir(fmt.Sprintf("%s/new", pkg.GetBaseDir("cache")))).Equal(nil)
		})
	})
}
