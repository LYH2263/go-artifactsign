async function refresh() {
  const s = await fetch('/api/stats').then(r => r.json());
  document.getElementById('stats-line').textContent =
    `signs=${s.Signs} verifies=${s.Verifies} digests=${s.Digests} revokes=${s.Revokes} trusts=${s.Trusts}`;
}
document.getElementById('btn-refresh').onclick = refresh;
document.getElementById('key-form').onsubmit = async (e) => {
  e.preventDefault();
  const id = document.getElementById('key-id').value;
  const text = document.getElementById('key-mat').value;
  const material = Array.from(new TextEncoder().encode(text));
  const res = await fetch('/api/keys', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id, material }),
  });
  document.getElementById('key-out').textContent = await res.text();
  refresh();
};
let lastSig = null;
document.getElementById('sign-form').onsubmit = async (e) => {
  e.preventDefault();
  const name = document.getElementById('sign-name').value;
  const text = document.getElementById('sign-payload').value;
  const payload = Array.from(new TextEncoder().encode(text));
  const res = await fetch('/api/sign', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ payload, name }),
  });
  const body = await res.text();
  document.getElementById('sign-out').textContent = body;
  try { lastSig = JSON.parse(body); } catch (_) { lastSig = null; }
  refresh();
};
document.getElementById('verify-form').onsubmit = async (e) => {
  e.preventDefault();
  const key_id = document.getElementById('ver-key').value;
  const text = document.getElementById('ver-payload').value;
  const payload = Array.from(new TextEncoder().encode(text));
  const sig = lastSig && lastSig.Signature ? lastSig.Signature : [];
  const res = await fetch('/api/verify', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ key_id, payload, sig }),
  });
  document.getElementById('verify-out').textContent = await res.text();
};
refresh();
