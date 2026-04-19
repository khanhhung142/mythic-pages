import { drag } from 'd3-drag';
import {
  forceCenter,
  forceCollide,
  forceLink,
  forceManyBody,
  forceSimulation,
} from 'd3-force';
import { pointer, select, type Selection } from 'd3-selection';
import { zoom, zoomIdentity } from 'd3-zoom';
import type { Locale } from '../i18n/config';
import {
  RELATION_KINDS,
  type GraphEdge,
  type GraphNode,
  type RelationKind,
} from '../lib/relations-graph';

type SimNode = GraphNode & { x?: number; y?: number; fx?: number | null; fy?: number | null };
type SimLink = GraphEdge & { source: SimNode | string; target: SimNode | string };

export type GraphMountPayload = {
  nodes: GraphNode[];
  edges: GraphEdge[];
  ghostCount: number;
  mode: 'global' | 'local';
  lang?: Locale;
};

const NODE_R = 15;
const HAN_DY = 26;
const NAME_DY = 28;

const STROKE: Record<RelationKind, string> = {
  family: 'var(--vermilion)',
  teachers: 'var(--ink)',
  allies: '#b8860b',
  cohabitors: 'var(--ink-soft)',
  allied_historical: 'var(--ink)',
  enemies: '#6b2c1a',
  artifacts: 'var(--ink)',
  mythic_events: 'var(--ink-light)',
  historic_events: 'var(--ink-light)',
  related_sites: 'var(--ink-mute)',
};

function categoryFill(cat: string | undefined): string {
  if (!cat) return 'var(--graph-node-default-fill)';
  return `var(--graph-cat-${cat}-fill, var(--graph-node-default-fill))`;
}

function categoryStroke(cat: string | undefined): string {
  if (!cat) return 'var(--graph-node-default-stroke)';
  return `var(--graph-cat-${cat}-stroke, var(--graph-node-default-stroke))`;
}

/**
 * Zoom/pan so the whole graph fits in the SVG viewBox on first load.
 * Uses the same convention as typical d3 "fit to box": center of bounds → viewport center.
 */
function computeFitTransform(
  nodes: SimNode[],
  width: number,
  height: number
) {
  const pad = 72;
  if (nodes.length === 0) {
    return zoomIdentity;
  }

  let minX = Infinity;
  let minY = Infinity;
  let maxX = -Infinity;
  let maxY = -Infinity;

  for (const n of nodes) {
    const x = n.x ?? 0;
    const y = n.y ?? 0;
    const top = n.han ? y - NODE_R - HAN_DY : y - NODE_R;
    const bottom = y + NODE_R + NAME_DY;
    minX = Math.min(minX, x - NODE_R);
    maxX = Math.max(maxX, x + NODE_R);
    minY = Math.min(minY, top);
    maxY = Math.max(maxY, bottom);
  }

  const rawW = maxX - minX;
  const rawH = maxY - minY;
  const cx = (minX + maxX) / 2;
  const cy = (minY + maxY) / 2;

  const graphW = Math.max(rawW, 28);
  const graphH = Math.max(rawH, 28);

  const innerW = Math.max(width - 2 * pad, 100);
  const innerH = Math.max(height - 2 * pad, 100);
  let k = Math.min(innerW / graphW, innerH / graphH);
  k = Math.min(Math.max(k, 0.18), 3.2);

  return zoomIdentity.translate(width / 2, height / 2).scale(k).translate(-cx, -cy);
}

