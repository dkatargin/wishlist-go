import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const version = process.env.VITE_VERSION || process.env.APP_VERSION || "rolling";

export default defineConfig({
  plugins: [react()],
  publicDir: "public",
  base: `/static/${version}/`, // базовый путь для assets в S3
  server: {
    port: 3001, // порт dev-сервера
    open: true, // открывать браузер при запуске
  },
  build: {
    outDir: "dist", // папка для прод-сборки
    sourcemap: true, // генерировать sourcemap для дебага (можно true на время разработки)
    minify: true,
  },
  define: {
    __APP_VERSION__: JSON.stringify(version),
    VITE_BACKEND_HOST: JSON.stringify(
      process.env.VITE_BACKEND_HOST || "wish.exo.icu",
    ),
    VITE_BACKEND_PORT: JSON.stringify(process.env.VITE_BACKEND_PORT || "443"),
    VITE_BACKEND_SCHEME: JSON.stringify(
      process.env.VITE_BACKEND_SCHEME || "https",
    ),
  },
});
