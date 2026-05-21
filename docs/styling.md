# Styling System

The site uses three styling layers: CSS custom properties (design tokens), Tailwind CSS, and scoped component styles.

## Design Tokens (CSS Variables)

File: `src/styles/global.css`

### Color Palette

| Variable | Value | Usage |
|----------|-------|-------|
| `--paper` | `#f4ead5` | Background — aged paper tone |
| `--paper-dark` | `#e8d9b8` | Card backgrounds, darker paper |
| `--paper-light` | `#faf4e3` | Summary boxes, sidebar cards |
| `--ink` | `#1a0f08` | Primary text color — deep brown/black |
| `--ink-soft` | `#3d2817` | Body text, secondary content |
| `--ink-mute` | `#6b5945` | Muted text, breadcrumbs |
| `--ink-light` | `#9c8a72` | Labels, metadata |
| `--vermilion` | `#a8321e` | Accent color — used for CTAs, tags, highlights |
| `--vermilion-dark` | `#7a1e10` | Darker accent variant |
| `--gold` | `#b8860b` | Secondary accent — category section, Hán tự text |
| `--jade` | `#2d5a3d` | Tertiary accent (defined but rarely used) |
| `--line` | `rgba(26,15,8,.15)` | Border lines |
| `--line-soft` | `rgba(26,15,8,.08)` | Subtle borders |
| `--shadow` | `rgba(26,15,8,.12)` | Box shadows |

### Design Aesthetic

The visual language evokes **aged Vietnamese manuscript**:
- Paper-tone backgrounds with subtle noise texture (SVG filter in body)
- Radial gradient overlays (vermilion top-left, gold bottom-right)
- Decorative Chinese characters as watermarks (e.g. `神話` on hero, `目錄` on catalog)
- Vermilion as primary accent (traditional Vietnamese color)

## Font System

Three font families loaded via Google Fonts:

| Font | Usage | Weight |
|------|-------|--------|
| **Cormorant Garamond** | Headings, decorative text, large titles | 400–700, italic |
| **Playfair Display** | Italic accents, subtitles, blockquotes | Italic 400, 600 |
| **Be Vietnam Pro** | Body text, UI elements | 300–600 |

Fonts loaded in two places:
- `global.css` — `@import url(...)` (used by BaseLayout pages)
- `EntryLayout.astro` — `<link>` tags in `<head>` (standalone)

### Typography Scale

**Headings** (Cormorant Garamond):
- Hero h1: `clamp(3.5rem, 8vw, 7rem)`
- Section titles: `clamp(2.5rem, 5vw, 4rem)`
- Entry title: `clamp(3.5rem, 8vw, 6.5rem)`
- Entry h2: `2.5rem`
- Entry h3: `1.6rem`

**Body** (Be Vietnam Pro):
- Base: `1.05rem`, line-height `1.85`, weight `300`

**UI elements**:
- Labels/tags: `0.7–0.75rem`, letter-spacing `0.15–0.4em`, uppercase
- Navigation: `0.85rem`, uppercase

## Styling Approach

### 1. Global CSS (`src/styles/global.css`)
- CSS reset (`*, margin, padding, box-sizing`)
- CSS variables (`:root`)
- Body background (paper texture + gradients)
- `.prose` styles for markdown content (defined but barely used — EntryLayout has its own)

### 2. Scoped Styles in `.astro` Components
Most styling lives inside `<style>` or `<style is:global>` blocks within `.astro` files:

| Component | Style Type | Scope |
|-----------|-----------|-------|
| `Header.astro` | `<style is:global>` | Nav styles leak globally |
| `Footer.astro` | `<style is:global>` | Footer styles leak globally |
| `EntryLayout.astro` | `<style is:global>` | ~400 lines of entry page styles |
| `EntriesListPage.astro` | `<style>` (scoped) | Catalog page styles |
| `index.astro` | `<style>` (scoped) | Home page styles |

### 3. Tailwind CSS
Configured but lightly used. The `tailwind.config.ts` has shadcn-style tokens (HSL variables) but the actual site mostly uses vanilla CSS. Tailwind is available for new components.

## Common CSS Patterns

### Image Placeholder
Used across all pages for missing images:

