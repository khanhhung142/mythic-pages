# Routing & Pages

Astro file-based routing now uses a **dynamic locale strategy**:

- all locale routes are generated from a single tree: `src/pages/[...lang]/*`
- default locale (`vi`) uses `lang: undefined` in `getStaticPaths()` for root URLs
- non-default locales use `lang: "<locale>"` for prefixed URLs
- no mirrored locale folders (e.g. `src/pages/en/*`) should be created

Legacy `/vi` redirects are still supported.

## Route Map

| URL Pattern | File | Layout | Description |
|------------|------|--------|-------------|
| `/` | `src/pages/[...lang]/index.astro` | `BaseLayout` (via `HomePage`) | Home (default locale, `lang=undefined`) |
| `/{lang}/` | `src/pages/[...lang]/index.astro` | same | Home (non-default locale, `lang=<locale>`) |
| `/about` | `src/pages/[...lang]/about.astro` | `BaseLayout` (via `AboutPage`) | About (default locale, `lang=undefined`) |
| `/{lang}/about` | `src/pages/[...lang]/about.astro` | same | About (non-default locale, `lang=<locale>`) |
| `/entries` | `src/pages/[...lang]/entries/index.astro` | `BaseLayout` (via `EntriesListPage`) | Entries catalog (default locale) |
| `/{lang}/entries` | `src/pages/[...lang]/entries/index.astro` | same | Entries catalog (non-default locales) |
| `/entries/category/[category]` | `src/pages/[...lang]/entries/category/[category].astro` | `BaseLayout` (via `EntriesListPage`) | Category page (default locale) |
| `/{lang}/entries/category/[category]` | `src/pages/[...lang]/entries/category/[category].astro` | same | Category page (non-default locales) |
| `/entries/[id]` | `src/pages/[...lang]/entries/[id].astro` | `EntryLayout` | Entry detail (default locale) |
| `/{lang}/entries/[id]` | `src/pages/[...lang]/entries/[id].astro` | `EntryLayout` | Entry detail (non-default locales) |

## Request Flow

```mermaid
graph TD
    subgraph buildTime [BuildTime_getStaticPaths]
        localeCfg[locales_and_defaultLocale] --> langRoutes[pages_[...lang]_generate_all_locale_paths]
        contentHelpers[getLocalizedEntries_or_getLocalizedEntry] --> publishedOnly[status_published_filter]
        publishedOnly --> staticPaths[paths_lang_id_category]
    end

    subgraph routing [Routing]
        rootRoute["/about"] --> rootFile["src/pages/[...lang]/about.astro_lang=undefined"]
        langRoute["/{lang}/about"] --> langFile["src/pages/[...lang]/about.astro_lang=locale"]
        rootEntries["/entries/[id]"] --> rootEntry["src/pages/[...lang]/entries/[id].astro_lang=undefined"]
        langEntries["/{lang}/entries/[id]"] --> langEntry["src/pages/[...lang]/entries/[id].astro_lang=locale"]
    end
```

## Dynamic-First Guardrail

All future route work must be locale-generic:

1. Extend dynamic pages/helpers, never duplicate locale directories.
2. Use `locales` and `defaultLocale` from `src/i18n/config.ts`.
3. Keep `getStaticPaths()` locale loops generic via shared helpers (for example `localeStaticPaths()` in `src/i18n/paths.ts`).
4. If a PR adds locale-specific page duplication, it violates project routing rules.
