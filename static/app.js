const $ = (id) => document.getElementById(id);
const list = $("list"), form = $("new-todo"), input = $("title");
let todos = [];
let filter = "all";
let toastTimer;

async function api(method, path, body) {
  const res = await fetch(path, {
    method,
    headers: body ? { "Content-Type": "application/json" } : {},
    body: body ? JSON.stringify(body) : undefined,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || res.statusText);
  }
  return res.status === 204 ? null : res.json();
}

function toast(msg) {
  const t = $("toast");
  t.textContent = msg;
  t.classList.add("show");
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => t.classList.remove("show"), 3500);
}
const showError = (err) => toast("⚠️ " + err.message);

function ago(iso) {
  const s = Math.max(0, (Date.now() - new Date(iso)) / 1000);
  if (s < 60) return "just now";
  const units = [["y", 31536000], ["mo", 2592000], ["d", 86400], ["h", 3600], ["m", 60]];
  for (const [u, n] of units) if (s >= n) return `${Math.floor(s / n)}${u} ago`;
}

function startEdit(li, span, todo) {
  const edit = document.createElement("input");
  edit.className = "edit";
  edit.value = todo.title;
  edit.maxLength = 500;
  span.replaceWith(edit);
  edit.focus();
  edit.select();
  let finished = false;
  const finish = (save) => {
    if (finished) return;
    finished = true;
    const title = edit.value.trim();
    if (save && title && title !== todo.title) {
      api("PATCH", `/api/todos/${todo.id}`, { title }).then(load).catch(showError);
    } else render();
  };
  edit.addEventListener("blur", () => finish(true));
  edit.addEventListener("keydown", (e) => {
    if (e.key === "Enter") finish(true);
    if (e.key === "Escape") finish(false);
  });
}

function item(todo) {
  const li = document.createElement("li");
  li.className = "item" + (todo.done ? " done" : "");

  const check = document.createElement("button");
  check.className = "check";
  check.setAttribute("aria-label", todo.done ? "Mark as not done" : "Mark as done");
  check.innerHTML = '<svg viewBox="0 0 24 24"><path d="M5 12.5l4.5 4.5L19 7.5"/></svg>';
  check.addEventListener("click", () =>
    api("PATCH", `/api/todos/${todo.id}`, { done: !todo.done }).then(load).catch(showError));

  const body = document.createElement("div");
  body.className = "body";
  const span = document.createElement("span");
  span.className = "text";
  span.textContent = todo.title;
  span.addEventListener("dblclick", () => startEdit(li, span, todo));
  const meta = document.createElement("small");
  meta.textContent = ago(todo.created_at);
  meta.title = new Date(todo.created_at).toLocaleString();
  body.append(span, meta);

  const del = document.createElement("button");
  del.className = "delete";
  del.setAttribute("aria-label", "Delete");
  del.innerHTML = '<svg viewBox="0 0 24 24"><path d="M6 6l12 12M18 6L6 18"/></svg>';
  del.addEventListener("click", () => {
    li.classList.add("leaving");
    setTimeout(() => api("DELETE", `/api/todos/${todo.id}`).then(load).catch(showError), 200);
  });

  li.append(check, body, del);
  return li;
}

function render() {
  const done = todos.filter((t) => t.done).length;
  const active = todos.length - done;
  $("c-all").textContent = todos.length;
  $("c-active").textContent = active;
  $("c-done").textContent = done;

  const pct = todos.length ? Math.round((done / todos.length) * 100) : 0;
  $("ring-fg").style.strokeDasharray = `${pct} 100`;
  $("ring-label").textContent = pct + "%";

  // Newest first.
  const shown = todos
    .filter((t) => filter === "all" || (filter === "done" ? t.done : !t.done))
    .slice()
    .reverse();
  list.replaceChildren(...shown.map(item));

  $("empty").hidden = shown.length > 0;
  $("empty-sub").textContent = todos.length === 0
    ? "Add your first task above to get started."
    : filter === "done" ? "Nothing completed yet — you've got this." : "Everything is done. Enjoy the moment!";
  $("status").textContent = todos.length
    ? `${active} task${active === 1 ? "" : "s"} left`
    : "No tasks";
}

async function load() {
  todos = await api("GET", "/api/todos");
  render();
}

$("filters").addEventListener("click", (e) => {
  const btn = e.target.closest("button[data-filter]");
  if (!btn) return;
  filter = btn.dataset.filter;
  for (const b of $("filters").children) b.classList.toggle("active", b === btn);
  render();
});

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  const title = input.value.trim();
  if (!title) return;
  try {
    await api("POST", "/api/todos", { title });
    input.value = "";
    await load();
  } catch (err) {
    showError(err);
  }
});

$("today").textContent = new Date().toLocaleDateString(undefined, {
  weekday: "long", month: "long", day: "numeric",
});
load().catch(showError);
