# Known Issues & Config Drift

Issues discovered during codebase audit. Address these before they cause confusion.

## 1. Unused Components (Dead Code)

**Severity**: Low — no runtime impact

| File | Status |
|------|--------|
| `src/components/EntryCard.astro` | Not imported anywhere |
| `src/components/InfoTable.astro` | Not imported anywhere |
| `src/components/RelationshipSection.astro` | Not imported anywhere |
| `src/components/SidebarCard.astro` | Not imported anywhere |
| `src/components/ThemeCloud.astro` | Not imported anywhere |
| `src/components/wiki/*.tsx` (5 files) | React components — no React installed |

**Action**: Delete or refactor. The `.astro` versions duplicate functionality already in `EntryLayout.astro`. The `.tsx` files are completely non-functional.

## 2. Tailwind Config Mismatch

**Severity**: Medium — misleading config

`tailwind.config.ts` references shadcn-style HSL variables (`--primary`, `--secondary`, `--border`, etc.) that are **never defined** in `global.css`. The actual design tokens are `--paper`, `--ink`, `--vermilion`, etc.

Font families in Tailwind config (`Inter`, `Lora`) don't match actual fonts (`Cormorant Garamond`, `Be Vietnam Pro`, `Playfair Display`).

**Action**: Either align Tailwind config with actual design tokens, or strip the shadcn tokens.

## 3. components.json Pointing to Wrong CSS

**Severity**: Low

`components.json` (shadcn config) points to `src/index.css` which doesn't exist. Actual global CSS is at `src/styles/global.css`.

**Action**: Delete `components.json` if not using shadcn, or update the path.

## 4. Duplicate Font Loading

**Severity**: Low — performance

Google Fonts loaded in two places:
- `src/styles/global.css` via `@import url(...)`
- `src/layouts/EntryLayout.astro` via `<link>` tags in `<head>`

`EntryLayout` doesn't use `BaseLayout` (which imports `global.css`), so it loads fonts independently. This means fonts are loaded twice if a user navigates from a BaseLayout page to an entry page.

**Action**: Consolidate font loading — either always use `<link>` in a shared head partial, or always use CSS `@import`.

## 5. EntryLayout Not Extending BaseLayout

**Severity**: Medium — maintenance burden

`EntryLayout.astro` is a complete standalone HTML document (~595 lines including styles). It duplicates `Header`, `Footer`, font loading, and `<html>` structure from `BaseLayout`.

**Action**: Consider refactoring `EntryLayout` to extend `BaseLayout` with a slot-based approach, keeping entry-specific styles scoped.

## 6. Home Page Category Links

**Severity**: Low — UX bug

In `src/pages/index.astro`, category cards all link to `/entries` instead of `/entries/category/{slug}`:

```astro
<!-- Current (broken) -->
<a href="/entries" class="cat-card">

<!-- Should be -->
<a href={`/entries/category/${cat.slug}`} class="cat-card">
```

**Action**: Fix the `href` to use category slug.

## 7. Non-functional Language Switcher

**Severity**: Low — placeholder feature

`Header.astro` has VI/EN buttons with no click handlers and no i18n routing.

**Action**: Either remove the buttons or implement i18n (see `adding-features.md`).

## 8. Placeholder Content

**Severity**: Low — cosmetic

- Hero section has lorem ipsum text
- Quote section has lorem ipsum
- All images are CSS-only placeholders (no actual images)
- Footer links (Đóng góp, Github, Email) point to `#`

**Action**: Replace with real content when ready.

## 9. ESLint/Vitest Config References React

**Severity**: Low — tooling

`eslint.config.js` imports React Hooks plugin. `vitest.config.ts` references React SWC. These dependencies are not in `package.json`.

**Action**: Remove React references from linting/testing config, or install React if planning to use it.

## 10. Playwright Config References Missing Package

**Severity**: Low — testing

`playwright.config.ts` imports from `lovable-agent-playwright-config` and `playwright-fixture.ts` imports from `@playwright/test`. Neither package is in `package.json`.

**Action**: Install Playwright dependencies or remove the config files.

## 11. Relations Not Linked

**Severity**: Medium — missing feature

Relations in frontmatter (family, allies, enemies) are plain text strings, not linked to other entry IDs. For example, `enemies: ["Giặc Ân"]` can't link to another entry page.

**Action**: Consider adding an `entry_id` field to relations, or auto-matching relation text to existing entry `name_vi` values.

## Priority Order

1. **Fix category links on home page** (#6) — quick fix, improves UX
2. **Delete unused components** (#1) — reduces confusion
3. **Align Tailwind config** (#2) — prevents misuse
4. **Refactor EntryLayout** (#5) — reduces duplication
5. **Link relations** (#11) — improves content interconnection
6. Other items as needed
