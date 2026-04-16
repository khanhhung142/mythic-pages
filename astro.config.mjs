import { defineConfig } from "astro/config";
import sitemap from "@astrojs/sitemap";

export default defineConfig({
  site: "https://vietmyth.vn",
  i18n: {
    defaultLocale: "vi",
    // Keep this list in sync with src/i18n/config.ts -> locales.
    locales: ["vi", "en"],
    routing: {
      prefixDefaultLocale: false,
    },
  },
  integrations: [
    sitemap({
      i18n: {
        defaultLocale: "vi",
        // Keep this map in sync with i18n.locales above.
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
