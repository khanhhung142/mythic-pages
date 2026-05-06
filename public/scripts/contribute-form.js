const MAX_TEXT_LENGTH = 30000;

function slugify(input) {
  return (input || '')
    .toString()
    .normalize('NFKD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '');
}

function getTurnstileToken() {
  if (!('turnstile' in window)) return null;
  try {
    return window.turnstile.getResponse();
  } catch {
    return null;
  }
}

function setBusy(isBusy) {
  const form = document.getElementById('contribute-form');
  const button = form?.querySelector('button[type="submit"]');
  if (!button) return;
  button.toggleAttribute('disabled', isBusy);
  button.classList.toggle('is-busy', isBusy);
}

function setStatus(message, kind) {
  const el = document.getElementById('form-status');
  if (!el) return;
  el.textContent = message || '';
  el.dataset.kind = kind || '';
}

function setMode(mode) {
  const modeInput = document.getElementById('mode');
  const tabs = document.querySelectorAll('.mode-tab');
  const newPanel = document.getElementById('mode-new');
  const editPanel = document.getElementById('mode-edit');
  if (!modeInput || !newPanel || !editPanel) return;

  modeInput.value = mode;
  const isNew = mode === 'new';
  newPanel.classList.toggle('is-hidden', !isNew);
  newPanel.toggleAttribute('hidden', !isNew);
  editPanel.classList.toggle('is-hidden', isNew);
  editPanel.toggleAttribute('hidden', isNew);

  tabs.forEach((tab) => {
    const isActive = tab.getAttribute('data-mode') === mode;
    tab.classList.toggle('is-active', isActive);
    tab.setAttribute('aria-selected', isActive ? 'true' : 'false');
  });

  setStatus('', '');
}

function toJson(form) {
  const data = new FormData(form);
  return Object.fromEntries(data.entries());
}

function getMessages() {
  const actions = document.querySelector('.actions');
  const dataset = actions?.dataset ?? {};
  return (key, fallback) => dataset[key] || fallback;
}

async function submitContribute(payload) {
  const res = await fetch('/api/contribute', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  const json = await res.json().catch(() => null);
  return { res, json };
}

function renderIssueUrl(issueUrl, msgCreatedIssue) {
  const a = document.createElement('a');
  a.href = issueUrl;
  a.target = '_blank';
  a.rel = 'noreferrer';
  a.textContent = issueUrl;
  const el = document.getElementById('form-status');
  if (!el) return;
  el.dataset.kind = 'ok';
  el.textContent = '';
  el.appendChild(document.createTextNode(msgCreatedIssue));
  el.appendChild(a);
}

function initContributeForm() {
  const form = document.getElementById('contribute-form');
  if (!form) return;

  const msg = getMessages();

  const titleInput = document.getElementById('new-title');
  const slugInput = document.getElementById('new-slug');

  titleInput?.addEventListener('input', () => {
    if (!slugInput) return;
    if (slugInput.value && slugInput.dataset.touched === '1') return;
    slugInput.value = slugify(titleInput.value);
  });
  slugInput?.addEventListener('input', () => {
    if (!slugInput) return;
    slugInput.dataset.touched = '1';
  });

  document.querySelectorAll('.mode-tab').forEach((tab) => {
    tab.addEventListener('click', () => setMode(tab.getAttribute('data-mode')));
  });

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    setStatus('', '');

    const payload = toJson(form);
    const mode = payload.mode === 'edit' ? 'edit' : 'new';
    payload.mode = mode;

    if (payload.website) {
      setStatus('OK.', 'ok');
      return;
    }

    const turnstileToken = getTurnstileToken();
    if (!turnstileToken) {
      setStatus(msg('msgTurnstileMissing', 'Turnstile missing.'), 'err');
      return;
    }
    payload.turnstileToken = turnstileToken;

    const required = [];
    if (mode === 'new') {
      required.push(['title', payload.title], ['markdownBody', payload.markdownBody]);
    } else {
      required.push(
        ['entryId', payload.entryId],
        ['changeSummary', payload.changeSummary],
        ['proposedMarkdown', payload.proposedMarkdown]
      );
    }
    const missing = required.filter(([, v]) => !String(v || '').trim()).map(([k]) => k);
    if (missing.length) {
      setStatus(msg('msgMissingRequired', 'Missing required fields.'), 'err');
      return;
    }

    const textFields = ['markdownBody', 'proposedMarkdown', 'changeSummary', 'sources', 'notes'];
    for (const key of textFields) {
      if (payload[key] && String(payload[key]).length > MAX_TEXT_LENGTH) {
        setStatus(msg('msgContentTooLong', 'Content too long.'), 'err');
        return;
      }
    }

    setBusy(true);
    setStatus(msg('msgSubmitting', 'Submitting…'), 'busy');
    try {
      const { res, json } = await submitContribute(payload);
      if (!res.ok) {
        setStatus(
          json && json.error ? String(json.error) : msg('msgSubmitFailed', 'Submit failed.'),
          'err'
        );
        return;
      }
      if (json && json.issueUrl) {
        renderIssueUrl(String(json.issueUrl), msg('msgCreatedIssue', 'Created issue: '));
        return;
      }
      setStatus(msg('msgSubmitted', 'Submitted.'), 'ok');
    } catch {
      setStatus(msg('msgNetworkError', 'Network error.'), 'err');
    } finally {
      setBusy(false);
    }
  });
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initContributeForm, { once: true });
} else {
  initContributeForm();
}

