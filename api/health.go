// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/yun/service"

	"github.com/rs/zerolog/log"
)

// HealthAction handles health check requests
func HealthAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Health check endpoint called")

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
