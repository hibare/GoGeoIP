import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { toast } from "vue-sonner";
import router from "@/router";
import {
  login as apiLogin,
  logout as apiLogout,
  getProfile,
} from "@/apis/auth";

import { type User } from "@/types/auth";

export const useUserStore = defineStore("user", () => {
  // State
  const user = ref<User | null>(null);
  const checkLoading = ref(true);
  const authLoading = ref(false);
  const hasCheckedAuth = ref(false);

  // Getters
  const fullName = computed(() => {
    if (!user.value) return "";
    const { first_name, last_name } = user.value;
    return `${first_name} ${last_name}`.trim();
  });

  const initials = computed(() => {
    if (!user.value) return "";
    const { first_name, last_name } = user.value;
    return `${first_name.charAt(0)}${last_name.charAt(0)}`.toUpperCase();
  });

  const isAuthenticated = computed(() => !!user.value);
  const email = computed(() => user.value?.email || "");

  const clearUser = () => {
    user.value = null;
  };

  // Actions
  const login = async (redirect?: string) => {
    authLoading.value = true;
    try {
      const { redirectUrl } = await apiLogin(redirect);
      window.location.href = redirectUrl;
    } catch (error) {
      authLoading.value = false;
      console.error("Login error:", error);
      toast.error("Login failed. Please try again.");
    }
    // Let loading continue though window redirect, only stop when there is an error
  };

  const logout = async () => {
    authLoading.value = true;

    try {
      await apiLogout();
    } catch (error) {
      console.error("Logout error:", error);
      toast.error("Logout failed. Please try again.");
    }

    // Navigate to login BEFORE clearing user state
    // This allows the router guard to see we are still authenticated but loading (logging out)
    await router.push({ name: "login" });

    // Clear user state
    clearUser();
    authLoading.value = false;
  };

  const checkAuth = async () => {
    if (hasCheckedAuth.value) return;

    checkLoading.value = true;
    try {
      const userData = await getProfile();
      user.value = userData;
    } catch (e) {
      console.error("Auth check failed:", e);
      clearUser();
    } finally {
      checkLoading.value = false;
      hasCheckedAuth.value = true;
    }
  };

  return {
    // State
    user,
    checkLoading,
    authLoading,
    hasCheckedAuth,

    // Getters
    fullName,
    initials,
    isAuthenticated,
    email,

    // Actions
    login,
    logout,
    checkAuth,
    clearUser,
  };
});
