<script setup lang="ts">
import { TooltipProvider } from "@/components/ui/tooltip";
import { Toaster } from "@/components/ui/sonner";
import Nav from "./components/Nav.vue";
</script>

<template>
  <TooltipProvider>
    <div
      class="min-h-screen bg-background text-foreground transition-colors flex flex-col"
    >
      <Nav />
      <main class="container mx-auto py-8 px-4 flex-1">
        <RouterView />
      </main>
      <footer class="border-t py-6 shrink-0">
        <div
          class="container mx-auto px-4 text-center text-sm text-muted-foreground"
        >
          GoGeoIP &copy; {{ new Date().getFullYear() }}
        </div>
      </footer>
      <Toaster position="top-right" richColors />
    </div>
  </TooltipProvider>
</template>
