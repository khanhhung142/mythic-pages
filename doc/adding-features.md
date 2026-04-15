# Adding Features — AI Playbook

Step-by-step patterns for common tasks. Follow these to stay consistent with existing codebase.

## Adding a New Entry (Content)

1. Create `src/content/entries/{slug}.md`
2. Add frontmatter matching the Zod schema (see `content-model.md`)
3. Required fields: `name_vi`, `category` (must be one of `CATEGORY_SLUGS`)
4. Set `status: published` to make it visible
5. Set `popularity` to control sort order (higher = shown first)
6. Write markdown body with `## Heading _italic_` pattern
7. No code changes needed — Astro auto-discovers new `.md` files

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

1. Add slug → label in `src/data/category-labels.ts`:
   ```typescript
   export const CATEGORY_LABELS: Record<string, string> = {
     // ... existing
     "new-slug": "Nhãn mới",
   };
   ```
2. Add category card in `src/pages/index.astro` (hardcoded `categories` array in frontmatter)
3. Category page auto-generates via `getStaticPaths()` from `CATEGORY_SLUGS`

## Adding a New Page

1. Create `src/pages/{path}.astro`
2. Import and use `BaseLayout`:
   ```astro
   ---
   import BaseLayout from '../layouts/BaseLayout.astro';
   ---
   <BaseLayout title="Page Title">
     <!-- content -->
   </BaseLayout>
   ```
3. For dynamic routes, export `getStaticPaths()`:
   ```astro
   ---
   export async function getStaticPaths() {
     return [
       { params: { slug: 'value' }, props: { /* data */ } },
     ];
   }
   ---
   ```
4. Add navigation link in `Header.astro` if needed

## Adding a New Component

1. Create `src/components/{Name}.astro`
2. Define props interface:
   ```astro
   ---
   interface Props {
     title: string;
     items: string[];
   }
   const { title, items } = Astro.props;
   ---
   <div class="my-component">
     <h3>{title}</h3>
     {items.map(item => <span>{item}</span>)}
   </div>

   <style>
   .my-component { /* scoped styles */ }
   </style>
   ```
3. Import in page: `import MyComponent from '../components/MyComponent.astro';`

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

1. Build search index at build time in a page's frontmatter:
   ```typescript
   const entries = await getCollection('entries', e => e.data.status === 'published');
   const searchData = entries.map(e => ({
     id: e.id,
     name: e.data.name_vi,
     summary: e.data.summary,
     category: e.data.category,
   }));
   ```
2. Pass as JSON to a `<script>` tag or island component
3. Filter client-side with vanilla JS or a small library

## Modifying the Entry Detail Page

The entry detail page is in `src/layouts/EntryLayout.astro`. Key areas:

- **Header section**: lines 65–96 (breadcrumb, title, tags)
- **Article content**: lines 100–148 (hero image, summary, markdown, sources, related)
- **Sidebar**: lines 150–200 (info table, relations, themes)
- **Styles**: lines 210–594 (all scoped global CSS)

To add a new sidebar section:
```astro
<!-- Add inside .sidebar-sticky, after existing .side-card blocks -->
<div class="side-card">
  <div class="side-label">New Section <span class="num">iv.</span></div>
  <!-- content -->
</div>
```

## Modifying the Home Page

`src/pages/index.astro` has four sections:
1. `.hero` — hero section with CTAs
2. `#featured` — featured entry cards
3. `#categories` — category grid (dark bg)
4. `.quote` — blockquote

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

## Adding i18n (Internationalization)

The site has non-functional VI/EN buttons. To implement:

1. Create content in both languages (e.g. `src/content/entries/en/thanh-giong.md`)
2. Or add translated fields to schema (`summary_en`, `name_display_en`)
3. Use Astro's i18n routing: `src/pages/en/entries/[id].astro`
4. Update Header language buttons to link to `/en/` prefix routes
5. Store locale preference in URL path (not cookies — static site)
