<script setup lang="ts">
import { computed } from "vue";
import { useDateFormat } from "@vueuse/core";
import { TooltipProvider } from "@/components/ui/tooltip";
import { Toaster } from "@/components/ui/sonner";
import Nav from "./components/Nav.vue";
import { VERSION, BUILD_TIMESTAMP, COMMIT_HASH } from "@/lib/constants";

const formattedBuildTimestamp = computed(() => {
  if (BUILD_TIMESTAMP === "unknown") return "unknown";
  return useDateFormat(BUILD_TIMESTAMP, "YYYY-MM-DD HH:mm");
});
</script>

<template>
  <TooltipProvider>
    <div class="min-h-screen bg-background text-foreground transition-colors flex flex-col grid-bg">
      <Nav />
      <main class="container mx-auto py-8 px-4 flex-1">
        <RouterView />
      </main>
      <footer class="border-t py-4 shrink-0 dark:border-neon-cyan/20 footer-glow">
        <div class="container mx-auto px-4 text-center text-xs font-mono">
          <div class="space-y-1">
            <div class="neon-text opacity-50">
              {{ VERSION }} | {{ formattedBuildTimestamp }} | {{ COMMIT_HASH.slice(0, 7) }}
            </div>
            <div class="text-muted-foreground tracking-widest uppercase text-[10px]">Waypoint &copy; {{ new Date().getFullYear() }}</div>
          </div>
        </div>
      </footer>
      <Toaster position="top-right" richColors />
    </div>
  </TooltipProvider>
</template>
