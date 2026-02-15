<template>
  <div class="min-h-screen flex items-center justify-center bg-background">
    <div class="text-center space-y-6 max-w-md mx-auto px-4">
      <!-- 404 Icon -->
      <div class="flex justify-center text-9xl text-foreground">
        <GhostIcon class="h-8 w-8" />
      </div>
      <!-- Error Message -->
      <div class="">
        <h1 class="text-2xl font-semibold text-foreground">Page Not Found</h1>
        <p class="text-muted-foreground pt-4">
          The page you're looking for doesn't exist or has been moved.
        </p>
      </div>

      <!-- Actions -->
      <div class="flex flex-col sm:flex-row gap-3 justify-center">
        <Button
          @click="goBack"
          variant="outline"
          class="flex items-center gap-2"
        >
          <ArrowLeftIcon class="h-4 w-4" />
          Go Back
        </Button>
        <Button @click="goHome" class="flex items-center gap-2">
          <HomeIcon class="h-4 w-4" />
          Go to Home
        </Button>
      </div>

      <!-- Additional Help -->
      <div class="pt-4 border-t border-border">
        <p class="text-sm text-muted-foreground">
          If you believe this is an error, please contact support or admin.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { Button } from "@/components/ui/button";
import { HomeIcon, ArrowLeftIcon, GhostIcon } from "lucide-vue-next";

const router = useRouter();

function goBack() {
  if (window.history.length > 1) {
    router.go(-1);
  } else {
    goHome();
  }
}

function goHome() {
  router.push({ name: "home" });
}
</script>
