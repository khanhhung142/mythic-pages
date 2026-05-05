import type { PagesFunction } from '@cloudflare/workers-types';


type ContributeMode = 'new' | 'edit';

type ContributePayload = {
  mode: ContributeMode;
  turnstileToken: string;
  displayName?: string;
  contactEmail?: string;
  sources?: string;
  notes?: string;
  website?: string; // honeypot

  // new
  title?: string;
  suggestedSlug?: string;
  category?: string;
  summary?: string;
  markdownBody?: string;

  // edit
  entryId?: string;
  changeSummary?: string;
  proposedMarkdown?: string;
};

type Env = {
  TURNSTILE_SECRET_KEY: string;
  GITHUB_TOKEN: string;
  GITHUB_REPO: string; // owner/repo
};

const MAX_TEXT_LENGTH = 30000;

function jsonResponse(body: unknown, init?: ResponseInit) {
  return new Response(JSON.stringify(body), {
    headers: {
      'content-type': 'application/json; charset=utf-8',
      ...(init?.headers ?? {}),
    },
    ...init,
  });
}

function isNonEmptyString(value: unknown): value is string {
  return typeof value === 'string' && value.trim().length > 0;
}

function getStringOrEmpty(value: unknown): string {
  return typeof value === 'string' ? value.trim() : '';
}

function validatePayload(payload: ContributePayload) {
  const mode: ContributeMode = payload.mode === 'edit' ? 'edit' : 'new';
  const errors: string[] = [];

  if (!isNonEmptyString(payload.turnstileToken)) errors.push('turnstileToken');

  if (mode === 'new') {
    if (!isNonEmptyString(payload.title)) errors.push('title');
    if (!isNonEmptyString(payload.markdownBody)) errors.push('markdownBody');
  } else {
    if (!isNonEmptyString(payload.entryId)) errors.push('entryId');
    if (!isNonEmptyString(payload.changeSummary)) errors.push('changeSummary');
    if (!isNonEmptyString(payload.proposedMarkdown)) errors.push('proposedMarkdown');
  }

  const textFields: Array<keyof ContributePayload> = [
    'markdownBody',
    'proposedMarkdown',
    'changeSummary',
    'sources',
    'notes',
    'summary',
  ];
  for (const key of textFields) {
    const v = payload[key];
    if (typeof v === 'string' && v.length > MAX_TEXT_LENGTH) errors.push(String(key));
  }

  return { mode, errors };
}

async function verifyTurnstile(args: { secret: string; token: string; ip?: string }) {
  const form = new FormData();
  form.set('secret', args.secret);
  form.set('response', args.token);
  if (args.ip) form.set('remoteip', args.ip);

  const res = await fetch('https://challenges.cloudflare.com/turnstile/v0/siteverify', {
    method: 'POST',
    body: form,
  });

  const json = (await res.json().catch(() => null)) as
    | null
    | { success?: boolean; 'error-codes'?: string[] };

  return Boolean(json?.success);
}

function buildIssueTitle(payload: ContributePayload, mode: ContributeMode): string {
  if (mode === 'edit') {
    return `[Contribute] Edit: ${getStringOrEmpty(payload.entryId) || 'unknown'}`;
  }
  return `[Contribute] New entry: ${getStringOrEmpty(payload.title) || 'untitled'}`;
}

