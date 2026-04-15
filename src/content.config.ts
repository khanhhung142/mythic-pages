import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

const entries = defineCollection({
  loader: glob({ pattern: '**/*.md', base: './src/content/entries' }),
  schema: z.object({
    name_vi: z.string(),
    name_han: z.string().optional(),
    aliases: z.array(z.string()).optional(),
    name_en: z.string().optional(),
    category: z.string(),
    subcategories: z.array(z.string()).optional(),
    gender: z.string().optional(),
    era: z.string().optional(),
    year_approx: z.number().optional(),
    year_end: z.number().optional(),
    region: z.string().optional(),
    locations: z.array(z.string()).optional(),
    coordinates: z.array(z.number()).optional(),
    relations: z.object({
      family: z.array(z.string()).optional(),
      allies: z.array(z.string()).optional(),
      enemies: z.array(z.string()).optional(),
      artifacts: z.array(z.string()).optional(),
    }).optional(),
    sources: z.array(z.object({
      title: z.string(),
      author: z.string().optional(),
      chapter: z.string().optional(),
      edition: z.string().optional(),
    })).optional(),
    summary: z.string().optional(),
    group: z.string().optional(),
    themes: z.array(z.string()).optional(),
    popularity: z.number().default(1),
    status: z.string().default('published'),
    author: z.string().optional(),
    updated_at: z.coerce.string().optional(),
  }),
});

export const collections = { entries };
