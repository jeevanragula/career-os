let state = {};
const $ = (id) => document.getElementById(id);

function esc(value) {
  return String(value ?? "").replace(/[&<>"]/g, (char) => ({
    "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;"
  }[char]));
}

function toast(message, good = false) {
  const el = $("toast");
  if (!el) return;
  el.textContent = message;
  el.className = "toast " + (good ? "good" : "error");
  el.style.display = "block";
  setTimeout(() => { el.style.display = "none"; }, 4000);
}

async function api(url, options = {}) {
  const response = await fetch(url, { ...options, credentials: "same-origin" });
  const text = await response.text();
  let data = {};
  try { data = text ? JSON.parse(text) : {}; } catch (_) {}
  if (!response.ok) throw new Error(data.error || ("Request failed (" + response.status + ")"));
  return data;
}

async function load() {
  try {
    state = await api("/api/dashboard");
    render();
    await Promise.all([loadAgents(), loadCandidates(), loadStatus()]);
  } catch (error) { toast(error.message); }
}

function render() {
  const stats = state.stats || {};
  $("stats").innerHTML = [
    ["opportunities", stats.jobs || 0],
    ["open actions", stats.open_tasks || 0],
    ["recommendations", stats.recommendations || 0],
    ["learning", stats.learning_in_progress || 0]
  ].map(([label, value]) =>
    '<div class="stat"><b>' + esc(value) + '</b><span>' + esc(label) + '</span></div>'
  ).join("");

  $("jobs").innerHTML = (state.opportunities || []).map((job) => {
    const skills = (job.matched_skills || []).map((skill) => '<span>' + esc(skill) + '</span>').join("");
    const gaps = job.gaps && job.gaps.length ? '<div class="meta gap">Gaps: ' + esc(job.gaps.join(", ")) + "</div>" : "";
    const openRole = job.url ? '<a class="btn secondary small" target="_blank" rel="noopener" href="' + esc(job.url) + '">Open role</a>' : "";
    return '<article class="job"><div class="job-main">' +
      '<div class="job-title">' + esc(job.title) + "</div>" +
      '<div class="meta">' + esc(job.company) + " · " + esc(job.location || "Location not specified") + "</div>" +
      "<p>" + esc(job.rationale || "No explanation yet.") + "</p>" +
      (skills ? '<div class="chips">' + skills + "</div>" : "") + gaps +
      '</div><div class="job-actions"><strong>' + Number(job.score || 0).toFixed(0) + "/100</strong>" +
      openRole +
      '<button class="btn small" data-action="analyze" data-id="' + esc(job.id) + '">Analyze</button>' +
      '<button class="btn secondary small" data-action="resume" data-id="' + esc(job.id) + '">Resume</button>' +
      "</div></article>";
  }).join("") || '<div class="empty big">No opportunities yet.<br><span>Run discovery above to build your inbox.</span></div>';

  $("tasks").innerHTML = (state.tasks || []).slice(0, 5).map((task) =>
    '<div class="item"><strong>' + esc(task.title) + '</strong><div class="meta">' + esc(task.area) +
    " · priority " + esc(task.priority) + '</div><button class="btn secondary small" data-action="done-task" data-id="' +
    esc(task.id) + '">Done</button></div>'
  ).join("") || '<div class="empty">No open actions.</div>';

  $("recs").innerHTML = (state.recommendations || []).slice(0, 4).map((rec) =>
    '<div class="item"><strong>' + esc(rec.title) + '</strong><div class="meta">' + esc(rec.agent_name) +
    "</div><p>" + esc(rec.rationale) + '</p><div class="actions">' +
    '<button class="btn small" data-action="recommendation" data-id="' + esc(rec.id) + '" data-status="accepted">Accept</button>' +
    '<button class="btn secondary small" data-action="recommendation" data-id="' + esc(rec.id) + '" data-status="dismissed">Dismiss</button>' +
    "</div></div>"
  ).join("") || '<div class="empty">Run an AI agent when AI is configured.</div>';

  $("goals").innerHTML = (state.goals || []).map((goal) =>
    '<div class="item"><strong>' + esc(goal.title) + '</strong><div class="meta">' + esc(goal.description) + "</div></div>"
  ).join("") || '<div class="empty">No active goals.</div>';
  $("learning").innerHTML = (state.learning || []).map((item) =>
    '<div class="item"><strong>' + esc(item.title) + '</strong><div class="meta">' + esc(item.skill) + " · " + esc(item.status) + "</div></div>"
  ).join("") || '<div class="empty">No learning items.</div>';
  $("network").innerHTML = (state.network || []).map((item) =>
    '<div class="item"><strong>' + esc(item.name) + '</strong><div class="meta">' + esc(item.role) + " · " + esc(item.company) + "</div></div>"
  ).join("") || '<div class="empty">No contacts yet.</div>';
  $("startups").innerHTML = (state.startups || []).map((item) =>
    '<div class="item"><strong>' + esc(item.name) + '</strong><div class="meta">' + esc(item.status) + "</div></div>"
  ).join("") || '<div class="empty">No startup tracks.</div>';
  $("skills").innerHTML = (state.skills || []).map((skill) =>
    '<div class="skill"><strong>' + esc(skill.name) + '</strong><span>' + esc(skill.level) + "/5</span></div>"
  ).join("");
}

