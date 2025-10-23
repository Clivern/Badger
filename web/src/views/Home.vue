<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <div class="flex-shrink-0 flex items-center">
              <img src="/logo.png" alt="Badger Logo" class="h-10 w-auto">
            </div>
            <div class="hidden md:ml-10 md:flex md:space-x-8">
              <a href="#" class="border-primary-500 text-gray-900 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                Dashboard
              </a>
              <a href="#" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                Gateway
              </a>
              <a href="#" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                Registry
              </a>
              <a href="#" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                Settings
              </a>
            </div>
          </div>
          <div class="flex items-center">
            <div class="flex items-center space-x-4">
              <span class="text-sm text-gray-700">{{ authStore.currentUser?.username }}</span>
              <button
                @click="handleLogout"
                class="btn-secondary text-sm"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- Page Header -->
      <div class="px-4 py-4 sm:px-0">
        <h1 class="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p class="mt-2 text-sm text-gray-600">
          Monitor and manage your MCP Gateway
        </p>
      </div>

      <!-- Stats Grid -->
      <div class="mt-8 px-4 sm:px-0">
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          <!-- Stat Card 1 -->
          <div class="card hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="flex-shrink-0 p-3 bg-green-100 rounded-lg">
                <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="ml-4 flex-1">
                <p class="text-sm font-medium text-gray-500">Status</p>
                <p class="text-2xl font-semibold text-gray-900">{{ stats.status }}</p>
              </div>
            </div>
          </div>

          <!-- Stat Card 2 -->
          <div class="card hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="flex-shrink-0 p-3 bg-blue-100 rounded-lg">
                <svg class="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <div class="ml-4 flex-1">
                <p class="text-sm font-medium text-gray-500">Total Requests</p>
                <p class="text-2xl font-semibold text-gray-900">{{ stats.totalRequests.toLocaleString() }}</p>
              </div>
            </div>
          </div>

          <!-- Stat Card 3 -->
          <div class="card hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="flex-shrink-0 p-3 bg-purple-100 rounded-lg">
                <svg class="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="ml-4 flex-1">
                <p class="text-sm font-medium text-gray-500">Avg Response</p>
                <p class="text-2xl font-semibold text-gray-900">{{ stats.avgResponse }}ms</p>
              </div>
            </div>
          </div>

          <!-- Stat Card 4 -->
          <div class="card hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="flex-shrink-0 p-3 bg-orange-100 rounded-lg">
                <svg class="w-8 h-8 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
                </svg>
              </div>
              <div class="ml-4 flex-1">
                <p class="text-sm font-medium text-gray-500">Active Gateways</p>
                <p class="text-2xl font-semibold text-gray-900">{{ stats.activeGateways }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Content Grid -->
      <div class="mt-8 px-4 sm:px-0">
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <!-- Recent Activity -->
          <div class="card">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h3>
            <div class="space-y-4">
              <div
                v-for="activity in recentActivities"
                :key="activity.id"
                class="flex items-start space-x-3 pb-4 border-b border-gray-200 last:border-0 last:pb-0"
              >
                <div class="flex-shrink-0">
                  <div :class="['w-2 h-2 rounded-full mt-2', activity.statusColor]" />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-gray-900">{{ activity.title }}</p>
                  <p class="text-sm text-gray-500">{{ activity.description }}</p>
                  <p class="text-xs text-gray-400 mt-1">{{ activity.time }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- System Info -->
          <div class="card">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">System Information</h3>
            <div class="space-y-3">
              <div class="flex justify-between py-2 border-b border-gray-200">
                <span class="text-sm text-gray-600">Version</span>
                <span class="text-sm font-medium text-gray-900">v0.4.0</span>
              </div>
              <div class="flex justify-between py-2 border-b border-gray-200">
                <span class="text-sm text-gray-600">Uptime</span>
                <span class="text-sm font-medium text-gray-900">{{ systemInfo.uptime }}</span>
              </div>
              <div class="flex justify-between py-2 border-b border-gray-200">
                <span class="text-sm text-gray-600">Memory Usage</span>
                <span class="text-sm font-medium text-gray-900">{{ systemInfo.memory }}</span>
              </div>
              <div class="flex justify-between py-2">
                <span class="text-sm text-gray-600">CPU Usage</span>
                <span class="text-sm font-medium text-gray-900">{{ systemInfo.cpu }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mt-8 px-4 sm:px-0">
        <div class="card">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
            <button class="btn-primary justify-center flex items-center space-x-2">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              <span>Add Gateway</span>
            </button>
            <button class="btn-secondary justify-center flex items-center space-x-2">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span>Configure</span>
            </button>
            <button class="btn-secondary justify-center flex items-center space-x-2">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
              <span>View Analytics</span>
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { healthAPI } from '@/api'

const router = useRouter()
const authStore = useAuthStore()

const stats = ref({
  status: 'Online',
  totalRequests: 42853,
  avgResponse: 145,
  activeGateways: 3
})

const recentActivities = ref([
  {
    id: 1,
    title: 'Gateway Connected',
    description: 'New MCP gateway registered successfully',
    time: '2 minutes ago',
    statusColor: 'bg-green-500'
  },
  {
    id: 2,
    title: 'Request Processed',
    description: 'API request completed in 120ms',
    time: '5 minutes ago',
    statusColor: 'bg-blue-500'
  },
  {
    id: 3,
    title: 'Configuration Updated',
    description: 'Gateway settings modified',
    time: '15 minutes ago',
    statusColor: 'bg-yellow-500'
  },
  {
    id: 4,
    title: 'Health Check',
    description: 'All systems operational',
    time: '30 minutes ago',
    statusColor: 'bg-green-500'
  }
])

const systemInfo = ref({
  uptime: '7 days, 14 hours',
  memory: '245 MB / 2 GB',
  cpu: '12%'
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(async () => {
  try {
    // Check API health
    await healthAPI.check()
  } catch (error) {
    console.error('Failed to check API health:', error)
  }
})
</script>

