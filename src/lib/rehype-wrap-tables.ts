import type { Root } from 'hast';

/** Wrap each `<table>` in `<div class="table-scroll">` for horizontal scroll without breaking table layout. */
export function rehypeWrapTables() {
  return (tree: Root) => {
    walk(tree);
  };
}

function walk(node: Root | { type?: string; tagName?: string; children?: unknown[] }) {
  const children = node.children as
    | Array<{ type?: string; tagName?: string; children?: unknown[] }>
    | undefined;
  if (!children) return;

  for (let i = 0; i < children.length; i++) {
    const child = children[i];
    if (child?.type === 'element' && child.tagName === 'table') {
      children[i] = {
        type: 'element',
        tagName: 'div',
        properties: { class: 'table-scroll' },
        children: [child],
      };
    } else {
      walk(child);
    }
  }
}
