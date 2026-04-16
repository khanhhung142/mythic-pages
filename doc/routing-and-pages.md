# Routing & Pages

Astro file-based routing now uses a **dynamic locale strategy**:

- default locale (`vi`) routes are unprefixed in `src/pages/*`
- non-default locales are generated via `src/pages/[lang]/*`
- no mirrored locale folders (e.g. `src/pages/en/*`) should be created

Legacy `/vi` redirects are still supported.

## Route Map

| URL Pattern | File | Layout | Description |
|------------|------|--------|-------------|
| `/` | `src/pages/index.astro` | `BaseLayout` (via `HomePage`) | Home (default locale) |
| `/{lang}/` | `src/pages/[lang]/index.astro` | same | Home for any non-default locale |
| `/about` | `src/pages/about.astro` | `BaseLayout` (via `AboutPage`) | About (default locale) |
| `/{lang}/about` | `src/pages/[lang]/about.astro` | same | About for non-default locales |
| `/entries` | `src/pages/entries/index.astro` | `BaseLayout` (via `EntriesListPage`) | Entries catalog (default locale) |
| `/{lang}/entries` | `src/pages/[lang]/entries/index.astro` | same | Entries catalog for non-default locales |
| `/entries/category/[category]` | `src/pages/entries/category/[category].astro` | `BaseLayout` (via `EntriesListPage`) | Category page (default locale) |
| `/{lang}/entries/category/[category]` | `src/pages/[lang]/entries/category/[category].astro` | same | Category page for non-default locales |
| `/entries/[id]` | `src/pages/entries/[id].astro` | `EntryLayout` | Entry detail (default locale) |
| `/{lang}/entries/[id]` | `src/pages/[lang]/entries/[id].astro` | `EntryLayout` | Entry detail for non-default locales |

## Request Flow

```mermaid
graph TD
    subgraph buildTime [BuildTime_getStaticPaths]
        localeCfg[locales_and_defaultLocale] --> langRoutes[pages_[lang]_generate_non_default_paths]
        contentHelpers[getLocalizedEntries_or_getLocalizedEntry] --> publishedOnly[status_published_filter]
        publishedOnly --> staticPaths[paths_lang_id_category]
    end

    subgraph routing [Routing]
        rootRoute["/about"] --> rootFile["src/pages/about.astro"]
        langRoute["/{lang}/about"] --> langFile["src/pages/[lang]/about.astro"]
        rootEntries["/entries/[id]"] --> rootEntry["src/pages/entries/[id].astro"]
        langEntries["/{lang}/entries/[id]"] --> langEntry["src/pages/[lang]/entries/[id].astro"]
    end
```

## Dynamic-First Guardrail

All future route work must be locale-generic:

1. Extend dynamic pages/helpers, never duplicate locale directories.
2. Use `locales` and `defaultLocale` from `src/i18n/config.ts`.
3. Keep `getStaticPaths()` locale loops generic (`locales.filter(...)`).
4. If a PR adds locale-specific page duplication, it violates project routing rules.
