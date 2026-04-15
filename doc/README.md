# Mythic Pages — Documentation Index

> Vietnamese mythology wiki built with Astro 6. Static site, no backend.
> This doc system is designed for AI agents to read and build new features.

## Quick Facts

| Item | Value |
|------|-------|
| Framework | Astro 6.1.5 (static SSG) |
| Language | TypeScript 5, Astro components, Markdown |
| Styling | CSS variables + Tailwind 3 + scoped `<style>` |
| Content | Astro Content Collections + Zod schemas |
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
| 3 | [routing-and-pages.md](./routing-and-pages.md) | URL routes, page files, layout hierarchy, mermaid diagrams |
| 4 | [components.md](./components.md) | Component catalog — active vs unused, props, where each is used |
| 5 | [styling.md](./styling.md) | CSS variables (design tokens), Tailwind config, font system, responsive |
| 6 | [data-flow.md](./data-flow.md) | Build-time content pipeline: MD → Zod → Collection → HTML |
| 7 | [adding-features.md](./adding-features.md) | AI playbook — how to add pages, components, content, styles |
| 8 | [known-issues.md](./known-issues.md) | Config drift, unused files, broken configs to be aware of |

## For AI Agents

When building a new feature:

1. Read `architecture.md` to understand project structure
2. Read `content-model.md` if touching content/data
3. Read `routing-and-pages.md` if adding new pages
4. Read `components.md` if creating/modifying UI
5. Read `styling.md` for design tokens and CSS conventions
6. Read `adding-features.md` for step-by-step patterns
7. Check `known-issues.md` to avoid existing pitfalls
