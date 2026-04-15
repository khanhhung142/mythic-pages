# Data Flow

All data is resolved at **build time**. There is no runtime API, no database, no client-side state.

## Content Pipeline Overview

```mermaid
graph TD
    A["Markdown files<br/>src/content/entries/*.md"] -->|"glob loader"| B["Astro Content Collections"]
    B -->|"Zod validation"| C{"Valid?"}
    C -->|"Yes"| D["In-memory collection"]
    C -->|"No"| E["Build error"]
    D -->|"getCollection()"| F["Page frontmatter scripts"]
    F -->|"filter/sort/slice"| G["Template data"]
    G -->|"Astro rendering"| H["Static HTML in dist/"]
```

## Step-by-Step Flow

### 1. Content Loading

```
src/content/entries/*.md
        ↓
glob({ pattern: '**/*.md', base: './src/content/entries' })
        ↓
Each .md file parsed: YAML frontmatter → data, body → markdown
        ↓
entry.id = filename without .md (e.g. "thanh-giong")
```

### 2. Schema Validation

```
entry.data (frontmatter) → Zod schema validation
        ↓
Required fields: name_vi, category
Defaults applied: popularity=1, status='published'
Optional fields: all others
        ↓
Type-safe entry object: CollectionEntry<'entries'>
```

### 3. Page Data Resolution

Each page fetches data differently:

```mermaid
graph LR
    subgraph "Home Page"
        H1["getCollection('entries', published)"] --> H2["Fisher-Yates shuffle"] --> H3["featured = [0]<br/>sideEntries = [1..3]"]
    end

    subgraph "Catalog Page"
        C1["getCollection('entries', published)"] --> C2["sort by popularity desc,<br/>then name_vi asc (vi locale)"] --> C3["pass all to EntriesListPage"]
    end

    subgraph "Category Page"
        F1["getCollection('entries', published)"] --> F2["sort (same as catalog)"] --> F3["filter by category"] --> F4["pass filtered to EntriesListPage"]
    end

    subgraph "Entry Detail"
        E1["getCollection('entries')"] --> E2["filter published"] --> E3["getStaticPaths():<br/>one path per entry"]
        E2 --> E4["related = top 3 by popularity<br/>(excluding current)"]
        E3 --> E5["render(entry) → Content component"]
    end
```

### 4. Sort Logic

Used in catalog and category pages:

```typescript
entries.sort((a, b) => {
  // Primary: popularity descending
  if ((b.data.popularity ?? 0) !== (a.data.popularity ?? 0)) {
    return (b.data.popularity ?? 0) - (a.data.popularity ?? 0);
  }
  // Secondary: name_vi ascending (Vietnamese locale)
  return (a.data.name_vi ?? '').localeCompare(b.data.name_vi ?? '', 'vi');
});
```

### 5. Related Entries Logic

In `[id].astro` → `getStaticPaths()`:

```typescript
related: published
  .filter(e => e.id !== entry.id)        // exclude current
  .sort((a, b) => (b.data.popularity ?? 1) - (a.data.popularity ?? 1))  // by popularity
  .slice(0, 3)                           // top 3
```

### 6. Markdown Rendering

```typescript
const { Content } = await render(entry);
// Content is an Astro component that renders the markdown body
// Passed as <slot> to EntryLayout
```

The rendered markdown receives typography from `EntryLayout`'s `.entry-content` styles.

## Category Labels Resolution

```mermaid
graph LR
    A["entry.data.category<br/>(slug: 'anh-hung')"] --> B["CATEGORY_LABELS[slug]<br/>(map lookup)"]
    B --> C["Vietnamese label<br/>('Anh hùng')"]
    B -->|"not found"| D["fallback: raw slug"]
```

Used in: `EntriesListPage` (tags), `EntryLayout` (breadcrumb, info table), `index.astro` (featured)

## Data Flow per Page

| Page | Input | Transform | Output |
|------|-------|-----------|--------|
| `/` | All published entries | Shuffle → take first 4 | Featured card + 3 side cards |
| `/entries` | All published entries | Sort by popularity/name | Full card grid |
| `/entries/category/X` | All published entries | Sort → filter by category | Filtered card grid |
| `/entries/Y` | Single entry + all published | render() + top 3 related | Full article + sidebar + related |

## Key Gotcha

The home page uses **random shuffle** (`Math.random()`) — so featured entries change on every build. This is intentional for variety but means builds are non-deterministic.
