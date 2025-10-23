<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
    <div class="max-w-md w-full">
      <!-- Logo and Header -->
      <div class="text-center mb-12">
        <div class="flex justify-center mb-6">
          <img src="/logo.png" alt="Badger Logo" class="h-32 w-auto">
        </div>
      </div>

      <!-- Login Form -->
      <div class="bg-white border border-gray-200 p-8">
        <form class="space-y-6" @submit.prevent="handleLogin">
          <!-- Username Field -->
          <div>
            <label for="username" class="block text-xs text-gray-500 mb-2">
              USERNAME
            </label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              required
              class="input-field"
              placeholder="Enter your username"
              :disabled="loading"
            >
          </div>

          <!-- Password Field -->
          <div>
            <label for="password" class="block text-xs text-gray-500 mb-2">
              PASSWORD
            </label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              class="input-field"
              placeholder="Enter your password"
              :disabled="loading"
            >
          </div>

          <!-- Remember Me -->
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                id="remember-me"
                v-model="form.rememberMe"
                type="checkbox"
                class="h-4 w-4 text-gray-900 focus:ring-gray-900 border-gray-300"
              >
              <label for="remember-me" class="ml-2 block text-sm text-gray-600">
                Remember me
              </label>
            </div>

            <div class="text-sm">
              <a href="#" class="text-gray-600 hover:text-gray-900">
                Forgot password?
              </a>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="border border-red-200 bg-red-50 p-3">
            <p class="text-sm text-red-900">
              {{ error }}
            </p>
          </div>

          <!-- Submit Button -->
          <div>
            <button
              type="submit"
              class="w-full btn-primary py-3 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="loading"
            >
              <span v-if="!loading">Sign in</span>
              <span v-else class="flex items-center justify-center">
                <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Signing in...
              </span>
            </button>
          </div>
        </form>
      </div>

      <!-- Footer -->
      <p class="text-center text-xs text-gray-400 mt-8">
        © 2025 Badger
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  username: '',
  password: '',
  rememberMe: false
})

const loading = ref(false)
const error = ref(null)

const handleLogin = async () => {
  loading.value = true
  error.value = null

  try {
    const success = await authStore.login(form.username, form.password)
    if (success) {
      router.push('/')
    } else {
      error.value = authStore.error || 'Login failed. Please try again.'
    }
  } catch (err) {
    error.value = err.message || 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}
</script>