```css
.img-placeholder {
  background: linear-gradient(135deg, var(--paper-dark), #d4c19a);
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(26,15,8,.15);
}
.img-placeholder::before {  /* diagonal lines pattern */ }
.img-placeholder::after {   /* 圖 character centered */ }
```

### Eyebrow Label
Small uppercase text with vermilion line:

```css
.eyebrow {
  font-size: .75rem;
  letter-spacing: .4em;
  text-transform: uppercase;
  color: var(--vermilion);
}
.eyebrow::before { content: ''; width: 40px; height: 1px; background: var(--vermilion); }
```

### Tag Pill
```css
.tag {
  font-size: .7rem;
  letter-spacing: .15em;
  text-transform: uppercase;
  padding: .5rem 1rem;
  border: 1px solid var(--ink);
}
.tag.primary { background: var(--vermilion); color: var(--paper); }
```

### Section Label (with trailing line)
```css
.section-label {
  font-size: .7rem;
  letter-spacing: .3em;
  text-transform: uppercase;
  color: var(--vermilion);
}
.section-label::after { content: ''; flex: 1; height: 1px; background: var(--line); }
```

## Entry content tables

Markdown GFM pipe tables in entry bodies are styled in `EntryLayout.astro` under `.entry-content` (not in `global.css` `.prose`).

| Behavior | Implementation |
|----------|----------------|
| Layout containment | `.entry-content` uses `min-width: 0`, `max-width: 100%`, `overflow-x: clip` so wide tables cannot blow out the grid column |
| Outer frame | `border` on `.table-scroll` (not on `<table>`), `border-radius: 2px`, `min-width: 0`, horizontal scroll inside wrapper |
| Table width in scroll | `width: max-content; min-width: 100%` inside `.table-scroll` — scrolls when wider than article column |
| Cell borders & zebra | `border: 1px solid var(--line)` on `th`/`td`; even rows `var(--paper-light)` |
| Header row | `th` uses `var(--paper-dark)` background, `Cormorant Garamond` |
| Italics in cells | `th em` / `td em` reset to body font so headers stay readable |
| Wide tables (desktop/tablet) | Rehype wraps each `<table>` in `<div class="table-scroll">` (`src/lib/rehype-wrap-tables.ts`); wrapper scrolls horizontally, table keeps `display: table` |
| Row labels on scroll | First column `position: sticky; left: 0` with matching row background |
| Mobile comparison view (≤768px) | Inline script in `EntryLayout.astro` reads each table and inserts `.comp-cards` after the wrapper: tab pills per source column (col 2+), one panel with attribute rows (`dt`/`dd`). Table hidden; cards shown. No markdown changes. |
| Mobile density | `@media (max-width: 600px)` — smaller table cells; `.comp-cards` stacks label above value |

**`.comp-cards` structure** (generated client-side): `.comp-tabs` → `.comp-tab` buttons; `.comp-panels` → `.comp-panel` with `.comp-source-name` + `.comp-rows` / `.comp-row` (`dt` label, `dd` value). Tab label = header text before `(` or truncated to 24 chars.

Authoring: see `content-model.md` (GFM pipe tables). First column = row attribute; following columns = sources/variants to compare.

## Responsive Breakpoints

| Breakpoint | Target |
|-----------|--------|
| `max-width: 1024px` | Tablet — entry page goes single column |
| `max-width: 900px` | Small tablet — hero stacks, nav hides, footer 2-col |
| `max-width: 768px` | Phone — entry tables hidden; `.comp-cards` source-column view shown |
| `max-width: 600px` | Small phone — entries grid 1-col; denser table cells and stacked comp-card rows |

## Tailwind Config Notes

File: `tailwind.config.ts`

- `darkMode: ["class"]` — dark mode ready but not used
- `content: ["./src/**/*.{astro,ts,tsx,js,jsx}"]`
- Extends with shadcn-style HSL color tokens (none actually defined in CSS)
- Font families set to Inter/Lora (NOT matching actual fonts used)
- The shadcn tokens (`--primary`, `--secondary`, etc.) are NOT defined in `global.css`

**This is a config drift issue** — see `known-issues.md`.
