import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const currentUser = computed(() => user.value)

  // Actions
  async function login(username, password) {
    loading.value = true
    error.value = null

    try {
      // For demo purposes, accept any non-empty credentials
      // In production, this should call your actual API endpoint
      if (username && password) {
        const demoToken = btoa(`${username}:${Date.now()}`)
        token.value = demoToken
        user.value = {
          username,
          email: `${username}@badger.local`,
          role: 'admin'
        }
        localStorage.setItem('token', demoToken)
        return true
      } else {
        throw new Error('Invalid credentials')
      }
    } catch (err) {
      error.value = err.message || 'Login failed'
      return false
    } finally {
      loading.value = false
    }
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
  }

  async function checkAuth() {
    if (token.value) {
      // In production, validate token with backend
      // For now, we'll just check if it exists
      try {
        const tokenData = atob(token.value).split(':')
        user.value = {
          username: tokenData[0],
          email: `${tokenData[0]}@badger.local`,
          role: 'admin'
        }
      } catch {
        logout()
      }
    }
  }

  // Initialize auth state
  checkAuth()

  return {
    user,
    token,
    loading,
    error,
    isAuthenticated,
    currentUser,
    login,
    logout,
    checkAuth
  }
})

