import { defineConfig } from "astro/config";

// `build.format: 'directory'` → `dist/entries/index.html`, `dist/entries/ho-tinh/index.html`
// Many static hosts only resolve /entries if that folder has index.html (not root-level entries.html).
// `trailingSlash: 'ignore'` → both `/entries` and `/entries/` work in dev/preview.
export default defineConfig({
  trailingSlash: "ignore",
  build: {
    format: "directory",
  },
});
