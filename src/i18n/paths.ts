import { defaultLocale, locales } from './config';
import type { Locale } from './config';

export function localeStaticPaths() {
  return locales.map((locale) => ({
    params: { lang: locale === defaultLocale ? undefined : locale },
    props: { lang: locale },
  }));
}

/**
 * Build a localized site path (no origin).
 * `path` is the path without locale prefix: "/", "/about", "/entries", "/entries/foo", or "/#hash".
 */
export function localePath(lang: Locale, path: string): string {
  const hashIdx = path.indexOf('#');
  const hash = hashIdx >= 0 ? path.slice(hashIdx) : '';
  const pathOnly = hashIdx >= 0 ? path.slice(0, hashIdx) : path;
  const trimmed = pathOnly.replace(/\/+$/, '') || '/';
  const normalized = trimmed === '/' ? '' : trimmed.replace(/^\/+/, '');
  if (lang === defaultLocale) {
    const base = normalized ? `/${normalized}` : "/";
    return base + hash;
  }
  const base = normalized ? `/${lang}/${normalized}` : `/${lang}/`;
  return base + hash;
}

export function localeFromPathname(pathname: string): Locale {
  const p = pathname.replace(/\/+$/, '') || '/';
  const segment = p.split('/').filter(Boolean)[0];
  if (segment && locales.includes(segment as Locale)) return segment as Locale;
  return defaultLocale;
}

/** Same story, other locale — for Header locale links. */
export function alternateLocalePath(pathname: string, target: Locale): string {
  const raw = pathname.replace(/\/+$/, '') || '/';
  const segments = raw.split('/').filter(Boolean);
  const firstSegment = segments[0] as Locale | undefined;
  const hasLocalePrefix = Boolean(firstSegment && locales.includes(firstSegment));

  let unprefixedPath = raw;
  if (hasLocalePrefix) {
    const rest = segments.slice(1).join('/');
    unprefixedPath = rest ? `/${rest}` : '/';
  }

  if (target === defaultLocale) return unprefixedPath;
  if (unprefixedPath === '/') return `/${target}/`;
  return `/${target}${unprefixedPath}`;
}
