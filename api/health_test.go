// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUnitHealthEndpoint tests the health check endpoint
func TestUnitHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("HealthAction should return OK status", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/_health", nil)
		c.Request = req

		HealthAction(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"status":"ok"}`, strings.TrimSpace(w.Body.String()))
	})
}
