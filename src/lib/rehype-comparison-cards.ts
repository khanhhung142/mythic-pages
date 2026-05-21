import type { Element, ElementContent, Root } from 'hast';

type HastParent = Root | Element;

/** Replace GFM tables with `.comp-cards` (CSS-only tabs) at build time — no client JS. */
export function rehypeComparisonCards() {
  let tableIndex = 0;

  return (tree: Root) => {
    tableIndex = 0;
    walk(tree, () => tableIndex++);
  };
}

function walk(node: HastParent, nextTableId: () => number) {
  const children = node.children;
  if (!children) return;

  for (let i = 0; i < children.length; i++) {
    const child = children[i];
    if (!child || child.type !== 'element') continue;

    if (child.tagName === 'table') {
      const cards = tableToCompCards(child, nextTableId());
      if (cards) children[i] = cards;
      continue;
    }

    if (child.tagName === 'div' && hasClass(child, 'table-scroll')) {
      const table = child.children?.find(
        (c): c is Element => c.type === 'element' && c.tagName === 'table',
      );
      if (table) {
        const cards = tableToCompCards(table, nextTableId());
        if (cards) {
          children[i] = cards;
          continue;
        }
      }
    }

    walk(child, nextTableId);
  }
}

function hasClass(node: Element, className: string): boolean {
  const cls = node.properties?.class;
  if (typeof cls === 'string') return cls.split(/\s+/).includes(className);
  if (Array.isArray(cls)) return cls.includes(className);
  return false;
}

function tableToCompCards(table: Element, tableId: number): Element | null {
  const headerCells = getHeaderCells(table);
  if (headerCells.length < 2) return null;

  const sourceHeaders = headerCells.slice(1);
  const bodyRows = getBodyRows(table);
  if (!bodyRows.length) return null;

  const groupName = `comp-${tableId}`;
  const children: ElementContent[] = [];

  if (sourceHeaders.length > 1) {
    const tabLabels: ElementContent[] = [];
    sourceHeaders.forEach((headerCell, sourceIndex) => {
      const radioId = `${groupName}-${sourceIndex}`;
      const radioProps: Element['properties'] = {
        type: 'radio',
        name: groupName,
        id: radioId,
        class: 'comp-radio',
        'data-index': String(sourceIndex),
      };
      if (sourceIndex === 0) radioProps.checked = true;

      children.push({
        type: 'element',
        tagName: 'input',
        properties: radioProps,
        children: [],
      });

      tabLabels.push({
        type: 'element',
        tagName: 'label',
        properties: {
          for: radioId,
          class: 'comp-tab',
          'data-index': String(sourceIndex),
        },
        children: [{ type: 'text', value: getShortTabLabel(hastToText(headerCell)) }],
      });
    });

    children.push({
      type: 'element',
      tagName: 'div',
      properties: { class: 'comp-tabs', role: 'tablist' },
      children: tabLabels,
    });
  }

  sourceHeaders.forEach((headerCell, sourceIndex) => {
    const columnIndex = sourceIndex + 1;
    const panelChildren: ElementContent[] = [
      {
        type: 'element',
        tagName: 'div',
        properties: { class: 'comp-source-name' },
        children: cloneChildren(headerCell.children),
      },
      buildRowsList(bodyRows, columnIndex),
    ];

    children.push({
      type: 'element',
      tagName: 'div',
      properties: {
        class: 'comp-panel',
        'data-panel': String(sourceIndex),
      },
      children: panelChildren,
    });
  });

  return {
    type: 'element',
    tagName: 'div',
    properties: { class: 'comp-cards' },
    children,
  };
}

function getHeaderCells(table: Element): Element[] {
  const thead = table.children?.find(
    (c): c is Element => c.type === 'element' && c.tagName === 'thead',
  );
  if (thead) {
    const row = thead.children?.find(
      (c): c is Element => c.type === 'element' && c.tagName === 'tr',
    );
    if (row) return getRowCells(row);
  }

  const firstRow = table.children?.find(
    (c): c is Element => c.type === 'element' && c.tagName === 'tr',
  );
  return firstRow ? getRowCells(firstRow) : [];
}

function getBodyRows(table: Element): Element[] {
  const tbody = table.children?.find(
    (c): c is Element => c.type === 'element' && c.tagName === 'tbody',
  );
  if (tbody) {
    return (tbody.children ?? []).filter(
      (c): c is Element => c.type === 'element' && c.tagName === 'tr',
    );
  }

  const rows = (table.children ?? []).filter(
    (c): c is Element => c.type === 'element' && c.tagName === 'tr',
  );
  return rows.length > 1 ? rows.slice(1) : [];
}

function getRowCells(row: Element): Element[] {
  return (row.children ?? []).filter(
    (c): c is Element =>
      c.type === 'element' && (c.tagName === 'th' || c.tagName === 'td'),
  );
}

function buildRowsList(bodyRows: Element[], columnIndex: number): Element {
  const rows: ElementContent[] = [];

  for (const row of bodyRows) {
    const cells = getRowCells(row);
    if (cells.length <= columnIndex) continue;

    const labelCell = cells[0];
    const valueCell = cells[columnIndex];
    if (!labelCell || !valueCell) continue;

    rows.push({
      type: 'element',
      tagName: 'div',
      properties: { class: 'comp-row' },
      children: [
        {
          type: 'element',
          tagName: 'dt',
          properties: {},
          children: cloneChildren(labelCell.children),
        },
        {
          type: 'element',
          tagName: 'dd',
          properties: {},
          children: cloneChildren(valueCell.children),
        },
      ],
    });
  }

  return {
    type: 'element',
    tagName: 'dl',
    properties: { class: 'comp-rows' },
    children: rows,
  };
}

function cloneChildren(children: ElementContent[] | undefined): ElementContent[] {
  return children ? [...children] : [];
}

function hastToText(node: ElementContent): string {
  if (node.type === 'text') return node.value ?? '';
  if (node.type === 'element' && node.children) {
    return node.children.map(hastToText).join('');
  }
  return '';
}

function getShortTabLabel(text: string): string {
  const trimmed = text.replace(/\s+/g, ' ').trim();
  const beforeParen = trimmed.split('(')[0].trim();
  if (beforeParen.length <= 24) return beforeParen;
  return `${beforeParen.slice(0, 22)}…`;
}
