<template>
  <div class="min-h-screen flex items-center justify-center bg-background">
    <div class="text-center space-y-6 max-w-md mx-auto px-4">
      <div class="flex justify-center">
        <GhostIcon class="h-16 w-16 neon-icon-cyan" />
      </div>
      <p class="neon-error-code">404</p>
      <div>
        <h1 class="text-2xl neon-heading">Page Not Found</h1>
        <p class="neon-subheading pt-4">
          // the page you're looking for doesn't exist or has been moved
        </p>
      </div>
      <div class="flex flex-col sm:flex-row gap-3 justify-center">
        <Button @click="goBack" variant="outline" class="flex items-center gap-2">
          <ArrowLeftIcon class="h-4 w-4" />
          Go Back
        </Button>
        <Button @click="goHome" class="flex items-center gap-2">
          <HomeIcon class="h-4 w-4" />
          Go to Home
        </Button>
      </div>
      <div class="pt-4 border-t border-border dark:border-neon-cyan/20">
        <p class="neon-subheading">// if you believe this is an error, contact support or admin</p>
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
