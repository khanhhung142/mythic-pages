import { getCollection } from 'astro:content';
import { defaultLocale } from './config';
import type { Locale } from './config';

type EntryCollectionKey = `entries${Capitalize<Locale>}`;

function collectionName(locale: Locale): EntryCollectionKey {
  return `entries${locale.charAt(0).toUpperCase()}${locale.slice(1)}` as EntryCollectionKey;
}

export async function getLocalizedEntries(locale: Locale) {
  const name = collectionName(locale);

  const entries = await getCollection(name, (e) => e.data.status === 'published');

  if (locale === defaultLocale) return entries;

  // Fallback: for missing localized entries, include default-locale versions.
  const baseEntries = await getCollection(collectionName(defaultLocale), (e) => e.data.status === 'published');
  const enIds = new Set(entries.map((e) => e.id));
  const fallbacks = baseEntries.filter((e) => !enIds.has(e.id));

  return [...entries, ...fallbacks];
}

export async function getLocalizedEntry(locale: Locale, id: string) {
  const name = collectionName(locale);

  const entries = await getCollection(name);
  let entry = entries.find((e) => e.id === id);

  if (!entry && locale !== defaultLocale) {
    const baseEntries = await getCollection(collectionName(defaultLocale));
    entry = baseEntries.find((e) => e.id === id);
  }

  return entry ?? null;
}

export async function getAllEntryIds() {
  const defaultEntries = await getCollection(collectionName(defaultLocale), (e) => e.data.status === 'published');
  return defaultEntries.map((e) => e.id);
}
