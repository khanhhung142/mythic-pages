import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';
import { locales } from './i18n/config';

const entrySchema = z.object({
  id: z.string().optional(),
  name_vi: z.string(),
  name_han: z.string().optional(),
  aliases: z.array(z.string()).optional(),
  name_en: z.string().optional(),
  category: z.string(),
  subcategories: z.array(z.string()).optional(),
  type: z.string().optional(),
  gender: z.string().optional(),
  era: z.string().optional(),
  era_mythic: z.string().optional(),
  era_historic: z.string().optional(),
  year_approx: z.union([z.number(), z.string(), z.null()]).optional(),
  year_note: z.string().optional(),
  year_end: z.number().optional(),
  region: z.string().optional(),
  locations: z.array(z.string()).optional(),
  location_modern: z.string().optional(),
  coordinates: z.array(z.number()).optional(),
  geography: z.record(z.string(), z.string()).optional(),
  relations: z
    .object({
      family: z.array(z.string()).optional(),
      allies: z.array(z.string()).optional(),
      enemies: z.array(z.string()).optional(),
      artifacts: z.array(z.string()).optional(),
      teachers: z.array(z.string()).optional(),
      allied_historical: z.array(z.string()).optional(),
      cohabitors: z.array(z.string()).optional(),
      mythic_events: z.array(z.string()).optional(),
      historic_events: z.array(z.string()).optional(),
      related_sites: z.array(z.string()).optional(),
    })
    .optional(),
  sources: z
    .array(
      z.object({
        title: z.string(),
        author: z.string().optional(),
        chapter: z.string().optional(),
        edition: z.string().optional(),
        notes: z.string().optional(),
      })
    )
    .optional(),
  summary: z.string().optional(),
  group: z.string().optional(),
  themes: z.array(z.string()).optional(),
  text_primary: z.string().optional(),
  cult_center: z.string().optional(),
  heritage_status: z.array(z.string()).optional(),
  archaeological_note: z.string().optional(),
  scholarly_debates: z.array(z.string()).optional(),
  popularity: z.number().default(1),
  status: z.string().default('published'),
  updated_at: z.coerce.string().optional(),
});

function toCollectionKey(locale: string) {
  return `entries${locale.charAt(0).toUpperCase()}${locale.slice(1)}`;
}

const collections = Object.fromEntries(
  locales.map((locale) => [
    toCollectionKey(locale),
    defineCollection({
      loader: glob({ pattern: '**/*.md', base: `./src/content/${locale}/entries` }),
      schema: entrySchema,
    }),
  ])
);

export { collections };
