// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// RequestSizeLimit creates a middleware that limits the size of request bodies
// maxBytes specifies the maximum allowed size in bytes
func RequestSizeLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/api/v1/") {
				log.Info().Str("path", r.URL.Path).Msg("Skipping request size limit for non-API route")
				next.ServeHTTP(w, r)
				return
			}

			// Limit the request body size
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
