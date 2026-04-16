import { defineConfig } from "astro/config";

export default defineConfig({
  i18n: {
    defaultLocale: 'vi',
    locales: ['vi', 'en'],
    routing: {
      prefixDefaultLocale: true,
      redirectToDefaultLocale: true,
    },
  },
  trailingSlash: "ignore",
  build: {
    format: "directory",
  },
});
