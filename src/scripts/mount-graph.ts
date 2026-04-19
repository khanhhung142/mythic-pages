import { drag } from 'd3-drag';
import {
  forceCenter,
  forceCollide,
  forceLink,
  forceManyBody,
  forceSimulation,
} from 'd3-force';
import { pointer, select } from 'd3-selection';
import { zoom, zoomIdentity } from 'd3-zoom';
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
};

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
        .distance(90)
        .strength(0.6)
    )
    .force('charge', forceManyBody().strength(-220))
    .force('center', forceCenter(width / 2, height / 2))
    .force('collide', forceCollide(28));

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

function updateLinks(linkLines: any, hitLines: any) {
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

  svg.attr('viewBox', `0 0 ${width} ${height}`).attr('role', 'img');

  svg.selectAll('*').remove();

  const root = svg.append('g').attr('class', 'graph-viewport');

  let currentTransform = zoomIdentity;
  const z = zoom<SVGSVGElement, unknown>().on('zoom', (ev) => {
    currentTransform = ev.transform;
    root.attr('transform', ev.transform.toString());
  });

  svg.call(z as never).call(z.transform as never, zoomIdentity);

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

  const tooltip = root
    .append('g')
    .attr('class', 'edge-tooltip')
    .style('display', 'none')
    .style('pointer-events', 'none');

  const tipBg = tooltip.append('rect').attr('rx', 4).attr('fill', 'var(--paper)').attr('stroke', 'var(--line)');
  const tipText = tooltip
    .append('text')
    .attr('font-size', 11)
    .attr('font-family', 'Be Vietnam Pro, sans-serif')
    .attr('fill', 'var(--ink-soft)');

  hitLines
    .on('mouseenter', (_ev, d) => {
      if (!d.qualifier) return;
      const sx = (d.source as SimNode).x ?? 0;
      const sy = (d.source as SimNode).y ?? 0;
      const tx = (d.target as SimNode).x ?? 0;
      const ty = (d.target as SimNode).y ?? 0;
      const mx = (sx + tx) / 2;
      const my = (sy + ty) / 2;
      tipText.text(`(${d.qualifier})`).attr('x', mx + 6).attr('y', my - 6);
      const bbox = (tipText.node() as SVGTextElement).getBBox();
      tipBg
        .attr('x', bbox.x - 4)
        .attr('y', bbox.y - 2)
        .attr('width', bbox.width + 8)
        .attr('height', bbox.height + 4);
      tooltip.style('display', null);
    })
    .on('mouseleave', () => {
      tooltip.style('display', 'none');
    });

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
    .attr('r', 6)
    .attr('fill', (d) => (d.isGhost ? 'none' : 'var(--vermilion)'))
    .attr('stroke', (d) => (d.isGhost ? 'var(--ink-light)' : 'var(--vermilion)'))
    .attr('stroke-width', (d) => (d.isGhost ? 1 : 0))
    .attr('stroke-dasharray', (d) => (d.isGhost ? '2 3' : null))
    .attr('cx', (d) => d.x ?? 0)
    .attr('cy', (d) => d.y ?? 0);

  const dragBeh = drag<SVGCircleElement, SimNode>()
    .on('start', (event) => {
      select(event.sourceEvent.target as SVGCircleElement).raise();
    })
    .on('drag', (event, d) => {
      const [px, py] = pointer(event, svg.node() as Element);
      const [x, y] = currentTransform.invert([px, py]);
      d.x = x;
      d.y = y;
      updateLinks(linkLines, hitLines);
      nodeCircles.attr('cx', (n) => n.x ?? 0).attr('cy', (n) => n.y ?? 0);
      hanLabels.attr('x', (n) => n.x ?? 0).attr('y', (n) => (n.y ?? 0) - 14);
      nameLabels.attr('x', (n) => n.x ?? 0).attr('y', (n) => (n.y ?? 0) + 18);
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
    .attr('font-size', 10)
    .attr('fill', 'var(--vermilion)')
    .attr('font-family', 'Cormorant Garamond, serif')
    .text((d) => d.han ?? '')
    .attr('x', (d) => d.x ?? 0)
    .attr('y', (d) => (d.y ?? 0) - 14);

  const nameLabels = labelsG
    .selectAll<SVGTextElement, SimNode>('text.nm')
    .data(simNodes)
    .join('text')
    .attr('class', 'nm node-label')
    .attr('text-anchor', 'middle')
    .attr('font-size', 10)
    .attr('fill', 'var(--ink-soft)')
    .attr('font-family', 'Be Vietnam Pro, sans-serif')
    .text((d) => d.name)
    .attr('x', (d) => d.x ?? 0)
    .attr('y', (d) => (d.y ?? 0) + 18);

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