function buildIssueBody(payload: ContributePayload, mode: ContributeMode): string {
  const displayName = getStringOrEmpty(payload.displayName);
  const contactEmail = getStringOrEmpty(payload.contactEmail);
  const sources = getStringOrEmpty(payload.sources);
  const notes = getStringOrEmpty(payload.notes);

  const lines: string[] = [];
  lines.push('## Submission');
  lines.push(`- Mode: \`${mode}\``);
  if (displayName) lines.push(`- Display name: ${displayName}`);
  if (contactEmail) lines.push(`- Contact email: ${contactEmail}`);
  lines.push('');

  if (mode === 'new') {
    lines.push('## New entry');
    lines.push(`- Title: ${getStringOrEmpty(payload.title)}`);
    const suggestedSlug = getStringOrEmpty(payload.suggestedSlug);
    if (suggestedSlug) lines.push(`- Suggested slug: \`${suggestedSlug}\``);
    const category = getStringOrEmpty(payload.category);
    if (category) lines.push(`- Category: ${category}`);
    const summary = getStringOrEmpty(payload.summary);
    if (summary) {
      lines.push('');
      lines.push('### Summary');
      lines.push(summary);
    }
    lines.push('');
    lines.push('### Proposed markdown');
    lines.push('```md');
    lines.push(getStringOrEmpty(payload.markdownBody));
    lines.push('```');
    lines.push('');
  } else {
    lines.push('## Edit proposal');
    lines.push(`- Entry ID: \`${getStringOrEmpty(payload.entryId)}\``);
    lines.push('');
    lines.push('### Change summary');
    lines.push(getStringOrEmpty(payload.changeSummary));
    lines.push('');
    lines.push('### Proposed markdown');
    lines.push('```md');
    lines.push(getStringOrEmpty(payload.proposedMarkdown));
    lines.push('```');
    lines.push('');
  }

  if (sources) {
    lines.push('## Sources / links');
    lines.push(sources);
    lines.push('');
  }

  if (notes) {
    lines.push('## Notes for maintainers');
    lines.push(notes);
    lines.push('');
  }

  lines.push('---');
  lines.push('_Created via contribute form (Cloudflare Pages Functions + Turnstile)._');

  return lines.join('\n');
}

async function createGitHubIssue(args: {
  token: string;
  repo: string;
  title: string;
  body: string;
  labels?: string[];
}) {
  const res = await fetch(`https://api.github.com/repos/${args.repo}/issues`, {
    method: 'POST',
    headers: {
      authorization: `Bearer ${args.token}`,
      accept: 'application/vnd.github+json',
      'content-type': 'application/json',
      'user-agent': 'vietmyth-contribute-bot',
    },
    body: JSON.stringify({
      title: args.title,
      body: args.body,
      ...(args.labels ? { labels: args.labels } : {}),
    }),
  });

  const json = (await res.json().catch(() => null)) as any;
  return { ok: res.ok, status: res.status, json };
}

export const onRequestPost: PagesFunction<Env> = async (context) => {
  const { request, env } = context;

  let payload: ContributePayload;
  try {
    payload = (await request.json()) as ContributePayload;
  } catch {
    return jsonResponse({ error: 'invalid_json' }, { status: 400 });
  }

  if (payload.website) {
    return jsonResponse({ ok: true }, { status: 200 });
  }

  const { mode, errors } = validatePayload(payload);
  if (errors.length) {
    return jsonResponse({ error: 'invalid_request', details: errors }, { status: 400 });
  }

  const ip = request.headers.get('CF-Connecting-IP') ?? undefined;
  const isHuman = await verifyTurnstile({
    secret: env.TURNSTILE_SECRET_KEY,
    token: payload.turnstileToken,
    ip,
  });
  if (!isHuman) {
    return jsonResponse({ error: 'turnstile_failed' }, { status: 400 });
  }

  const repo = env.GITHUB_REPO;
  if (!repo || !repo.includes('/')) {
    return jsonResponse({ error: 'github_repo_missing' }, { status: 500 });
  }

  const title = buildIssueTitle(payload, mode);
  const body = buildIssueBody(payload, mode);
  const labels = ['contribution', mode === 'new' ? 'new-entry' : 'edit-entry'];

  const first = await createGitHubIssue({
    token: env.GITHUB_TOKEN,
    repo,
    title,
    body,
    labels,
  });

  let issue = first;
  if (!issue.ok && issue.status === 422) {
    issue = await createGitHubIssue({
      token: env.GITHUB_TOKEN,
      repo,
      title,
      body,
    });
  }

  if (!issue.ok) {
    return jsonResponse({ error: 'github_failed', status: issue.status }, { status: 502 });
  }

  const issueUrl = issue.json?.html_url;
  if (!issueUrl) {
    return jsonResponse({ error: 'github_failed' }, { status: 502 });
  }

  return jsonResponse({ issueUrl }, { status: 200 });
};

