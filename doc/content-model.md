# Content Model

Myth entries live as Markdown under **locale-specific folders**:

- `src/content/vi/entries/*.md` — Vietnamese (canonical set)
- `src/content/en/entries/*.md` — English translations (optional per entry)

Each file = one mythological entry. Same filename (`id`) across locales denotes the same story. If a localized file is missing, the site serves default-locale content for that route (see `src/i18n/content.ts`).

## Schema Definition

File: `src/content.config.ts`

Collections are generated dynamically from `locales` and share one schema (`entrySchema`):

```typescript
const collections = Object.fromEntries(
  locales.map((locale) => [
    `entries${locale.charAt(0).toUpperCase()}${locale.slice(1)}`,
    defineCollection({
      loader: glob({ pattern: '**/*.md', base: `./src/content/${locale}/entries` }),
      schema: entrySchema,
    }),
  ])
);
```

`entrySchema` fields:

```typescript
z.object({
  name_vi: z.string(),                          // Vietnamese name (required)
  name_han: z.string().optional(),               // Hán tự name
  aliases: z.array(z.string()).optional(),        // Alternative names
  name_en: z.string().optional(),                // English name
  category: z.string(),                          // Category slug (required)
  subcategories: z.array(z.string()).optional(),
  gender: z.string().optional(),                 // "nam" | "nu" | "khong-xac-dinh"
  era: z.string().optional(),                    // Historical era text
  year_approx: z.number().optional(),            // Approximate year (negative = BCE)
  year_end: z.number().optional(),
  region: z.string().optional(),                 // "bac" | "trung" | "nam"
  locations: z.array(z.string()).optional(),      // Place names
  coordinates: z.array(z.number()).optional(),    // [lat, lng]
  relations: z.object({
    family: z.array(z.string()).optional(),
    teachers: z.array(z.string()).optional(),
    allies: z.array(z.string()).optional(),
    cohabitors: z.array(z.string()).optional(),
    allied_historical: z.array(z.string()).optional(),
    enemies: z.array(z.string()).optional(),
    artifacts: z.array(z.string()).optional(),
    mythic_events: z.array(z.string()).optional(),
    historic_events: z.array(z.string()).optional(),
    related_sites: z.array(z.string()).optional(),
  }).optional(),
  sources: z.array(z.object({
    title: z.string(),
    author: z.string().optional(),
    chapter: z.string().optional(),
    edition: z.string().optional(),
  })).optional(),
  summary: z.string().optional(),                // Short description
  group: z.string().optional(),                  // Grouping label (e.g. "Tứ Bất Tử")
  themes: z.array(z.string()).optional(),         // Theme tags (slugs)
  popularity: z.number().default(1),             // Sorting weight
  status: z.string().default('published'),       // "published" = visible
  author: z.string().optional(),
  updated_at: z.coerce.string().optional(),
})
```

## Category System

File: `src/data/category-labels.ts`

| Slug | VI label | EN label (via `getCategoryLabel`) |
|------|----------|-------------------------------------|
| `than-linh` | Thần linh | Deities |
| `anh-hung` | Anh hùng | Heroes |
| `yeu-quai` | Yêu quái | Demons |
| `linh-vat` | Linh vật | Sacred Beasts |
| `dia-danh` | Địa danh | Places |
| `vat-pham` | Vật phẩm | Artifacts |
| `le-hoi` | Lễ hội | Festivals |
| `tich-co` | Tích cổ | Ancient Tales |

`CATEGORY_SLUGS` = array of all slug keys. Used for `getStaticPaths()` in category pages.

## relations.* strings and the graph

Each value is a **free-text string** (often with a qualifier in parentheses, e.g. `Lạc Long Quân (chồng)`). At build time, `src/lib/relations-graph.ts` resolves names to entry IDs using normalized `name_vi`, `name_en`, and `aliases` (exact match). Strings that do not match any entry become **ghost nodes** on the relation graph. See [relations-graph.md](./relations-graph.md).

## Region & gender labels (UI)

Region and gender **slugs** in frontmatter are fixed; **display strings** for the sidebar and filters come from `t(lang, 'region.*')` and `t(lang, 'gender.*')` in `src/i18n/config.ts`, not only from hardcoded Vietnamese tables.

## Frontmatter Example

```yaml
---
name_vi: Thánh Gióng
name_han: 聖揀
aliases:
  - Phù Đổng Thiên Vương
  - Sóc Thiên Vương
name_en: Saint Gióng
category: anh-hung
gender: nam
era: Hùng Vương thứ 6
year_approx: -1718
region: bac
locations:
  - Phù Đổng
coordinates:
  - 21.0667
  - 105.9833
relations:
  family:
    - Mẹ Gióng
  allies:
    - Hùng Vương thứ 6
  enemies:
    - Giặc Ân
  artifacts:
    - Ngựa sắt
    - Roi sắt
sources:
  - title: Lĩnh Nam Chích Quái
    author: Trần Thế Pháp
    chapter: Truyện Đổng Thiên Vương
summary: "Cậu bé ba tuổi không nói không cười..."
group: Tứ Bất Tử
themes:
  - chien-tranh
  - bao-ve-to-quoc
popularity: 5
status: published
author: claude+admin
updated_at: 2026-04-09
---

## Câu _chuyện_

Markdown body here...
```

## Content Conventions

- **File naming**: `kebab-case.md` matching the entry's ID (e.g. `thanh-giong.md`); use the **same** `id` in `vi` and `en` when both exist
- **Markdown headings**: Use `## Heading _italic_` style — the italic part gets styled differently via `EntryLayout.astro`
- **Only `status: 'published'`** entries appear on the site
- **`popularity`** drives sort order (higher = shown first) and "related entries" selection
- **Relations** are plain text strings, not linked to other entry IDs (yet)
- **Themes** are slug strings, displayed via `slugToLabel()` which replaces `-` with spaces

## Existing Entries

Markdown files live under `src/content/{locale}/entries/`. Regenerate this table from the repo when inventory matters for agents:

| File | Name | Category | Popularity |
|------|------|----------|-----------|
| `au-co.md` | Âu Cơ | — | — |
| `ho-tinh.md` | Hồ Tinh | yeu-quai | — |
| `lac-long-quan.md` | Lạc Long Quân | than-linh | 10 |
| `moc-tinh.md` | Mộc Tinh | — | — |
| `ngu-tinh.md` | Ngư Tinh | — | — |
| `thanh-giong.md` | Thánh Gióng | anh-hung | 5 |
