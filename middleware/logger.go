// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Logger creates a new logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log incoming request
		log.Printf("Incoming Request: %s %s", c.Request.Method, c.Request.URL.String())

		// Call next handler
		c.Next()

		// Log outgoing response
		log.Printf("Outgoing Response: %d %s", c.Writer.Status(), http.StatusText(c.Writer.Status()))
	}
}
