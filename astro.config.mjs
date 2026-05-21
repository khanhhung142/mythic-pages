import { defineConfig, envField } from "astro/config";
import sitemap from "@astrojs/sitemap";
import { rehypeWrapTables } from "./src/lib/rehype-wrap-tables.ts";

export default defineConfig({
  markdown: {
    rehypePlugins: [rehypeWrapTables],
  },
  site: "https://vietmyth.vn",

  env: {
    schema: {
      PUBLIC_TURNSTILE_SITE_KEY: envField.string({
        context: "client",
        access: "public",
        optional: true,
      }),
    },
  },

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