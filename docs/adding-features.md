# Adding Features — AI Playbook

Step-by-step patterns for common tasks. Follow these to stay consistent with existing codebase.

## Adding a New Entry (Content)

1. Add `src/content/vi/entries/{slug}.md` (required for the story to exist in the canonical set)
2. Optionally add `src/content/{locale}/entries/{slug}.md` (same filename) for localized page body; if omitted, `/{locale}/entries/{slug}` still works via default-locale fallback
3. Add frontmatter matching the Zod schema (see `content-model.md`)
4. Required fields: `name_vi`, `category` (must be one of `CATEGORY_SLUGS`)
5. Set `status: published` to make it visible
6. Set `popularity` to control sort order (higher = shown first)
7. Write markdown body with `## Heading _italic_` pattern
8. Optional: add a GFM comparison table (`| col | col |` with `|---|---|`) for cross-tradition summaries — first column = row label; renders with scroll on mobile (see `content-model.md`, `styling.md`)
9. No routing code changes needed — collections pick up new `.md` files on build

```yaml
---
name_vi: Sơn Tinh
name_han: 山精
name_en: Mountain Spirit
category: than-linh
gender: nam
era: Hùng Vương thứ 18
region: bac
locations:
  - Tản Viên
relations:
  family: []
  allies:
    - Hùng Vương thứ 18
  enemies:
    - Thủy Tinh
sources:
  - title: Lĩnh Nam Chích Quái
    author: Trần Thế Pháp
summary: "..."
themes:
  - thien-nhien
  - chien-tranh
popularity: 7
status: published
---

## Câu _chuyện_

Content here...
```

## Adding a New Category

1. Add slug → **per-locale** labels in `src/data/category-labels.ts`:
   ```typescript
   export const CATEGORY_LABELS: Record<string, Record<Locale, string>> = {
     // ... existing
     "new-slug": { vi: "Nhãn mới", en: "New label" },
   };
   ```
2. Add category card metadata in `src/components/HomePage.astro` (`categoryMeta` object) so the home grid has copy for both locales
3. Category list pages auto-generate via `getStaticPaths()` from `CATEGORY_SLUGS` in `entries/category/[category].astro` and `[lang]/entries/category/[category].astro`

## Adding a New Page

1. Default locale page goes under `src/pages/`.
2. Non-default locale pages must go under `src/pages/[...lang]/...` (dynamic), not `src/pages/en/...`.
3. Import and use `BaseLayout` with `title` and `lang` (fixed `'vi'` for root wrappers, dynamic `lang` for `[lang]` routes):
   ```astro
   ---
   import BaseLayout from '../layouts/BaseLayout.astro';
   const lang = 'vi' as const;
   ---
   <BaseLayout title="..." lang={lang}>
   ```
4. For dynamic routes, export `getStaticPaths()` from `[...lang]` files using `localeStaticPaths()` from `src/i18n/paths.ts`
5. Use `t(lang, 'key')` for user-visible strings; add keys to every supported locale in `src/i18n/config.ts`
6. Update `Header.astro` / footer links only if you add a top-level section

### Example: static About page (`/about` and `/{lang}/about`)

Reference: [`src/pages/about.astro`](../src/pages/about.astro) and [`src/pages/[lang]/about.astro`](../src/pages/[lang]/about.astro).

- `src/pages/about.astro` is default-locale wrapper (`lang="vi"`)
- `src/pages/[lang]/about.astro` renders non-default locales via `getStaticPaths()`
- All copy via `about.*` keys in `src/i18n/config.ts`
- Nav label “Về dự án” / “About” in [`Header.astro`](src/components/Header.astro) uses `localePath(lang, '/about')`

After adding a similar page, update `doc/routing-and-pages.md` and this file if the pattern changes.

## Adding a relation kind (for the graph)

1. Extend `relations` in `src/content.config.ts` (Zod) with a new optional `z.array(z.string())` field.
2. Add the kind to `RelationKind`, `RELATION_KINDS`, and `RELATION_KEYS` in `src/lib/relations-graph.ts`.
3. Add stroke styling in `src/scripts/mount-graph.ts` (`STROKE` map).
4. Add localized label keys under `entry.*` in `src/i18n/config.ts` and map them in `RelationsPage.astro` (`kindUiKey`).

