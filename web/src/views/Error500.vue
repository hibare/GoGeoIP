<template>
  <div class="min-h-screen flex items-center justify-center bg-background">
    <div class="text-center space-y-6 max-w-md mx-auto px-4">
      <!-- 500 Icon -->
      <div class="flex justify-center text-9xl text-foreground">
        <ServerOff class="h-8 w-8" />
      </div>

      <!-- Error Message -->
      <div class="">
        <h1 class="text-2xl font-semibold text-foreground">Server Error</h1>
        <div
          class="mt-4 p-3 border-2 border-dashed border-destructive/50 rounded-md bg-destructive/5"
        >
          <p class="text-destructive text-sm font-medium">
            {{ errorMessage }}
          </p>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex flex-col sm:flex-row gap-3 justify-center">
        <Button @click="goHome" class="flex items-center gap-2">
          <Home class="h-4 w-4" />
          Go to Home
        </Button>
      </div>

      <!-- Additional Help -->
      <div class="pt-4 border-t border-border">
        <p class="text-sm text-muted-foreground">
          If the problem persists, please contact support or admin.
        </p>
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