async function loadStatus() {
  try {
    const status = await api("/api/status");
    const rows = [
      ["Database", status.database, "Core data store"],
      ["Tavily discovery", status.tavily, "Required for web opportunity search"],
      ["AI agents", status.ai, "Required for Analyze, Resume and agent runs"]
    ];
    $("status").innerHTML = rows.map(([name, ready, description]) =>
      '<div class="status-row"><div><strong>' + esc(name) + '</strong><div class="meta">' + esc(description) +
      '</div></div><span class="status-pill ' + (ready ? "ok" : "warn") + '">' +
      (ready ? "Ready" : "Not configured") + "</span></div>"
    ).join("");
    $("heroTitle").textContent = status.tavily ? "Your opportunity radar is ready" : "Connect Tavily to start discovering opportunities";
    $("heroText").textContent = status.tavily ? "Search for roles and CareerOS will build your opportunity inbox." :
      "Add TAVILY_API_KEY to your .env, restart CareerOS, then run discovery.";
  } catch (error) {
    $("status").innerHTML = '<div class="empty">Unable to load system status: ' + esc(error.message) + "</div>";
  }
}

async function loadAgents() {
  try {
    const data = await api("/api/agents");
    $("agents").innerHTML = (data.agents || []).map((agent) =>
      '<div class="agent"><strong>' + esc(agent.name) + '</strong><div class="meta">' + esc(agent.purpose) +
      '</div><button class="btn small" data-action="run-agent" data-name="' + esc(agent.name) + '">Run</button></div>'
    ).join("");
  } catch (error) { $("agents").innerHTML = '<div class="empty">' + esc(error.message) + "</div>"; }
}

async function loadCandidates() {
  try {
    const data = await api("/api/discovery/candidates");
    $("candidates").innerHTML = (data.candidates || []).map((candidate) =>
      '<div class="item"><strong>' + esc(candidate.name) + '</strong><div class="meta">' + esc(candidate.evidence) + "</div></div>"
    ).join("") || '<div class="empty">No discovered companies yet.</div>';
  } catch (error) { $("candidates").innerHTML = '<div class="empty">' + esc(error.message) + "</div>"; }
}

async function runDiscovery() {
  const body = {
    roles: $("roles").value.split(",").map((value) => value.trim()).filter(Boolean),
    locations: $("locations").value.split(",").map((value) => value.trim()).filter(Boolean)
  };
  const buttons = [$("discover"), $("discoverHero"), $("discoverBottom")].filter(Boolean);
  buttons.forEach((button) => { button.disabled = true; });
  $("discoveryStatus").textContent = "Starting discovery…";
  try {
    const data = await api("/api/discover/automatic", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body)
    });
    const runID = data.run_id;
    $("discoveryStatus").textContent = "Discovery running…";
    let attempts = 0;
    while (attempts++ < 600) {
      await new Promise((resolve) => setTimeout(resolve, 2000));
      const run = await api("/api/discovery/runs/" + encodeURIComponent(runID));
      if (run.status === "failed") throw new Error(run.error || "Discovery failed");
      if (run.status === "completed") {
        const result = run.query || {};
        $("discoveryStatus").textContent =
          "Search completed: " + (run.discovered_count || 0) +
          " Tavily results. Review them below before crawling.";
        toast("Tavily search completed", true);
        await load();
        return;
      }
      $("discoveryStatus").textContent = "Discovery running… " + (run.discovered_count || 0) + " companies/signals found.";
    }
    throw new Error("Discovery is still running. Refresh later to see the result.");
  } catch (error) {
    $("discoveryStatus").textContent = "Discovery failed: " + error.message;
    toast(error.message);
  } finally {
    buttons.forEach((button) => { button.disabled = false; });
  }
}

async function runAgent(name) {
  toast("Running " + name + "…");
  try {
    await api("/api/agents/run", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ name }) });
    toast(name + " completed", true);
    await load();
  } catch (error) { toast(error.message); }
}

async function doneTask(id) {
  try { await api("/api/tasks/" + encodeURIComponent(id), { method: "POST" }); await load(); }
  catch (error) { toast(error.message); }
}

async function updateRecommendation(id, status) {
  try {
    await api("/api/recommendations/" + encodeURIComponent(id) + "?status=" + encodeURIComponent(status), { method: "POST" });
    await load();
  } catch (error) { toast(error.message); }
}

async function analyze(id) {
  try {
    await api("/api/analyze", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ job_id: id }) });
    toast("Job analyzed", true);
    await load();
  } catch (error) { toast(error.message); }
}

async function resume(id) {
  try {
    const data = await api("/api/resume", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ job_id: id }) });
    const popup = window.open("", "_blank");
    if (!popup) throw new Error("Popup blocked. Allow popups for localhost:8080 and try again.");
    popup.document.write('<pre style="white-space:pre-wrap;font:14px system-ui;padding:24px">' + esc(data.content) + "</pre>");
    popup.document.close();
  } catch (error) { toast(error.message); }
}

async function handleAction(event) {
  const target = event.target.closest("[data-action]");
  if (!target) return;
  const action = target.dataset.action;
  if (action === "analyze") await analyze(target.dataset.id);
  else if (action === "resume") await resume(target.dataset.id);
  else if (action === "done-task") await doneTask(target.dataset.id);
  else if (action === "recommendation") await updateRecommendation(target.dataset.id, target.dataset.status);
  else if (action === "run-agent") await runAgent(target.dataset.name);
}

function bindEvents() {
  $("refresh").addEventListener("click", load);
  $("discover").addEventListener("click", runDiscovery);
  $("discoverHero").addEventListener("click", runDiscovery);
  $("discoverBottom").addEventListener("click", runDiscovery);
  document.addEventListener("click", handleAction);
}

document.addEventListener("DOMContentLoaded", () => {
  bindEvents();
  load();
});