function layout(
  nodes: SimNode[],
  edges: GraphEdge[],
  width: number,
  height: number
): { simNodes: SimNode[]; simLinks: SimLink[] } {
  const simNodes = nodes.map((n) => ({ ...n }));
  const idToSim = new Map(simNodes.map((n) => [n.id, n]));
  const simLinks = edges.map((e) => ({
    ...e,
    source: idToSim.get(e.source as string) ?? e.source,
    target: idToSim.get(e.target as string) ?? e.target,
  })) as SimLink[];

  const sim = forceSimulation(simNodes as SimNode[])
    .force(
      'link',
      forceLink(simLinks)
        .id((d: unknown) => (d as SimNode).id)
        .distance(138)
        .strength(0.55)
    )
    .force('charge', forceManyBody().strength(-340))
    .force('center', forceCenter(width / 2, height / 2))
    .force('collide', forceCollide(NODE_R * 3.8));

  for (let i = 0; i < 320; i++) sim.tick();
  sim.stop();

  return { simNodes, simLinks };
}

function endpointId(end: SimLink['source'] | SimLink['target']): string {
  if (end == null) return '';
  if (typeof end === 'string') return end;
  return (end as SimNode).id;
}

function catOpacity(d: SimNode, activeCategory: string | null): number {
  if (!activeCategory) return 1;
  if (d.category && d.category !== activeCategory) return 0.15;
  if (!d.category) return 0.15;
  return 1;
}

function applyVisibility(
  root: ReturnType<typeof select>,
  visibleKinds: Set<string>,
  activeCategory: string | null
) {
  const showKind = (k: string) => visibleKinds.has(k);

  root.selectAll<SVGLineElement, SimLink>('.edge-line').each(function (d) {
    select(this).style('display', showKind(d.kind) ? null : 'none');
  });
  root.selectAll<SVGLineElement, SimLink>('.edge-hit').each(function (d) {
    select(this).style('display', showKind(d.kind) ? null : 'none');
  });

  const visibleNodeIds = new Set<string>();
  root.selectAll<SVGLineElement, SimLink>('.edge-line').each(function (d) {
    if (select(this).style('display') === 'none') return;
    visibleNodeIds.add(endpointId(d.source));
    visibleNodeIds.add(endpointId(d.target));
  });

  root.selectAll<SVGGElement, SimNode>('.node-group').each(function (d) {
    const show = visibleNodeIds.has(d.id);
    const el = select(this);
    el.style('display', show ? null : 'none');
    if (show) el.style('opacity', String(catOpacity(d, activeCategory)));
  });

  root.selectAll<SVGTextElement, SimNode>('.node-label').each(function (d) {
    const show = visibleNodeIds.has(d.id);
    const el = select(this);
    el.style('display', show ? null : 'none');
    if (show) el.style('opacity', String(catOpacity(d, activeCategory)));
  });
}

function updateLinks(
  linkLines: Selection<SVGLineElement, SimLink, SVGElement, unknown>,
  hitLines: Selection<SVGLineElement, SimLink, SVGElement, unknown>
) {
  linkLines
    .attr('x1', (d: SimLink) => (d.source as SimNode).x ?? 0)
    .attr('y1', (d: SimLink) => (d.source as SimNode).y ?? 0)
    .attr('x2', (d: SimLink) => (d.target as SimNode).x ?? 0)
    .attr('y2', (d: SimLink) => (d.target as SimNode).y ?? 0);
  hitLines
    .attr('x1', (d: SimLink) => (d.source as SimNode).x ?? 0)
    .attr('y1', (d: SimLink) => (d.source as SimNode).y ?? 0)
    .attr('x2', (d: SimLink) => (d.target as SimNode).x ?? 0)
    .attr('y2', (d: SimLink) => (d.target as SimNode).y ?? 0);
}

