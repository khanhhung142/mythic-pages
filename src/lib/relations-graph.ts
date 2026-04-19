import type { Locale } from '../i18n/config';
import { localePath } from '../i18n/paths';

export type RelationKind =
  | 'family'
  | 'teachers'
  | 'allies'
  | 'cohabitors'
  | 'allied_historical'
  | 'enemies'
  | 'artifacts'
  | 'mythic_events'
  | 'historic_events'
  | 'related_sites';

export type GraphNode = {
  id: string;
  name: string;
  name_en?: string;
  han?: string;
  category?: string;
  isGhost: boolean;
  href?: string;
};

export type GraphEdge = {
  source: string;
  target: string;
  kind: RelationKind;
  qualifier?: string;
};

export type BuiltGraph = {
  nodes: GraphNode[];
  edges: GraphEdge[];
  unresolved: string[];
  ghostCount: number;
};

export type EntryLike = {
  id: string;
  data: {
    name_vi: string;
    name_en?: string;
    name_han?: string;
    aliases?: string[];
    category?: string;
    relations?: Partial<Record<RelationKind, string[]>>;
    status?: string;
  };
};

/** All relation kinds (order for filters / mini-graph priority). */
export const RELATION_KINDS: readonly RelationKind[] = [
  'family',
  'teachers',
  'allies',
  'cohabitors',
  'allied_historical',
  'enemies',
  'artifacts',
  'mythic_events',
  'historic_events',
  'related_sites',
] as const;

const RELATION_KEYS: { key: RelationKind; field: keyof NonNullable<EntryLike['data']['relations']> }[] = [
  { key: 'family', field: 'family' },
  { key: 'teachers', field: 'teachers' },
  { key: 'allies', field: 'allies' },
  { key: 'cohabitors', field: 'cohabitors' },
  { key: 'allied_historical', field: 'allied_historical' },
  { key: 'enemies', field: 'enemies' },
  { key: 'artifacts', field: 'artifacts' },
  { key: 'mythic_events', field: 'mythic_events' },
  { key: 'historic_events', field: 'historic_events' },
  { key: 'related_sites', field: 'related_sites' },
];

/** Strip parentheticals and normalize for exact matching (name_vi / name_en / aliases). */
export function normalizeForMatch(s: string): string {
  return s
    .replace(/\([^)]*\)/g, '')
    .replace(/[·,—–-]/g, ' ')
    .normalize('NFC')
    .trim()
    .toLowerCase()
    .replace(/\s+/g, ' ');
}

export function extractQualifier(raw: string): string | undefined {
  const m = raw.match(/\(([^)]+)\)/);
  return m ? m[1].trim() : undefined;
}

function ghostId(raw: string): string {
  const base = normalizeForMatch(raw) || raw;
  let h = 2166136261;
  for (let i = 0; i < base.length; i++) {
    h ^= base.charCodeAt(i);
    h = Math.imul(h, 16777619);
  }
  return `ghost::${(h >>> 0).toString(16)}`;
}

function displayName(entry: EntryLike['data'], lang: Locale): string {
  if (lang === 'en') {
    return entry.name_en?.trim() || entry.name_vi;
  }
  return entry.name_vi;
}

function buildNameToIdMap(entries: EntryLike[]): Map<string, string> {
  const map = new Map<string, string>();
  for (const e of entries) {
    const keys = [
      normalizeForMatch(e.data.name_vi),
      ...(e.data.name_en ? [normalizeForMatch(e.data.name_en)] : []),
      ...(e.data.aliases?.map((a) => normalizeForMatch(a)) ?? []),
    ].filter(Boolean);
    for (const k of keys) {
      if (!map.has(k)) map.set(k, e.id);
      else if (map.get(k) !== e.id) {
        console.warn(
          `[relations-graph] duplicate name key "${k}" → keep ${map.get(k)}, skip ${e.id}`
        );
      }
    }
  }
  return map;
}

