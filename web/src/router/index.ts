import { createRouter, createWebHistory } from "vue-router";
import { useUserStore } from "@/store/auth";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("@/views/HomeView.vue"),
      meta: {
        title: "Home",
        requiresAuth: false,
      },
    },
    {
      path: "/history",
      name: "history",
      component: () => import("@/views/HistoryView.vue"),
      meta: {
        title: "History",
        requiresAuth: false,
      },
    },
    {
      path: "/profile",
      name: "profile",
      component: () => import("@/views/ProfileView.vue"),
      meta: {
        title: "Profile",
        requiresAuth: true,
      },
    },
    {
      path: "/api-keys",
      name: "api-keys",
      component: () => import("@/views/APIKeysView.vue"),
      meta: {
        title: "API Keys",
        requiresAuth: true,
      },
    },
    {
      path: "/login",
      name: "login",
      component: () => import("@/views/LoginView.vue"),
      meta: {
        title: "Login",
        requiresAuth: false,
      },
    },
    {
      path: "/500",
      name: "server-error",
      component: () => import("@/views/Error500.vue"),
      meta: {
        title: "Server Error",
      },
    },
    {
      path: "/:pathMatch(.*)*",
      name: "not-found",
      component: () => import("@/views/Error404.vue"),
      meta: {
        title: "Page Not Found",
      },
    },
  ],
});

// Navigation guard for authentication
router.beforeEach(async (to, _from, next) => {
  const requiresAuth = to.meta.requiresAuth as boolean;
  const userStore = useUserStore();

  // Make sure we've checked auth state at least once
  if (!userStore.hasCheckedAuth) {
    await userStore.checkAuth();
  }

  if (requiresAuth) {
    // For routes that require auth, check if user is authenticated
    if (userStore.isAuthenticated) {
      // User is authenticated, allow navigation
      next();
    } else {
      // User is not authenticated, redirect to login
      next({ path: "/login", query: { redirect: to.path } });
    }
  } else {
    // Public route
    // If user is authenticated and trying to access login page, redirect to dashboard
    // UNLESS we are currently logging out (authLoading is true)
    if (
      to.path === "/login" &&
      userStore.isAuthenticated &&
      !userStore.authLoading
    ) {
      next({ path: "/" });
    } else {
      next();
    }
  }
});

router.afterEach((to) => {
  const title = to.meta.title as string;
  if (title) {
    document.title = `${title}`;
  }
});

export default router;
