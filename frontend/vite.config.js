import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const version = process.env.VITE_VERSION || process.env.APP_VERSION || "rolling";

export default defineConfig({
  plugins: [react()],
  publicDir: "public",
  base: `/static/${version}/`, // базовый путь для assets в S3
  server: {
    port: 3002, // порт dev-сервера
    open: false, // открывать браузер при запуске
  },
  build: {
    outDir: "dist", // папка для прод-сборки
    sourcemap: true, // генерировать sourcemap для дебага (можно true на время разработки)
    minify: true,
  },
  define: {
    __APP_VERSION__: JSON.stringify(version),
    "import.meta.env.VITE_BACKEND_HOST": JSON.stringify(
      process.env.VITE_BACKEND_HOST || "wish.dimhost.ru",
    ),
    "import.meta.env.VITE_BACKEND_PORT": JSON.stringify(process.env.VITE_BACKEND_PORT || "443"),
    "import.meta.env.VITE_BACKEND_SCHEME": JSON.stringify(
      process.env.VITE_BACKEND_SCHEME || "https",
    ),
    "import.meta.env.VITE_DEPLOYMENT_TYPE": JSON.stringify(
      process.env.VITE_DEPLOYMENT_TYPE || "production",
    ),
    "import.meta.env.VITE_AUTH_MOCKUP": JSON.stringify(
      process.env.VITE_AUTH_MOCKUP || "",
    ),
  },
});
