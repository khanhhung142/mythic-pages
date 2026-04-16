# Adding Features — AI Playbook

Step-by-step patterns for common tasks. Follow these to stay consistent with existing codebase.

## Adding a New Entry (Content)

1. Add `src/content/vi/entries/{slug}.md` (required for the story to exist in the canonical set)
2. Optionally add `src/content/en/entries/{slug}.md` with the **same** filename for a full English page body; if omitted, `/en/entries/{slug}` still works and shows the VI markdown via fallback
3. Add frontmatter matching the Zod schema (see `content-model.md`)
4. Required fields: `name_vi`, `category` (must be one of `CATEGORY_SLUGS`)
5. Set `status: published` to make it visible
6. Set `popularity` to control sort order (higher = shown first)
7. Write markdown body with `## Heading _italic_` pattern
8. No routing code changes needed — collections pick up new `.md` files on build

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
2. Add category card metadata in `src/pages/[lang]/index.astro` (`categoryMeta` object) so the home grid has copy for both locales
3. Category list pages auto-generate via `getStaticPaths()` from `CATEGORY_SLUGS` × locales

## Adding a New Page

1. Prefer placing pages under `src/pages/[lang]/...` so URLs stay locale-prefixed
2. Import and use `BaseLayout` with `title` and `lang`:
   ```astro
   ---
   import BaseLayout from '../../layouts/BaseLayout.astro';
   import type { Locale } from '../../i18n/config';
   const { lang } = Astro.params as { lang: Locale };
   ---
   <BaseLayout title="..." lang={lang}>
   ```
3. For dynamic routes, export `getStaticPaths()` including `lang` when needed
4. Use `t(lang, 'key')` for user-visible strings; add keys to `src/i18n/config.ts` for both `vi` and `en`
5. Update `Header.astro` / footer links only if you add a top-level section

### Example: static About page (`/[lang]/about`)

Reference implementation: [`src/pages/[lang]/about.astro`](src/pages/[lang]/about.astro).

- `export async function getStaticPaths() { return locales.map((lang) => ({ params: { lang } })); }`
- `const { lang } = Astro.params as { lang: Locale };`
- `BaseLayout title={pageTitle} lang={lang}` with `pageTitle` built from `t(lang, 'about.title')` and `site.title`
- All copy via `about.*` keys in `src/i18n/config.ts` (vi + en)
- Nav label “Về dự án” / “About” in [`Header.astro`](src/components/Header.astro) points to `/${lang}/about`

After adding a similar page, update `doc/routing-and-pages.md` and this file if the pattern changes.

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

`src/pages/[lang]/index.astro` has four sections:
1. `.hero` — hero section with CTAs
2. `#featured` — featured entry cards
3. `#categories` — category grid (links to `/${lang}/entries/category/...`)
4. `.quote` — blockquote (`#quote`)

Each section's styles are in the same file's `<style>` block.

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

Implemented behavior (locales, fallback, `t()`, collections) is documented in [i18n.md](./i18n.md). When adding UI copy, always add **both** `vi` and `en` keys in `src/i18n/config.ts` unless the string is intentionally locale-only.
