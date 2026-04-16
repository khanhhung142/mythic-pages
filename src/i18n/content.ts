import { getCollection } from 'astro:content';
import type { Locale } from './config';

export async function getLocalizedEntries(locale: Locale) {
  const collectionName = locale === 'vi' ? 'entriesVi' : 'entriesEn';

  const entries = await getCollection(collectionName, (e) => e.data.status === 'published');

  if (locale === 'vi') return entries;

  // Fallback: for missing EN entries, include VI versions
  const viEntries = await getCollection('entriesVi', (e) => e.data.status === 'published');
  const enIds = new Set(entries.map((e) => e.id));
  const fallbacks = viEntries.filter((e) => !enIds.has(e.id));

  return [...entries, ...fallbacks];
}

export async function getLocalizedEntry(locale: Locale, id: string) {
  const collectionName = locale === 'vi' ? 'entriesVi' : 'entriesEn';

  const entries = await getCollection(collectionName);
  let entry = entries.find((e) => e.id === id);

  if (!entry && locale !== 'vi') {
    const viEntries = await getCollection('entriesVi');
    entry = viEntries.find((e) => e.id === id);
  }

  return entry ?? null;
}

export async function getAllEntryIds() {
  const viEntries = await getCollection('entriesVi', (e) => e.data.status === 'published');
  return viEntries.map((e) => e.id);
}
