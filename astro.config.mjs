import { defineConfig } from "astro/config";

export default defineConfig({
  i18n: {
    defaultLocale: "vi",
    locales: ["vi", "en"],
    routing: {
      prefixDefaultLocale: false,
    },
  },
  trailingSlash: "ignore",
  build: {
    format: "directory",
  },
  redirects: {
    "/vi": "/",
  },
});
