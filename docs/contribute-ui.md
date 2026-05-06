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

Important: deploy via **Cloudflare Pages (Git integration)**. Do not deploy this repo with `wrangler deploy` as a standalone Worker — Pages env/secrets and Pages Functions routing will not apply.

### Build-time vs runtime env (agent gotcha)

- `PUBLIC_TURNSTILE_SITE_KEY` is read in `src/components/ContributePage.astro` via `import.meta.env.PUBLIC_TURNSTILE_SITE_KEY`.
  - This means it must exist in **Pages build environment** (Pages project settings), then you must **redeploy** for the HTML bundle to include it.
  - Setting this variable under **Workers → Variables and Secrets** will not affect an Astro static build on Pages.
- `TURNSTILE_SECRET_KEY`, `GITHUB_TOKEN`, `GITHUB_REPO` are read at **runtime** in Pages Functions via `context.env.*`.

## 4) How request verification works

`functions/api/contribute.ts` does:

1. Parses JSON body
2. Checks honeypot field (`website`) — if filled, it does nothing
3. Verifies Turnstile token via `https://challenges.cloudflare.com/turnstile/v0/siteverify`
4. Calls GitHub API to create Issue:
   - `POST https://api.github.com/repos/{GITHUB_REPO}/issues`

## 5) Troubleshooting

### Click “Submit” but nothing happens (CSP blocks inline script)

Fix options (pick one):

- **Option A (keep CSP, allow specific inline script)**: add hash allowlist to `script-src`
  - Add this hash (from console) into CSP `script-src`:
    - `sha256-eJGI0Ik4oYe/PKLDOt4wcN76wYs8h+Ew05pMzdY6xG8=`
- **Option B (recommended in this repo)**: avoid inline script entirely
  - Contribute JS shipped as external file: `public/scripts/contribute-form.js`
  - Page loads it via `<script src="/scripts/contribute-form.js" defer></script>`
  - This works with CSP setups that allow `script-src 'self'` (common)

### “Turnstile missing”

- `PUBLIC_TURNSTILE_SITE_KEY` not set in Pages env, or not available to the **build**.
- After changing env vars in Pages, trigger a **new deployment** (retry build / push a commit).

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
