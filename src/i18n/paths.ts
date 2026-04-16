import type { Locale } from "./config";

/**
 * Build a localized site path (no origin).
 * `path` is the path without locale prefix: "/", "/about", "/entries", "/entries/foo", or "/#hash".
 */
export function localePath(lang: Locale, path: string): string {
  const hashIdx = path.indexOf("#");
  const hash = hashIdx >= 0 ? path.slice(hashIdx) : "";
  const pathOnly = hashIdx >= 0 ? path.slice(0, hashIdx) : path;
  const trimmed = pathOnly.replace(/\/+$/, "") || "/";
  const normalized =
    trimmed === "/" ? "" : trimmed.replace(/^\/+/, "");
  if (lang === "vi") {
    const base = normalized ? `/${normalized}` : "/";
    return base + hash;
  }
  const base = normalized ? `/en/${normalized}` : "/en/";
  return base + hash;
}

export function localeFromPathname(pathname: string): Locale {
  const p = pathname.replace(/\/$/, "") || "/";
  if (p === "/en" || p.startsWith("/en/")) return "en";
  return "vi";
}

/** Same story, other locale — for Header VI/EN links. */
export function alternateLocalePath(pathname: string, target: Locale): string {
  const raw = pathname.replace(/\/$/, "") || "/";
  const isEn = raw === "/en" || raw.startsWith("/en/");
  if (target === "en") {
    if (isEn) return pathname;
    if (raw === "/") return "/en/";
    return `/en${raw}`;
  }
  if (!isEn) return pathname;
  if (raw === "/en") return "/";
  const rest = raw.slice(3);
  return rest ? `/${rest}` : "/";
}
