const list = document.getElementById("list");
const form = document.getElementById("new-todo");
const input = document.getElementById("title");
const statusEl = document.getElementById("status");

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

function showError(err) {
  statusEl.textContent = "Error: " + err.message;
  statusEl.className = "error";
}

function render(todos) {
  list.replaceChildren();
  for (const todo of todos) {
    const li = document.createElement("li");
    li.className = todo.done ? "done" : "";

    const cb = document.createElement("input");
    cb.type = "checkbox";
    cb.checked = todo.done;
    cb.addEventListener("change", () =>
      api("PATCH", `/api/todos/${todo.id}`, { done: cb.checked }).then(load).catch(showError)
    );

    const span = document.createElement("span");
    span.textContent = todo.title;

    const del = document.createElement("button");
    del.className = "delete";
    del.textContent = "\u00d7";
    del.title = "Delete";
    del.addEventListener("click", () =>
      api("DELETE", `/api/todos/${todo.id}`).then(load).catch(showError)
    );

    li.append(cb, span, del);
    list.append(li);
  }
  const left = todos.filter((t) => !t.done).length;
  statusEl.className = "";
  statusEl.textContent = todos.length ? `${left} of ${todos.length} remaining` : "Nothing to do yet.";
}

async function load() {
  render(await api("GET", "/api/todos"));
}

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

load().catch(showError);