## Tuning the relation matcher

Edit `normalizeForMatch()` and lookup construction in `src/lib/relations-graph.ts`. Prefer exact matching; use `aliases` in frontmatter to disambiguate duplicate names. See [relations-graph.md](./relations-graph.md).

## Adding a New Component

1. Create `src/components/{Name}.astro`
2. Define props interface; pass `lang: Locale` when rendering UI strings:
   ```astro
   ---
   import { t } from '../i18n/config';
   import type { Locale } from '../i18n/config';
   interface Props { lang?: Locale; }
   const { lang = 'vi' } = Astro.props;
   ---
   <div>{t(lang, 'some.key')}</div>
   ```
3. Import in page from the appropriate path depth

### Style Conventions for New Components

- Use CSS variables from `global.css` (see `styling.md`)
- Font: `'Cormorant Garamond'` for headings, `'Be Vietnam Pro'` for body
- Colors: `var(--vermilion)` for accents, `var(--ink)` for text, `var(--paper)` for bg
- Labels: uppercase, small size (`0.7rem`), letter-spacing (`0.15em+`)
- Prefer scoped `<style>` over `<style is:global>`
- Follow existing responsive breakpoints: 1024px, 900px, 600px

## Adding Client-Side Interactivity

Astro components are server-only by default. For interactivity:

### Option A: Vanilla JS (no framework needed)
```astro
<button id="my-btn">Click</button>
<script>
  document.getElementById('my-btn')?.addEventListener('click', () => {
    // handle click
  });
</script>
```

### Option B: Add React (requires setup)
1. `npm install @astrojs/react react react-dom`
2. Add to `astro.config.mjs`:
   ```js
   import react from '@astrojs/react';
   export default defineConfig({ integrations: [react()] });
   ```
3. Create `.tsx` components and use `client:load` or `client:visible` directive:
   ```astro
   ---
   import SearchBar from '../components/SearchBar.tsx';
   ---
   <SearchBar client:load />
   ```

## Adding Search Functionality

Pattern for a search feature:

1. Build search index at build time using `getLocalizedEntries(lang)` (or both collections if you need raw splits)
2. Pass as JSON to a `<script>` tag or island component
3. Filter client-side with vanilla JS or a small library

## Modifying the Entry Detail Page

The entry detail page is in `src/layouts/EntryLayout.astro`. Key areas:

- **Header section**: breadcrumb, title, tags
- **Article content**: hero image, summary, markdown, sources, related
- **Sidebar**: info table, relations, themes
- **Styles**: scoped global CSS at bottom of file

To add a new sidebar section, follow existing `.side-card` blocks and use `t(lang, ...)` for labels.

## Modifying the Home Page

`src/components/HomePage.astro` (used by `src/pages/index.astro` and `src/pages/[lang]/index.astro`) has four sections:
1. `.hero` — hero section with CTAs
2. `#featured` — featured entry cards
3. `#categories` — category grid (links via `localePath(lang, '/entries/category/...')`)
4. `.quote` — blockquote (`#quote`)

Each section's styles are in `HomePage.astro`'s `<style>` block.

## Adding Images

Currently all images are placeholders. To add real images:

1. Place images in `public/images/entries/{entry-id}.jpg`
2. Replace `.img-placeholder` divs with `<img>` tags
3. Or add an `image` field to the content schema:
   ```typescript
   image: z.string().optional(),  // path relative to public/
   ```
4. Reference in templates: `<img src={entry.data.image} alt={entry.data.name_vi} />`

## i18n Reference

Implemented behavior (locales, fallback, `t()`, collections) is documented in [i18n.md](./i18n.md).

## Agent Contract: Dynamic-Only

All agents contributing to this repo must follow:

1. No locale-mirrored page trees (`src/pages/en`, `src/pages/ja`, etc.).
2. Locale routing/features must be implemented once in dynamic paths/helpers.
3. New locale onboarding must be config/content-driven, not page-copy-driven.
4. Any non-dynamic i18n implementation is considered a regression.
