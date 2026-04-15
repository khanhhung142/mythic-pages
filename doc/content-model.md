# Content Model

All content lives as Markdown files in `src/content/entries/`. Each file = one mythological entry.

## Schema Definition

File: `src/content.config.ts`

```typescript
const entries = defineCollection({
  loader: glob({ pattern: '**/*.md', base: './src/content/entries' }),
  schema: z.object({
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
      allies: z.array(z.string()).optional(),
      enemies: z.array(z.string()).optional(),
      artifacts: z.array(z.string()).optional(),
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
  }),
});
```

## Category System

File: `src/data/category-labels.ts`

| Slug | Vietnamese Label | Meaning |
|------|-----------------|---------|
| `than-linh` | Thần linh | Deities |
| `anh-hung` | Anh hùng | Heroes |
| `yeu-quai` | Yêu quái | Demons |
| `linh-vat` | Linh vật | Sacred beasts |
| `dia-danh` | Địa danh | Places |
| `vat-pham` | Vật phẩm | Artifacts |
| `le-hoi` | Lễ hội | Festivals |
| `tich-co` | Tích cổ | Ancient tales |

`CATEGORY_SLUGS` = array of all slug keys. Used for `getStaticPaths()` in category pages.

## Region Labels (in EntryLayout)

| Slug | Label |
|------|-------|
| `bac` | Bắc Bộ |
| `trung` | Trung Bộ |
| `nam` | Nam Bộ |

## Gender Labels (in EntryLayout)

| Slug | Label |
|------|-------|
| `nam` | Nam |
| `nu` | Nữ |
| `khong-xac-dinh` | Không xác định |

## Frontmatter Example

```yaml
---
id: thanh-giong
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

- **File naming**: `kebab-case.md` matching the entry's ID (e.g. `thanh-giong.md`)
- **Markdown headings**: Use `## Heading _italic_` style — the italic part gets styled differently via `EntryLayout.astro`
- **Only `status: 'published'`** entries appear on the site
- **`popularity`** drives sort order (higher = shown first) and "related entries" selection
- **Relations** are plain text strings, not linked to other entry IDs (yet)
- **Themes** are slug strings, displayed via `slugToLabel()` which replaces `-` with spaces

## Existing Entries

| File | Name | Category | Popularity |
|------|------|----------|-----------|
| `au-co.md` | Âu Cơ | — | — |
| `ho-tinh.md` | Hồ Tinh | yeu-quai | — |
| `lac-long-quan.md` | Lạc Long Quân | than-linh | 10 |
| `moc-tinh.md` | Mộc Tinh | — | — |
| `ngu-tinh.md` | Ngư Tinh | — | — |
| `thanh-giong.md` | Thánh Gióng | anh-hung | 5 |
