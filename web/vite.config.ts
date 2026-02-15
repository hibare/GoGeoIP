import path from "node:path";
import fs from "node:fs";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "tailwindcss";
import autoprefixer from "autoprefixer";

export default defineConfig({
  plugins: [vue()],
  css: {
    postcss: {
      plugins: [tailwindcss, autoprefixer],
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: (() => {
    const keyPath = path.resolve(__dirname, "../certs/key.pem");
    const certPath = path.resolve(__dirname, "../certs/cert.pem");

    // Only enable HTTPS if certificate files exist (for development)
    if (fs.existsSync(keyPath) && fs.existsSync(certPath)) {
      return {
        https: {
          key: fs.readFileSync(keyPath),
          cert: fs.readFileSync(certPath),
        },
        proxy: {
          "^/api": {
            target: "https://localhost:5000",
            changeOrigin: true,
            secure: false,
          },
        },
      };
    } else {
      return {
        proxy: {
          "^/api": {
            target: "http://localhost:5000",
            changeOrigin: true,
            secure: false,
          },
        },
      };
    }
  })(),
});
