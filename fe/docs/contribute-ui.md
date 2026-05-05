# Contribute UI + Cloudflare Pages Functions (Turnstile + GitHub Issues)

This site is static Astro, but Cloudflare Pages supports **Pages Functions** for server-side endpoints.

Contribute flow implemented here:

- UI: `src/pages/[...lang]/contribute.astro` → `src/components/ContributePage.astro`
- API: `POST /api/contribute` → `functions/api/contribute.ts`
- Output: creates a **GitHub Issue** in a **private** repo

## Security model (no secrets in frontend)

- Frontend uses **Turnstile site key** only (public): `PUBLIC_TURNSTILE_SITE_KEY`
- All secrets live **server-side** in Pages Functions:
  - `TURNSTILE_SECRET_KEY`
  - `GITHUB_TOKEN`

Never embed secrets in `.astro` templates, JS bundles, or `PUBLIC_*` env vars.

## 1) Create Cloudflare Turnstile widget

In Cloudflare Dashboard:

- Turnstile → Create widget
- Mode: Managed
- Allowed hostnames:
  - Production domain (e.g. `vietmyth.vn`)
  - Cloudflare Pages preview domains if you test on previews

You will get:

- **Site key** (public) → `PUBLIC_TURNSTILE_SITE_KEY`
- **Secret key** (private) → `TURNSTILE_SECRET_KEY`

## 2) Create GitHub token (private repo)

Use **fine-grained personal access token** (recommended).

- Repository access: select the private repo
- Permissions:
  - **Issues: Read & Write**

No need for Pull requests/Contents because we only create Issues.

## 3) Configure Cloudflare Pages environment variables

Cloudflare Pages → your project → Settings → Environment variables.

Add **Environment variables**:

- `PUBLIC_TURNSTILE_SITE_KEY` = `<turnstile_site_key>`
- `GITHUB_REPO` = `owner/repo`

Add **Secrets**:

- `TURNSTILE_SECRET_KEY` = `<turnstile_secret_key>`
- `GITHUB_TOKEN` = `<github_fine_grained_pat>`

Make sure these are set for the environments you use (Preview + Production).

## 4) How request verification works

`functions/api/contribute.ts` does:

1. Parses JSON body
2. Checks honeypot field (`website`) — if filled, it does nothing
3. Verifies Turnstile token via `https://challenges.cloudflare.com/turnstile/v0/siteverify`
4. Calls GitHub API to create Issue:
   - `POST https://api.github.com/repos/{GITHUB_REPO}/issues`

## 5) Troubleshooting

### “Turnstile missing”

- `PUBLIC_TURNSTILE_SITE_KEY` not set in Pages env, or not available in build env.

### “turnstile_failed”

- Wrong `TURNSTILE_SECRET_KEY`
- Hostname not allowed in Turnstile widget settings (common for preview domains)

### “github_failed”

- `GITHUB_TOKEN` missing/invalid
- Token not granted access to repo
- Issues permission not set to Read & Write

### Labels not applied

Issue creation tries to add labels (`contribution`, `new-entry` / `edit-entry`).
If labels do not exist, GitHub may respond 422; function falls back to creating Issue without labels.

