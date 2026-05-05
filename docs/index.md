# Mythic Pages — Documentation Index

> Vietnamese mythology wiki built with Astro 6. Static site, no backend.
> This doc system is designed for AI agents to read and build new features.

## Quick Facts

| Item | Value |
|------|-------|
| Framework | Astro 6.1.5 (static SSG) |
| Language | TypeScript 5, Astro components, Markdown |
| Styling | CSS variables + Tailwind 3 + scoped `<style>` |
| Content | Astro Content Collections + Zod schemas (dynamic `entries{Locale}` from `locales`) |
| i18n | Dynamic-locale routing: default locale at root (`/`), non-default locales under `/{lang}/...`; fallback to default-locale markdown when localized file is missing |
| Backend | None — pure static site |
| Auth | None |
| State | None — all data resolved at build time |
| Deploy | Any static host (Netlify, Vercel, Cloudflare Pages, etc.) |

## Documentation Map

Read these files in order for full understanding:

| # | File | What It Covers |
|---|------|---------------|
| 1 | [architecture.md](./architecture.md) | Tech stack, project structure, build pipeline, deploy config |
| 2 | [content-model.md](./content-model.md) | Zod schema, frontmatter fields, category system, writing content |
| 3 | [routing-and-pages.md](./routing-and-pages.md) | URL routes, VI root vs `/en/`, layout hierarchy, mermaid diagrams |
| 4 | [i18n.md](./i18n.md) | Locales, collections, EN↔VI fallback, UI strings, lang switch |
| 5 | [components.md](./components.md) | Component catalog — active vs unused, props, where each is used |
| 6 | [styling.md](./styling.md) | CSS variables (design tokens), Tailwind config, font system, responsive |
| 7 | [data-flow.md](./data-flow.md) | Build-time content pipeline: MD → Zod → Collection → HTML |
| 8 | [adding-features.md](./adding-features.md) | AI playbook — how to add pages, components, content, styles |
| 9 | [relations-graph.md](./relations-graph.md) | Relation network graph: matcher, D3 island, routes, mini widget |
| 10 | [known-issues.md](./known-issues.md) | Config drift, unused files, broken configs to be aware of |

## Cloudflare deploy docs (frontend)

Cloudflare Pages / Pages Functions docs live under `fe/docs/`:

- [`fe/docs/contribute-ui.md`](../fe/docs/contribute-ui.md) — Contribute UI + Turnstile + GitHub Issues (private repo)

## For AI Agents

When building a new feature:

1. Read `architecture.md` to understand project structure
2. Read `content-model.md` if touching content/data
3. Read `routing-and-pages.md` and `i18n.md` if adding routes or locales
4. Read `components.md` if creating/modifying UI
5. Read `styling.md` for design tokens and CSS conventions
6. Read `adding-features.md` for step-by-step patterns
7. Read `relations-graph.md` if touching `relations.*` graph or `src/lib/relations-graph.ts`
8. Check `known-issues.md` to avoid existing pitfalls

## Dynamic-First Rule (Mandatory)

This repository is now **dynamic-first** for i18n and routing.

- Do not create mirrored page trees like `src/pages/en/...`, `src/pages/ja/...`, etc.
- Add locale-aware behavior by extending shared dynamic paths (`src/pages/[...lang]/...`) and shared helpers in `src/i18n/*`.
- Any new route or feature must work for **all locales in `src/i18n/config.ts`**, not only `vi`/`en`.
- When adding locale-specific content, add files under `src/content/{locale}/...`; routing must stay generic.
- If a proposed change introduces locale-specific duplication, it should be considered incorrect and refactored.
