# Data Flow

All data is resolved at **build time**. There is no runtime API, no database, no client-side state.

## Content Pipeline Overview

```mermaid
graph TD
    A["Markdown files<br/>src/content/{locale}/entries/*.md"] -->|"glob loaders"| B["Astro Content Collections<br/>entries{Locale}"]
    B -->|"Zod validation"| C{"Valid?"}
    C -->|"Yes"| D["In-memory collections"]
    C -->|"No"| E["Build error"]
    D -->|"getLocalizedEntries / getLocalizedEntry"| F["src/i18n/content.ts"]
    F --> G["Page frontmatter scripts<br/>src/pages/... and [lang]/..."]
    G -->|"filter/sort/slice"| H["Template data"]
    H -->|"Astro rendering"| I["Static HTML in dist/"]
```

## Step-by-Step Flow

### 1. Content Loading

```
src/content/{locale}/entries/*.md   → collection entries{Locale}
        ↓
glob({ pattern: '**/*.md', base: './src/content/{locale}/entries' })
        ↓
Each .md file parsed: YAML frontmatter → data, body → markdown
        ↓
entry.id = filename without .md (e.g. "thanh-giong")
```

### 2. Schema Validation

```
entry.data (frontmatter) → Zod schema validation (shared entrySchema)
        ↓
Required fields: name_vi, category
Defaults applied: popularity=1, status='published'
Optional fields: all others
        ↓
Type-safe entry objects: CollectionEntry<`entries${Capitalize<Locale>}`>
```

### 3. Locale fallback (generic)

For **any non-default locale**, `getLocalizedEntries`:

1. Loads published entries from `entries{Locale}`
2. Loads published entries from default-locale collection
3. Appends default-locale entries whose `id` is missing in requested locale

So localized catalogs and static paths include every story from the canonical default locale even before translation exists.

### 4. Page Data Resolution

```mermaid
graph LR
    subgraph "Home Page (index / HomePage)"
        H1["getLocalizedEntries(lang)"] --> H2["Fisher-Yates shuffle"] --> H3["featured = [0]<br/>sideEntries = [1..3]"]
    end

    subgraph "Catalog Page"
        C1["getLocalizedEntries(lang)"] --> C2["sort by popularity desc,<br/>then name_vi asc (vi locale)"] --> C3["pass all to EntriesListPage + lang"]
    end

    subgraph "Category Page"
        F1["getLocalizedEntries(lang)"] --> F2["sort (same as catalog)"] --> F3["filter by category"] --> F4["pass filtered to EntriesListPage + lang"]
    end

    subgraph "Entry Detail"
        E1["getLocalizedEntry(lang, id)"] --> E2["getStaticPaths():<br/>paths per locale + id"]
        E2 --> E3["related = top 3 by popularity<br/>(excluding current)"]
        E3 --> E4["render(entry) → Content component"]
    end
```

### 5. Sort Logic

Used in catalog and category pages:

```typescript
entries.sort((a, b) => {
  if ((b.data.popularity ?? 0) !== (a.data.popularity ?? 0)) {
    return (b.data.popularity ?? 0) - (a.data.popularity ?? 0);
  }
  return (a.data.name_vi ?? '').localeCompare(b.data.name_vi ?? '', 'vi');
});
```

### 6. Related Entries Logic

In `entries/[id].astro` and `[lang]/entries/[id].astro` → `getStaticPaths()` / props:

```typescript
related: published
  .filter(e => e.id !== entry.id)
  .sort((a, b) => (b.data.popularity ?? 1) - (a.data.popularity ?? 1))
  .slice(0, 3)
```

### 7. Markdown Rendering

```typescript
const { Content } = await render(entry);
```

The rendered markdown receives typography from `EntryLayout`'s `.entry-content` styles.

## Category Labels Resolution

```mermaid
graph LR
    A["entry.data.category<br/>(slug: 'anh-hung')"] --> B["getCategoryLabel(slug, lang)"]
    B --> C["CATEGORY_LABELS[slug][lang]"]
    B -->|"missing"| D["fallback: vi label, then slug"]
```

Used in: `EntriesListPage`, `EntryLayout`, `HomePage`

## Data Flow per Page

| Page | Input | Transform | Output |
|------|-------|-----------|--------|
| `/`, `/en/` | Localized published entries | Shuffle → take first 4 | Featured card + 3 side cards |
| `/entries`, `/{lang}/entries` | Localized entries | Sort by popularity/name | Full card grid |
| `/.../entries/category/X` | Localized entries | Sort → filter by category | Filtered card grid |
| `/.../entries/Y` | `getLocalizedEntry` + published | render() + top 3 related | Full article + sidebar + related |

## Key Gotcha

The home page uses **random shuffle** (`Math.random()`) — so featured entries change on every build. This is intentional for variety but means builds are non-deterministic.
