// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HealthAction Controller
func HealthAction(c *gin.Context) {

	log.Info(`Incoming Request to Health Action`)

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
