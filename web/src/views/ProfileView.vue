<template>
  <div class="space-y-8">
    <div class="flex flex-col space-y-1">
      <h1 class="text-3xl font-bold">Profile</h1>
      <p class="text-muted-foreground">Manage your account settings</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Account Information</CardTitle>
        <CardDescription
          >View and manage your account details</CardDescription
        >
      </CardHeader>
      <CardContent class="space-y-6">
        <div class="flex items-center gap-6">
          <div
            class="flex h-20 w-20 items-center justify-center rounded-full bg-primary text-primary-foreground text-3xl"
          >
            {{ userStore.initials }}
          </div>
          <div class="space-y-1">
            <p class="text-xl font-medium">{{ userStore.fullName }}</p>
            <p class="text-sm text-muted-foreground">
              {{ userStore.email }}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useUserStore } from "@/store/auth";

const userStore = useUserStore();
</script>
