# Architecture

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Framework | Astro | ^6.1.5 |
| Language | TypeScript | ^5.0.0 |
| Styling | Tailwind CSS | ^3.4.17 |
| CSS Processing | PostCSS + Autoprefixer | ^8.5.1 / ^10.4.21 |
| Type Checking | @astrojs/check | ^0.9.0 |
| Content | Astro Content Collections + Zod | built-in |
| Node | 22 (see `.node-version`) |
| Package Manager | npm (also has `bun.lock`) |

## Project Structure

```
mythic-pages/
├── astro.config.mjs          # Astro config: i18n, trailingSlash, build format
├── package.json               # Dependencies and scripts
├── tailwind.config.ts         # Tailwind theme (shadcn-like tokens, mostly unused)
├── tsconfig.json              # Extends astro/tsconfigs/base
├── postcss.config.js          # PostCSS plugins
├── components.json            # shadcn scaffold config (legacy, not active)
├── public/                    # Static assets
│   ├── favicon.ico
│   ├── placeholder.svg
│   └── robots.txt
└── src/
    ├── content.config.ts      # Zod schema + entriesVi / entriesEn collections
    ├── env.d.ts               # Astro type references
    ├── i18n/
    │   ├── config.ts          # locales, defaultLocale, ui strings, t()
    │   └── content.ts         # getLocalizedEntries, getLocalizedEntry, getAllEntryIds
    ├── data/
    │   └── category-labels.ts # Category slug → label per locale (vi/en)
    ├── content/
    │   ├── vi/entries/        # Vietnamese markdown (canonical set of entries)
    │   │   └── *.md
    │   └── en/entries/        # English markdown (optional per entry; missing → VI fallback)
    │       └── *.md
    ├── layouts/
    │   ├── BaseLayout.astro   # Minimal shell: <html lang>, global.css, Header, Footer
    │   └── EntryLayout.astro  # Full entry page: standalone <html>, sidebar, typography
    ├── pages/                 # File-based routing
    │   ├── index.astro        # Redirect / → /vi/
    │   └── [lang]/            # locale segment: vi | en
    │       ├── index.astro    # Home page
    │       └── entries/
    │           ├── index.astro
    │           ├── [id].astro
    │           └── category/
    │               └── [category].astro
    ├── styles/
    │   └── global.css         # CSS variables, reset, base typography
    ├── components/
    │   ├── Header.astro       # Fixed nav bar + lang switch
    │   ├── Footer.astro       # Site footer
    │   ├── EntriesListPage.astro  # Shared list page (catalog + category filter)
    │   ├── EntryCard.astro    # ⚠️ UNUSED — not imported anywhere
    │   ├── InfoTable.astro    # ⚠️ UNUSED — sidebar info table
    │   ├── RelationshipSection.astro  # ⚠️ UNUSED
    │   ├── SidebarCard.astro  # ⚠️ UNUSED
    │   ├── ThemeCloud.astro   # ⚠️ UNUSED
    │   └── wiki/              # ⚠️ UNUSED React components (no React installed)
    │       ├── InfoTable.tsx
    │       ├── RelatedEntries.tsx
    │       ├── RelationshipSection.tsx
    │       ├── SidebarCard.tsx
    │       └── ThemeCloud.tsx
    └── test/
        ├── setup.ts
        └── example.test.ts
```

## Build Pipeline

```mermaid
graph LR
    A["Markdown files<br/>src/content/vi|en/entries/*.md"] --> B[Zod validation<br/>src/content.config.ts]
    B --> C[Astro Content Collections<br/>entriesVi / entriesEn]
    C --> D["getLocalized* helpers<br/>src/i18n/content.ts"]
    D --> E[Astro Pages<br/>src/pages/[lang]/**/*.astro]
    E --> F[Static HTML<br/>dist/]
    G[global.css + Tailwind] --> E
    H[Components + t(lang,key)] --> E
```

## Build Commands

| Command | What It Does |
|---------|-------------|
| `npm run dev` | `astro dev` — local dev server |
| `npm run build` | `astro build` — generate static site to `dist/` |
| `npm run preview` | `astro preview` — preview built site |
| `npm run check` | `astro check` — TypeScript type checking |

## Astro Config

File: `astro.config.mjs`

```js
export default defineConfig({
  i18n: {
    defaultLocale: 'vi',
    locales: ['vi', 'en'],
    routing: {
      prefixDefaultLocale: true,   // /vi/... not bare /
      redirectToDefaultLocale: true, // / → default locale
    },
  },
  trailingSlash: "ignore",  // both /entries and /entries/ work
  build: {
    format: "directory",    // dist/vi/entries/index.html style
  },
});
```

- No integrations installed (no `@astrojs/react`, no `@astrojs/tailwind`)
- No server adapter → pure static output
- Root `/` also redirects explicitly in `src/pages/index.astro` (301 → `/vi/`)

## Deploy

- Output: `dist/` directory (gitignored)
- No `vercel.json`, `netlify.toml`, or Dockerfile
- Compatible with any static hosting
- SEO: `robots.txt` allows common crawlers, entry pages set `<meta description>`
