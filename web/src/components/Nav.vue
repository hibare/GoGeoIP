<template>
  <header class="border-b dark:border-neon-cyan/20 bg-card/50 dark:bg-black/70 backdrop-blur-sm sticky top-0 z-50 shrink-0">
    <div class="container mx-auto flex h-16 items-center justify-between px-4">
      <RouterLink to="/" class="flex items-center gap-2">
        <img src="/logo.png" alt="Waypoint" class="h-8 w-auto neon-logo transition-all duration-300" />
        <span class="text-xl font-semibold tracking-tight neon-heading transition-all duration-300">Waypoint</span>
      </RouterLink>
      <div class="flex items-center gap-4">
        <nav class="hidden md:flex items-center gap-1">
          <RouterLink
            to="/"
            class="px-3 py-2 text-sm font-medium font-mono rounded-md transition-colors"
            :class="route.path === '/' ? 'neon-text' : 'text-muted-foreground hover:text-foreground'"
          >
            Home
          </RouterLink>
          <RouterLink
            to="/history"
            class="px-3 py-2 text-sm font-medium font-mono rounded-md transition-colors"
            :class="route.path === '/history' ? 'neon-text' : 'text-muted-foreground hover:text-foreground'"
          >
            History
          </RouterLink>
        </nav>
        <div class="relative" v-if="userStore.isAuthenticated">
          <Button variant="ghost" size="sm" class="flex items-center gap-2" @click="showDropdown = !showDropdown">
            <div class="neon-avatar h-8 w-8 bg-primary text-primary-foreground">
              {{ userStore.initials }}
            </div>
            <span class="hidden sm:inline-block font-mono dark:text-neon-cyan/80">{{ userStore.fullName }}</span>
          </Button>
          <div
            v-if="showDropdown"
            class="absolute right-0 mt-2 w-48 rounded-md border dark:border-neon-cyan/20 bg-background dark:bg-black/90 shadow-lg backdrop-blur-sm"
          >
            <div class="py-1">
              <RouterLink to="/profile" class="block px-4 py-2 text-sm neon-menu-item" @click="showDropdown = false">Profile</RouterLink>
              <RouterLink to="/api-keys" class="block px-4 py-2 text-sm neon-menu-item" @click="showDropdown = false">API Keys</RouterLink>
              <button class="w-full text-left px-4 py-2 text-sm neon-menu-item-pink" @click="handleLogout">Sign out</button>
            </div>
          </div>
        </div>
        <Button v-if="!userStore.isAuthenticated" variant="default" @click="router.push('/login')">Login</Button>
        <ThemeToggle />
        <a
          href="https://github.com/hibare/Waypoint"
          target="_blank"
          rel="noopener noreferrer"
          class="text-muted-foreground hover:text-foreground dark:text-neon-cyan/50 dark:hover:text-neon-cyan transition-colors"
        >
          <GithubIcon class="h-5 w-5" />
        </a>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { GithubIcon } from "lucide-vue-next";
import { Button } from "@/components/ui/button";
import ThemeToggle from "@/components/ThemeToggle.vue";
import { useTheme } from "@/composables/useTheme";
import { useUserStore } from "@/store/auth";

const { init } = useTheme();
const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const showDropdown = ref(false);

function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement;
  if (!target.closest(".relative")) {
    showDropdown.value = false;
  }
}

function handleLogout() {
  showDropdown.value = false;
  userStore.logout();
  router.push("/login");
}

onMounted(() => {
  init();
  document.addEventListener("click", handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener("click", handleClickOutside);
});
</script>
