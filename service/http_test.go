// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	t.Run("WriteJSON with map data and 200 status", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		data := map[string]interface{}{
			"status":  "ok",
			"message": "test message",
		}

		// Act
		WriteJSON(w, http.StatusOK, data)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
		assert.Equal(t, "test message", response["message"])
	})

	t.Run("WriteJSON with struct data and 201 status", func(t *testing.T) {
		// Arrange
		type TestStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		w := httptest.NewRecorder()
		data := TestStruct{
			ID:   123,
			Name: "Test User",
		}

		// Act
		WriteJSON(w, http.StatusCreated, data)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response TestStruct
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 123, response.ID)
		assert.Equal(t, "Test User", response.Name)
	})

	t.Run("WriteJSON with error status codes", func(t *testing.T) {
		testCases := []struct {
			name       string
			statusCode int
			data       interface{}
		}{
			{
				name:       "400 Bad Request",
				statusCode: http.StatusBadRequest,
				data:       map[string]string{"error": "bad request"},
			},
			{
				name:       "404 Not Found",
				statusCode: http.StatusNotFound,
				data:       map[string]string{"error": "not found"},
			},
			{
				name:       "500 Internal Server Error",
				statusCode: http.StatusInternalServerError,
				data:       map[string]string{"error": "internal error"},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				w := httptest.NewRecorder()

				// Act
				WriteJSON(w, tc.statusCode, tc.data)

				// Assert
				assert.Equal(t, tc.statusCode, w.Code)
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			})
		}
	})

	t.Run("WriteJSON with empty data", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		data := map[string]interface{}{}

		// Act
		WriteJSON(w, http.StatusOK, data)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Empty(t, response)
	})

	t.Run("WriteJSON with nested data", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		data := map[string]interface{}{
			"user": map[string]interface{}{
				"id":   1,
				"name": "John Doe",
				"metadata": map[string]string{
					"role": "admin",
				},
			},
			"timestamp": "2023-01-01T00:00:00Z",
		}

		// Act
		WriteJSON(w, http.StatusOK, data)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		user, ok := response["user"].(map[string]interface{})
		assert.True(t, ok, "Expected user to be a map")
		assert.Equal(t, "John Doe", user["name"])
	})
}
