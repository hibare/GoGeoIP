<template>
  <div class="min-h-screen flex items-center justify-center bg-background">
    <div class="text-center space-y-6 max-w-md mx-auto px-4">
      <div class="flex justify-center">
        <ServerOff class="h-16 w-16 neon-icon-pink" />
      </div>
      <p class="neon-error-code">500</p>
      <div>
        <h1 class="text-2xl neon-heading">Server Error</h1>
        <div class="mt-4 p-3 border-2 border-dashed border-destructive/50 rounded-md bg-destructive/5 dark:border-neon-pink/40 dark:bg-neon-pink/5">
          <p class="text-destructive text-sm font-medium font-mono dark:text-neon-pink">{{ errorMessage }}</p>
        </div>
      </div>
      <div class="flex flex-col sm:flex-row gap-3 justify-center">
        <Button @click="goHome" class="flex items-center gap-2">
          <Home class="h-4 w-4" />
          Go to Home
        </Button>
      </div>
      <div class="pt-4 border-t border-border dark:border-neon-pink/20">
        <p class="neon-subheading">// if the problem persists, contact support or admin</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Button } from "@/components/ui/button";
import { Home, ServerOff } from "lucide-vue-next";

const route = useRoute();
const router = useRouter();

const errorMessage = computed(() => {
  const message = route.query.message as string;
  return message || "Something went wrong. Please try again later.";
});

function goHome() {
  router.push({ name: "home" });
}
</script>