function resolveTarget(
  raw: string,
  nameToId: Map<string, string>
): { id: string; ghost: boolean } {
  const n = normalizeForMatch(raw);
  if (!n) return { id: ghostId(raw), ghost: true };
  const id = nameToId.get(n);
  if (id) return { id, ghost: false };
  return { id: ghostId(raw), ghost: true };
}

function edgeKey(source: string, target: string, kind: RelationKind): string {
  return `${source}|${target}|${kind}`;
}

const KIND_PRIORITY: Record<RelationKind, number> = {
  family: 0,
  teachers: 1,
  artifacts: 2,
  allies: 3,
  cohabitors: 4,
  allied_historical: 5,
  enemies: 6,
  mythic_events: 7,
  historic_events: 8,
  related_sites: 9,
};

/**
 * Full graph: all published entries as nodes; edges from relations.* strings.
 */
export function buildGraph(entries: EntryLike[], lang: Locale): BuiltGraph {
  const published = entries.filter((e) => (e.data.status ?? 'published') === 'published');
  const nameToId = buildNameToIdMap(published);

  const nodeMap = new Map<string, GraphNode>();
  for (const e of published) {
    nodeMap.set(e.id, {
      id: e.id,
      name: displayName(e.data, lang),
      name_en: e.data.name_en,
      han: e.data.name_han?.charAt(0),
      category: e.data.category,
      isGhost: false,
      href: localePath(lang, `/entries/${e.id}`),
    });
  }

  const edgeMap = new Map<string, GraphEdge>();
  const unresolved: string[] = [];

  for (const e of published) {
    const rel = e.data.relations;
    if (!rel) continue;

    for (const { key, field } of RELATION_KEYS) {
      const items = rel[field];
      if (!items?.length) continue;

      for (const raw of items) {
        const qualifier = extractQualifier(raw);
        const { id: targetId, ghost } = resolveTarget(raw, nameToId);
        if (ghost) unresolved.push(raw);

        if (targetId === e.id) continue;

        if (ghost && !nodeMap.has(targetId)) {
          nodeMap.set(targetId, {
            id: targetId,
            name: raw.replace(/\([^)]*\)/g, '').trim() || raw,
            isGhost: true,
          });
        }

        const ek = edgeKey(e.id, targetId, key);
        const existing = edgeMap.get(ek);
        if (existing) {
          if (qualifier) {
            const prev = existing.qualifier ?? '';
            existing.qualifier = [prev, qualifier].filter(Boolean).join(' · ');
          }
        } else {
          edgeMap.set(ek, {
            source: e.id,
            target: targetId,
            kind: key,
            qualifier,
          });
        }
      }
    }
  }

  const ghostCount = [...nodeMap.values()].filter((n) => n.isGhost).length;

  return {
    nodes: [...nodeMap.values()],
    edges: [...edgeMap.values()],
    unresolved,
    ghostCount,
  };
}

/**
 * 1-hop neighborhood around centerId, capped at maxNodes (priority by relation kind).
 */
export function buildLocalSubgraph(
  entries: EntryLike[],
  lang: Locale,
  centerId: string,
  maxNodes = 25
): BuiltGraph {
  const full = buildGraph(entries, lang);
  const incident = full.edges.filter((e) => e.source === centerId || e.target === centerId);
  if (incident.length === 0) {
    const center = full.nodes.find((n) => n.id === centerId);
    return {
      nodes: center ? [center] : [],
      edges: [],
      unresolved: [],
      ghostCount: center?.isGhost ? 1 : 0,
    };
  }

  const sorted = [...incident].sort((a, b) => KIND_PRIORITY[a.kind] - KIND_PRIORITY[b.kind]);

  const keptIds = new Set<string>([centerId]);
  const keptEdges: GraphEdge[] = [];

  for (const e of sorted) {
    const other = e.source === centerId ? e.target : e.source;
    if (!keptIds.has(other) && keptIds.size >= maxNodes) continue;
    keptIds.add(other);
    keptEdges.push(e);
  }

  const nodes = full.nodes.filter((n) => keptIds.has(n.id));
  const ghostCount = nodes.filter((n) => n.isGhost).length;

  return {
    nodes,
    edges: keptEdges,
    unresolved: [],
    ghostCount,
  };
}
