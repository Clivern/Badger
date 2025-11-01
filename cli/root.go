// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var config string

var rootCmd = &cobra.Command{
	Use: "mut",
	Short: `🐗 A Fast, Secure MCP Gateway, Set up in Minutes.

If you have any suggestions, bug reports, or annoyances please report
them to our issue tracker at <https://github.com/clivern/mut/issues>`,
}

// Execute runs cmd tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
