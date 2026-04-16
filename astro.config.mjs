import { defineConfig } from "astro/config";
import sitemap from "@astrojs/sitemap";

export default defineConfig({
  site: "https://vietmyth.vn",
  i18n: {
    defaultLocale: "vi",
    locales: ["vi", "en"],
    routing: {
      prefixDefaultLocale: false,
    },
  },
  integrations: [
    sitemap({
      i18n: {
        defaultLocale: "vi",
        locales: { vi: "vi-VN", en: "en" },
      },
    }),
  ],
  trailingSlash: "ignore",
  build: {
    format: "directory",
  },
  redirects: {
    "/vi": "/",
  },
});
