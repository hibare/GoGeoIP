<template>
  <div class="min-h-screen flex items-center justify-center bg-background px-4">
    <Card class="w-full max-w-md">
      <CardHeader class="space-y-1 text-center">
        <div class="flex justify-center mb-4">
          <img src="/logo.png" alt="GoGeoIP" class="h-20 w-auto" />
        </div>
        <CardDescription>
          Sign in to access your IP lookup history and settings
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div
          v-if="error"
          class="text-sm text-red-500 text-center p-4 bg-red-50 rounded-md"
        >
          {{ error }}
        </div>
        <div
          v-else-if="loading"
          class="text-center text-sm text-muted-foreground py-4"
        >
          <div class="flex justify-center items-center gap-2">
            <div
              class="animate-spin h-4 w-4 border-2 border-primary border-t-transparent rounded-full"
            ></div>
            <span>Redirecting to login...</span>
          </div>
        </div>
        <Button v-if="!loading" class="w-full" @click="login"> Sign in </Button>
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
