<script setup lang="ts">
import { ref } from 'vue'
import { Settings, Key, Globe, RefreshCw } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'

const apiBaseUrl = ref(import.meta.env.VITE_API_BASE_URL || 'http://localhost:5000')
const apiKey = ref(import.meta.env.VITE_API_KEY || '')
const autoUpdate = ref(true)
const autoUpdateInterval = ref('24h')

function saveSettings() {
  localStorage.setItem('gogeoip_api_base_url', apiBaseUrl.value)
  localStorage.setItem('gogeoip_api_key', apiKey.value)
  localStorage.setItem('gogeoip_auto_update', String(autoUpdate.value))
  localStorage.setItem('gogeoip_auto_update_interval', autoUpdateInterval.value)
  alert('Settings saved!')
}
</script>

<template>
  <div class="space-y-8">
    <div class="flex flex-col items-center justify-center space-y-4">
      <h1 class="text-3xl font-bold">Settings</h1>
      <p class="text-muted-foreground">Configure your GoGeoIP application</p>
    </div>

    <div class="max-w-2xl mx-auto space-y-6">
      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <Key class="w-5 h-5" />
            API Configuration
          </CardTitle>
          <CardDescription>
            Configure the connection to your GoGeoIP server
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <label class="text-sm font-medium">API Base URL</label>
            <Input v-model="apiBaseUrl" placeholder="http://localhost:5000" />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">API Key</label>
            <Input v-model="apiKey" type="password" placeholder="Enter your API key" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <RefreshCw class="w-5 h-5" />
            Database Settings
          </CardTitle>
          <CardDescription>
            Configure MaxMind database auto-update
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="flex items-center gap-2">
            <input
              type="checkbox"
              id="autoUpdate"
              v-model="autoUpdate"
              class="w-4 h-4"
            />
            <label for="autoUpdate" class="text-sm font-medium">Enable Auto-Update</label>
          </div>
          <div v-if="autoUpdate" class="space-y-2">
            <label class="text-sm font-medium">Update Interval</label>
            <select
              v-model="autoUpdateInterval"
              class="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="6h">Every 6 hours</option>
              <option value="12h">Every 12 hours</option>
              <option value="24h">Every 24 hours</option>
              <option value="48h">Every 48 hours</option>
              <option value="168h">Every week</option>
            </select>
          </div>
        </CardContent>
      </Card>

      <div class="flex justify-end">
        <Button @click="saveSettings">Save Settings</Button>
      </div>
    </div>
  </div>
</template>
