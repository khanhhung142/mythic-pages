import { describe, it, expect } from 'vitest';
import {
  buildGraph,
  buildLocalSubgraph,
  normalizeForMatch,
  type EntryLike,
} from '../lib/relations-graph';

const auCo: EntryLike = {
  id: 'au-co',
  data: {
    name_vi: 'Âu Cơ',
    name_en: 'Âu Cơ',
    category: 'than-linh',
    status: 'published',
    relations: {},
  },
};

const lacLong: EntryLike = {
  id: 'lac-long-quan',
  data: {
    name_vi: 'Lạc Long Quân',
    name_en: 'Dragon Lord',
    category: 'than-linh',
    status: 'published',
    relations: {
      family: ['Âu Cơ (vợ)'],
    },
  },
};

describe('relations-graph', () => {
  it('matches lac-long-quan → Âu Cơ with family + qualifier vợ', () => {
    const g = buildGraph([auCo, lacLong], 'vi');
    const edge = g.edges.find(
      (e) => e.source === 'lac-long-quan' && e.target === 'au-co' && e.kind === 'family'
    );
    expect(edge).toBeDefined();
    expect(edge?.qualifier).toMatch(/vợ/);
  });

  it('creates ghost node for unknown relation string', () => {
    const ghosty: EntryLike = {
      id: 'x',
      data: {
        name_vi: 'X',
        category: 'anh-hung',
        status: 'published',
        relations: {
          family: ['Không Tồn Tại'],
        },
      },
    };
    const g = buildGraph([ghosty], 'vi');
    expect(g.nodes.some((n) => n.isGhost)).toBe(true);
    expect(g.edges.some((e) => e.target.startsWith('ghost::'))).toBe(true);
  });

  it('is idempotent: same counts for same input', () => {
    const entries = [auCo, lacLong];
    const a = buildGraph(entries, 'vi');
    const b = buildGraph(entries, 'vi');
    expect(a.nodes.length).toBe(b.nodes.length);
    expect(a.edges.length).toBe(b.edges.length);
  });

  it('normalizeForMatch strips qualifiers', () => {
    expect(normalizeForMatch('Lạc Long Quân (chồng)')).toBe(normalizeForMatch('Lạc Long Quân'));
  });

  it('buildLocalSubgraph keeps 1-hop neighbors', () => {
    const g = buildLocalSubgraph([auCo, lacLong], 'vi', 'lac-long-quan', 25);
    expect(g.nodes.some((n) => n.id === 'au-co')).toBe(true);
    expect(g.edges.length).toBeGreaterThanOrEqual(1);
  });
});
