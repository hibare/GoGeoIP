<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import ThemeToggle from '@/components/ThemeToggle.vue'
import { useTheme } from '@/composables/useTheme'

const { init } = useTheme()

onMounted(() => {
  init()
})
</script>

<template>
  <div class="min-h-screen bg-background text-foreground transition-colors flex flex-col">
    <header class="border-b bg-card/50 backdrop-blur-sm sticky top-0 z-50 shrink-0">
      <div class="container mx-auto flex h-16 items-center justify-between px-4">
        <div class="flex items-center gap-2">
          <img src="/gogeoip.png" alt="GoGeoIP" class="h-8 w-auto" />
          <span class="text-xl font-semibold tracking-tight">GoGeoIP</span>
        </div>
        <div class="flex items-center gap-4">
          <nav class="hidden md:flex items-center gap-6">
            <RouterLink
              to="/"
              class="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
            >
              Lookup
            </RouterLink>
            <RouterLink
              to="/history"
              class="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
            >
              History
            </RouterLink>
            <RouterLink
              to="/settings"
              class="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
            >
              Settings
            </RouterLink>
          </nav>
          <ThemeToggle />
        </div>
      </div>
    </header>
    <main class="container mx-auto py-8 px-4 flex-1">
      <RouterView />
    </main>
    <footer class="border-t py-6 shrink-0">
      <div class="container mx-auto px-4 text-center text-sm text-muted-foreground">
        GoGeoIP &copy; {{ new Date().getFullYear() }}
      </div>
    </footer>
  </div>
</template>
