<template>
  <div class="min-h-screen flex items-center justify-center bg-background px-4">
    <Card class="w-full max-w-md">
      <CardHeader class="space-y-1 text-center">
        <div class="flex justify-center mb-4">
          <img src="/logo.png" alt="Waypoint" class="h-20 w-auto neon-logo" />
        </div>
        <CardDescription class="neon-subheading">
          // sign in to access your IP lookup history and settings
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div
          v-if="error"
          class="text-sm text-center p-4 rounded-md text-red-500 bg-red-50 dark:bg-neon-pink/10 dark:border dark:border-neon-pink/30 dark:text-neon-pink font-mono"
        >
          {{ error }}
        </div>
        <div v-else-if="loading" class="text-center text-sm text-muted-foreground py-4">
          <div class="flex justify-center items-center gap-2 font-mono dark:text-neon-cyan/70">
            <div class="animate-spin h-4 w-4 border-2 border-primary border-t-transparent rounded-full dark:border-neon-cyan dark:border-t-transparent"></div>
            <span>Redirecting to login...</span>
          </div>
        </div>
        <Button v-if="!loading" class="w-full" @click="login">Sign in</Button>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import axios from "@/lib/axios";

const router = useRouter();
const route = useRoute();
const loading = ref(false);
const error = ref("");

const login = async () => {
  loading.value = true;
  error.value = "";

  const redirect = (route.query.redirect as string) || "/";

  try {
    const response = await axios.get("/api/v1/auth/login", {
      params: { redirect },
    });

    const redirectUrl = response.data.redirect_url;
    window.location.href = redirectUrl;
  } catch (e) {
    error.value = "Failed to initiate login. Please try again.";
    loading.value = false;
  }
};
</script>
