<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <div class="flex-shrink-0 flex items-center">
              <img src="/logo.png" alt="Yun Logo" class="h-10 w-auto">
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
      <div class="px-4 py-6 sm:px-0">
        <h1 class="text-2xl font-medium text-gray-900">Dashboard</h1>
      </div>

      <!-- Stats Grid -->
      <div class="px-4 sm:px-0">
        <div class="grid grid-cols-1 gap-px bg-gray-200 sm:grid-cols-2 lg:grid-cols-4">
          <!-- Stat Card 1 -->
          <div class="bg-white p-6">
            <p class="text-xs text-gray-500 mb-2">STATUS</p>
            <p class="text-2xl font-light text-gray-900">{{ stats.status }}</p>
          </div>

          <!-- Stat Card 2 -->
          <div class="bg-white p-6">
            <p class="text-xs text-gray-500 mb-2">TOTAL REQUESTS</p>
            <p class="text-2xl font-light text-gray-900">{{ stats.totalRequests.toLocaleString() }}</p>
          </div>

          <!-- Stat Card 3 -->
          <div class="bg-white p-6">
            <p class="text-xs text-gray-500 mb-2">AVG RESPONSE</p>
            <p class="text-2xl font-light text-gray-900">{{ stats.avgResponse }}ms</p>
          </div>

          <!-- Stat Card 4 -->
          <div class="bg-white p-6">
            <p class="text-xs text-gray-500 mb-2">ACTIVE GATEWAYS</p>
            <p class="text-2xl font-light text-gray-900">{{ stats.activeGateways }}</p>
          </div>
        </div>
      </div>

      <!-- Content Grid -->
      <div class="mt-8 px-4 sm:px-0">
        <div class="grid grid-cols-1 gap-px bg-gray-200 lg:grid-cols-2">
          <!-- Recent Activity -->
          <div class="bg-white p-8">
            <h3 class="text-sm font-medium text-gray-900 mb-6">RECENT ACTIVITY</h3>
            <div class="space-y-5">
              <div
                v-for="activity in recentActivities"
                :key="activity.id"
                class="border-b border-gray-100 pb-5 last:border-0 last:pb-0"
              >
                <p class="text-sm text-gray-900">{{ activity.title }}</p>
                <p class="text-sm text-gray-500 mt-1">{{ activity.description }}</p>
                <p class="text-xs text-gray-400 mt-1">{{ activity.time }}</p>
              </div>
            </div>
          </div>

          <!-- System Info -->
          <div class="bg-white p-8">
            <h3 class="text-sm font-medium text-gray-900 mb-6">SYSTEM INFORMATION</h3>
            <div class="space-y-4">
              <div class="flex justify-between py-3 border-b border-gray-100">
                <span class="text-sm text-gray-500">Version</span>
                <span class="text-sm text-gray-900">v0.4.0</span>
              </div>
              <div class="flex justify-between py-3 border-b border-gray-100">
                <span class="text-sm text-gray-500">Uptime</span>
                <span class="text-sm text-gray-900">{{ systemInfo.uptime }}</span>
              </div>
              <div class="flex justify-between py-3 border-b border-gray-100">
                <span class="text-sm text-gray-500">Memory Usage</span>
                <span class="text-sm text-gray-900">{{ systemInfo.memory }}</span>
              </div>
              <div class="flex justify-between py-3">
                <span class="text-sm text-gray-500">CPU Usage</span>
                <span class="text-sm text-gray-900">{{ systemInfo.cpu }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mt-8 px-4 sm:px-0 pb-8">
        <div class="bg-white p-8">
          <h3 class="text-sm font-medium text-gray-900 mb-6">QUICK ACTIONS</h3>
          <div class="flex flex-wrap gap-3">
            <button class="btn-primary">Add Gateway</button>
            <button class="btn-secondary">Configure</button>
            <button class="btn-secondary">View Analytics</button>
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