export function mountGraphWithPayload(
  svgId: string,
  payload: GraphMountPayload,
  options?: { filterRoot?: string | HTMLElement }
): void {
  const svg = select(`#${svgId}`);
  if (svg.empty()) return;

  const width = 1200;
  const height = payload.mode === 'local' ? 420 : 800;

  const { simNodes, simLinks } = layout(
    payload.nodes as SimNode[],
    payload.edges,
    width,
    height
  );

  const fitTransform = computeFitTransform(simNodes, width, height);

  svg.attr('viewBox', `0 0 ${width} ${height}`).attr('role', 'img');

  svg.selectAll('*').remove();

  const root = svg.append('g').attr('class', 'graph-viewport');

  let currentTransform = fitTransform;
  const z = zoom<SVGSVGElement, unknown>()
    .scaleExtent([0.12, 8])
    .on('zoom', (ev) => {
      currentTransform = ev.transform;
      root.attr('transform', ev.transform.toString());
    });

  svg.call(z as never);

  const edgesG = root.append('g').attr('class', 'edges');

  const linkLines = edgesG
    .selectAll<SVGLineElement, SimLink>('line.edge-line')
    .data(simLinks)
    .join('line')
    .attr('class', (d) => `edge-line edge-kind-${d.kind}`)
    .attr('stroke', (d) => STROKE[d.kind] ?? 'var(--ink-soft)')
    .attr('stroke-opacity', 0.45)
    .attr('stroke-width', 1)
    .attr('pointer-events', 'none');

  const hitLines = edgesG
    .selectAll<SVGLineElement, SimLink>('line.edge-hit')
    .data(simLinks)
    .join('line')
    .attr('class', (d) => `edge-hit edge-kind-${d.kind}`)
    .attr('stroke', 'transparent')
    .attr('stroke-width', 14)
    .style('cursor', 'default');

  updateLinks(linkLines, hitLines);

  const nodesG = root.append('g').attr('class', 'nodes');
  const labelsG = root.append('g').attr('class', 'labels');

  const nodeGroups = nodesG
    .selectAll<SVGGElement, SimNode>('g')
    .data(simNodes)
    .join('g')
    .attr('class', 'node-group')
    .style('cursor', (d) => (d.href ? 'pointer' : 'default'));

  const nodeCircles = nodeGroups
    .append('circle')
    .attr('r', NODE_R)
    .attr('fill', (d) => (d.isGhost ? 'none' : categoryFill(d.category)))
    .attr('stroke', (d) => (d.isGhost ? 'var(--graph-node-ghost-stroke)' : categoryStroke(d.category)))
    .attr('stroke-width', (d) => (d.isGhost ? 1.5 : 0))
    .attr('stroke-dasharray', (d) => (d.isGhost ? '2 3' : null))
    .attr('cx', (d) => d.x ?? 0)
    .attr('cy', (d) => d.y ?? 0);

  function refreshGeometry() {
    updateLinks(linkLines, hitLines);
    nodeCircles.attr('cx', (n) => n.x ?? 0).attr('cy', (n) => n.y ?? 0);
    hanLabels.attr('x', (n) => n.x ?? 0).attr('y', (n) => (n.y ?? 0) - HAN_DY);
    nameLabels.attr('x', (n) => n.x ?? 0).attr('y', (n) => (n.y ?? 0) + NAME_DY);
  }

  const dragBeh = drag<SVGCircleElement, SimNode>()
    .on('start', (event) => {
      select(event.sourceEvent.target as SVGCircleElement).raise();
    })
    .on('drag', (event, d) => {
      const [px, py] = pointer(event, svg.node() as Element);
      const [x, y] = currentTransform.invert([px, py]);
      d.x = x;
      d.y = y;
      refreshGeometry();
    });

  nodeCircles.call(dragBeh);

  nodeGroups
    .filter((d) => Boolean(d.href))
    .on('click', (ev, d) => {
      ev.preventDefault();
      if (d.href) window.location.href = d.href;
    });

  const hanLabels = labelsG
    .selectAll<SVGTextElement, SimNode>('text.han')
    .data(simNodes.filter((n) => n.han))
    .join('text')
    .attr('class', 'han node-label')
    .attr('text-anchor', 'middle')
    .attr('font-size', 13)
    .attr('fill', (d) => categoryFill(d.category))
    .attr('font-family', 'Cormorant Garamond, serif')
    .text((d) => d.han ?? '')
    .attr('x', (d) => d.x ?? 0)
    .attr('y', (d) => (d.y ?? 0) - HAN_DY);

  const nameLabels = labelsG
    .selectAll<SVGTextElement, SimNode>('text.nm')
    .data(simNodes)
    .join('text')
    .attr('class', 'nm node-label')
    .attr('text-anchor', 'middle')
    .attr('font-size', 12)
    .attr('fill', 'var(--ink-soft)')
    .attr('font-family', 'Be Vietnam Pro, sans-serif')
    .text((d) => d.name)
    .attr('x', (d) => d.x ?? 0)
    .attr('y', (d) => (d.y ?? 0) + NAME_DY);

  const filterRootSel = options?.filterRoot
    ? typeof options.filterRoot === 'string'
      ? select(options.filterRoot)
      : select(options.filterRoot)
    : null;

  let visibleKinds = new Set<string>(RELATION_KINDS);
  let activeCategory: string | null = null;

  function refreshFilters() {
    applyVisibility(root as never, visibleKinds, activeCategory);
  }

  if (filterRootSel && !filterRootSel.empty() && payload.mode === 'global') {
    filterRootSel.selectAll('[data-filter-kind]').on('click', function () {
      const kind = select(this).attr('data-filter-kind');
      if (!kind) return;
      if (kind === 'all') {
        visibleKinds = new Set(RELATION_KINDS);
        filterRootSel.selectAll('[data-filter-kind]').classed('pill-active', function () {
          return select(this).attr('data-filter-kind') === 'all';
        });
        refreshFilters();
        return;
      }
      filterRootSel.select('[data-filter-kind="all"]').classed('pill-active', false);
      if (visibleKinds.size === RELATION_KINDS.length) {
        visibleKinds = new Set([kind]);
      } else {
        if (visibleKinds.has(kind)) visibleKinds.delete(kind);
        else visibleKinds.add(kind);
        if (visibleKinds.size === 0) visibleKinds = new Set(RELATION_KINDS);
      }
      filterRootSel.selectAll('[data-filter-kind]').each(function () {
        const k = select(this).attr('data-filter-kind');
        if (k === 'all') select(this).classed('pill-active', visibleKinds.size === RELATION_KINDS.length);
        else select(this).classed('pill-active', visibleKinds.has(k ?? ''));
      });
      refreshFilters();
    });

    filterRootSel.selectAll('[data-filter-category]').on('click', function () {
      const cat = select(this).attr('data-filter-category');
      if (!cat) return;
      if (cat === 'all') {
        activeCategory = null;
        filterRootSel.selectAll('[data-filter-category]').classed('pill-active', function () {
          return select(this).attr('data-filter-category') === 'all';
        });
      } else {
        activeCategory = activeCategory === cat ? null : cat;
        filterRootSel.selectAll('[data-filter-category]').classed('pill-active', function () {
          const c = select(this).attr('data-filter-category');
          return c === 'all' ? activeCategory === null : c === activeCategory;
        });
      }
      refreshFilters();
    });

    refreshFilters();
  }

  currentTransform = fitTransform;
  root.attr('transform', fitTransform.toString());
  svg.call(z.transform as never, fitTransform);
  requestAnimationFrame(() => {
    currentTransform = fitTransform;
    root.attr('transform', fitTransform.toString());
    svg.call(z.transform as never, fitTransform);
  });

  (window as unknown as { __MYTHIC_GRAPH__?: unknown }).__MYTHIC_GRAPH__ = {
    nodes: simNodes,
    links: simLinks,
    kinds: [...RELATION_KINDS],
  };
}

/** Optional: read graph JSON from `<script type="application/json" id="…">`. */
export function mountGraph(
  svgId: string,
  dataScriptId: string,
  options?: { filterRoot?: string | HTMLElement }
): void {
  const el = document.getElementById(dataScriptId);
  if (!el?.textContent) return;
  try {
    const payload = JSON.parse(el.textContent) as GraphMountPayload;
    mountGraphWithPayload(svgId, payload, options);
  } catch {
    /* invalid JSON */
  }
}
