// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

// TestUnitHealthEndpoint
func TestUnitHealthEndpoint(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("#HealthAction", func() {
		g.It("It should satisfy all provided test cases", func() {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest(http.MethodGet, "/_health", nil)
			c.Request = req

			HealthAction(c)

			g.Assert(w.Code).Equal(http.StatusOK)
			g.Assert(strings.TrimSpace(w.Body.String())).Equal(`{"status":"ok"}`)
		})
	})
}
