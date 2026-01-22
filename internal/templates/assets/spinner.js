const ROUND_PLAN = [
  { id: "door-1", label: "Door Prize 1", count: 25, type: "door" },
  { id: "door-2", label: "Door Prize 2", count: 25, type: "door" },
  { id: "door-3", label: "Door Prize 3", count: 22, type: "door" },
  { id: "door-4", label: "Door Prize 4", count: 26, type: "door" },
  { id: "grand-1", label: "Grand Prize 1", count: 1, type: "grand" },
  { id: "grand-2", label: "Grand Prize 2", count: 1, type: "grand" },
  { id: "grand-3", label: "Grand Prize 3", count: 1, type: "grand" },
];

const TOTAL_WINNERS = ROUND_PLAN.reduce((sum, round) => sum + round.count, 0);
const DOORPRIZE_TOTAL = ROUND_PLAN.filter(
  (round) => round.type === "door",
).reduce((sum, round) => sum + round.count, 0);
const GRAND_PRIZE_TOTAL = ROUND_PLAN.filter(
  (round) => round.type === "grand",
).reduce((sum, round) => sum + round.count, 0);
const TAD_TARGET = 15;
const CHUNK_SIZE = 5;
const CHUNK_DURATION = 1400;
const CHUNK_PAUSE = 1000;
const GRAND_CHUNK_PAUSE = CHUNK_PAUSE * 5;
const SPIN_TICK_MS = 100;
let confettiTimer = null;
let spinRows = [];
const API_KEY = "Dptaspen@25!";

const state = {
  allEmployees: [],
  employees: [],
  pool: {
    regular: {
      organik: [],
      tad: [],
    },
    guaranteed: {
      organik: [],
      tad: [],
    },
  },
  remainingTad: TAD_TARGET,
  guaranteedRemaining: 0,
  doorPrizeDraws: 0,
  grandPrizeDraws: 0,
  results: [],
  currentRound: 0,
  spinning: false,
};

const elements = {
  dataStatus: document.getElementById("dataStatus"),
  alertBox: document.getElementById("alertBox"),
  totalEmployees: document.getElementById("totalEmployees"),
  totalOrganik: document.getElementById("totalOrganik"),
  totalTad: document.getElementById("totalTad"),
  totalWinners: document.getElementById("totalWinners"),
  quotaText: document.getElementById("quotaText"),
  quotaFill: document.getElementById("quotaFill"),
  remainingText: document.getElementById("remainingText"),
  currentRoundLabel: document.getElementById("currentRoundLabel"),
  currentRoundSub: document.getElementById("currentRoundSub"),
  spinnerName: document.getElementById("spinnerName"),
  spinnerSub: document.getElementById("spinnerSub"),
  spinBtn: document.getElementById("spinBtn"),
  resetBtn: document.getElementById("resetBtn"),
  roundHint: document.getElementById("roundHint"),
  roundList: document.getElementById("roundList"),
  results: document.getElementById("results"),
  resultsSummary: document.getElementById("resultsSummary"),
  existingWinnersStatus: document.getElementById("existingWinnersStatus"),
  reloadBtn: document.getElementById("reloadBtn"),
  confetti: document.querySelector(".confetti"),
  spinOverlay: document.getElementById("spinOverlay"),
  overlayTitle: document.getElementById("overlayTitle"),
  overlaySub: document.getElementById("overlaySub"),
  spinFrame: document.getElementById("spinFrame"),
  spinReel: document.getElementById("spinReel"),
  nextChunkBtn: document.getElementById("nextChunkBtn"),
  overlayHint: document.getElementById("overlayHint"),
};

function fetchWithKey(url, options = {}) {
  const headers = new Headers(options.headers || {});
  headers.set("X-API-Key", API_KEY);
  return fetch(url, { ...options, headers });
}

function setAlert(message) {
  if (!message) {
    elements.alertBox.classList.remove("is-visible");
    elements.alertBox.textContent = "";
    return;
  }
  elements.alertBox.textContent = message;
  elements.alertBox.classList.add("is-visible");
}

function setStatus(message) {
  elements.dataStatus.textContent = message;
}

function setExistingWinnersStatus(message) {
  if (!elements.existingWinnersStatus) return;
  elements.existingWinnersStatus.textContent = message || "";
}

function normalizeEmployee(raw) {
  if (!raw) return null;
  const typeRaw = String(
    raw.JENIS_KEPEGAWAIAN ?? raw.jenis_kepegawaian ?? raw.type ?? "",
  ).toUpperCase();
  if (typeRaw !== "TAD" && typeRaw !== "ORGANIK") {
    return null;
  }
  const isExcluded = parseBoolish(
    raw.IS_EXCLUDED ?? raw.is_excluded ?? raw.isExcluded ?? raw.excluded,
  );
  const guaranteedDoorprize = parseBoolish(
    raw.GUARANTEED_DOORPRIZE ??
      raw.guaranteed_doorprize ??
      raw.guaranteedDoorprize,
  );
  const fallbackId = `${Date.now()}-${Math.floor(Math.random() * 1e9)}`;
  return {
    id: raw.ID ?? raw.id ?? raw.NIP ?? raw.nip ?? fallbackId,
    name: raw.NAMA_KARYAWAN ?? raw.nama_karyawan ?? raw.name ?? "Tanpa Nama",
    position: raw.JABATAN ?? raw.jabatan ?? raw.position ?? "-",
    branch: raw.KANTOR_CABANG ?? raw.kantor_cabang ?? raw.branch ?? "-",
    type: typeRaw,
    isExcluded,
    guaranteedDoorprize,
  };
}

function parseBoolish(value) {
  if (typeof value === "boolean") return value;
  if (typeof value === "number") return value === 1;
  if (typeof value === "string") {
    const normalized = value.trim().toLowerCase();
    return normalized === "1" || normalized === "true" || normalized === "yes";
  }
  return false;
}

function parseDate(value) {
  if (!value) return null;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return null;
  }
  return date;
}

function normalizeExistingWinner(raw) {
  if (!raw) return null;
  const prizeType = String(raw.prize_type ?? raw.prizeType ?? "")
    .trim()
    .toLowerCase();
  if (prizeType !== "door" && prizeType !== "grand") {
    return null;
  }

  const typeRaw = String(
    raw.employment_type ?? raw.employmentType ?? raw.type ?? "",
  )
    .trim()
    .toUpperCase();
  const employmentType = typeRaw === "TAD" ? "TAD" : "ORGANIK";

  const employeeIDValue = raw.employee_id ?? raw.employeeId ?? raw.id ?? "";
  const employeeID = employeeIDValue ? String(employeeIDValue).trim() : "";
  const name = String(raw.name ?? raw.NAMA_KARYAWAN ?? "-").trim() || "-";
  const position = String(raw.position ?? raw.JABATAN ?? "-").trim() || "-";
  const branch = String(raw.branch ?? raw.KANTOR_CABANG ?? "-").trim() || "-";
  const roundID = String(raw.round_id ?? raw.roundId ?? "").trim();
  const roundLabel = String(raw.round_label ?? raw.roundLabel ?? "").trim();

  return {
    id: employeeID,
    name,
    position,
    branch,
    type: employmentType,
    prizeType,
    roundId: roundID,
    roundLabel,
    wonAt: parseDate(raw.won_at ?? raw.wonAt),
  };
}

function getNextRoundIndex(winnersByRound) {
  for (let i = 0; i < ROUND_PLAN.length; i += 1) {
    const round = ROUND_PLAN[i];
    const count = winnersByRound.get(round.id)?.length || 0;
    if (count < round.count) {
      return i;
    }
  }
  return ROUND_PLAN.length;
}

function applyExistingWinners(existingWinners) {
  if (!Array.isArray(existingWinners)) {
    return;
  }

  const normalized = existingWinners
    .map(normalizeExistingWinner)
    .filter(Boolean);

  if (normalized.length === 0) {
    setExistingWinnersStatus("Belum ada pemenang tersimpan.");
    return;
  }

  normalized.sort((a, b) => {
    const aTime = a.wonAt ? a.wonAt.getTime() : 0;
    const bTime = b.wonAt ? b.wonAt.getTime() : 0;
    if (aTime !== bTime) {
      return aTime - bTime;
    }
    return 0;
  });

  const winnerIds = new Set(normalized.map((item) => item.id).filter(Boolean));
  const filterPool = (list) =>
    list.filter((item) => !winnerIds.has(String(item.id)));

  state.pool.regular.organik = filterPool(state.pool.regular.organik);
  state.pool.regular.tad = filterPool(state.pool.regular.tad);
  state.pool.guaranteed.organik = filterPool(state.pool.guaranteed.organik);
  state.pool.guaranteed.tad = filterPool(state.pool.guaranteed.tad);
  state.guaranteedRemaining =
    state.pool.guaranteed.organik.length + state.pool.guaranteed.tad.length;

  const doorWinners = normalized.filter((item) => item.prizeType === "door");
  const grandWinners = normalized.filter((item) => item.prizeType === "grand");
  const tadDoorWinners = doorWinners.filter(
    (item) => item.type === "TAD",
  ).length;

  state.doorPrizeDraws = doorWinners.length;
  state.grandPrizeDraws = grandWinners.length;
  state.remainingTad = Math.max(0, TAD_TARGET - tadDoorWinners);

  state.results = normalized.map((item) => ({
    id: item.id,
    name: item.name,
    position: item.position,
    branch: item.branch,
    type: item.type,
  }));

  const winnersByRound = new Map();
  normalized.forEach((item) => {
    if (!item.roundId) return;
    const list = winnersByRound.get(item.roundId) || [];
    list.push(item);
    winnersByRound.set(item.roundId, list);
  });

  elements.results.innerHTML = "";
  ROUND_PLAN.forEach((round) => {
    const list = winnersByRound.get(round.id);
    if (list && list.length) {
      renderRoundResults(round, list);
    }
  });

  state.currentRound = getNextRoundIndex(winnersByRound);
  updateQuota();
  updateRoundList();
  updateRoundUI();
  updateSummary();
  updateSpinAvailability();
  setExistingWinnersStatus(`Memuat ${normalized.length} pemenang tersimpan.`);
}

function initPool(records) {
  state.allEmployees = records;
  const eligible = records.filter((item) => !item.isExcluded);
  const guaranteed = eligible.filter((item) => item.guaranteedDoorprize);
  const regular = eligible.filter((item) => !item.guaranteedDoorprize);

  state.employees = eligible;
  state.pool.regular.organik = regular.filter(
    (item) => item.type === "ORGANIK",
  );
  state.pool.regular.tad = regular.filter((item) => item.type === "TAD");
  state.pool.guaranteed.organik = guaranteed.filter(
    (item) => item.type === "ORGANIK",
  );
  state.pool.guaranteed.tad = guaranteed.filter((item) => item.type === "TAD");
  state.remainingTad = TAD_TARGET;
  state.guaranteedRemaining = guaranteed.length;
  state.doorPrizeDraws = 0;
  state.grandPrizeDraws = 0;
  state.results = [];
  state.currentRound = 0;
  state.spinning = false;

  elements.totalEmployees.textContent = eligible.length;
  elements.totalOrganik.textContent =
    state.pool.regular.organik.length + state.pool.guaranteed.organik.length;
  elements.totalTad.textContent =
    state.pool.regular.tad.length + state.pool.guaranteed.tad.length;
  elements.totalWinners.textContent = "0";
  updateQuota();
  updateRoundList();
  updateRoundUI();
  elements.results.innerHTML = "";
  elements.resultsSummary.textContent =
    "Belum ada pemenang. Tekan spin untuk mulai.";
  elements.resetBtn.disabled = false;
  setExistingWinnersStatus("");
  updateSpinAvailability();
}

function validateCounts() {
  const remainingDraws = TOTAL_WINNERS - state.results.length;
  const remainingGrandDraws = GRAND_PRIZE_TOTAL - state.grandPrizeDraws;
  const availableTad =
    state.pool.regular.tad.length + state.pool.guaranteed.tad.length;
  const availableOrganik =
    state.pool.regular.organik.length + state.pool.guaranteed.organik.length;
  const availableTotal = availableTad + availableOrganik;
  const availableGrandEligible = state.pool.regular.organik.length;
  const remainingOrganik = remainingDraws - state.remainingTad;
  const currentRound = ROUND_PLAN[state.currentRound];
  const remainingDoorPrizeDraws = DOORPRIZE_TOTAL - state.doorPrizeDraws;
  const guaranteedOrganik = state.pool.guaranteed.organik.length;

  if (availableTotal < remainingDraws) {
    return `Data karyawan kurang. Minimal ${remainingDraws}, tersedia ${availableTotal}.`;
  }
  if (state.guaranteedRemaining > remainingDoorPrizeDraws) {
    return `Guaranteed doorprize melebihi sisa doorprize (${remainingDoorPrizeDraws}).`;
  }
  if (state.remainingTad > remainingDoorPrizeDraws) {
    return `Quota TAD tersisa ${state.remainingTad} melebihi sisa doorprize (${remainingDoorPrizeDraws}).`;
  }
  if (state.remainingTad + guaranteedOrganik > remainingDoorPrizeDraws) {
    return `Kombinasi TAD (${state.remainingTad}) dan guaranteed organik (${guaranteedOrganik}) melebihi sisa doorprize (${remainingDoorPrizeDraws}).`;
  }
  if (availableGrandEligible < remainingGrandDraws) {
    return `Data organik non-guaranteed kurang untuk grand prize. Minimal ${remainingGrandDraws}, tersedia ${availableGrandEligible}.`;
  }
  if (availableTad < state.remainingTad) {
    return `Data TAD kurang. Minimal ${state.remainingTad}, tersedia ${availableTad}.`;
  }
  if (availableOrganik < remainingOrganik) {
    return `Data organik kurang. Minimal ${remainingOrganik}, tersedia ${availableOrganik}.`;
  }
  if (
    currentRound &&
    currentRound.type === "grand" &&
    state.guaranteedRemaining > 0
  ) {
    return "Guaranteed doorprize belum terpenuhi sebelum grand prize.";
  }
  return "";
}

function updateQuota() {
  const tadWinners = TAD_TARGET - state.remainingTad;
  elements.quotaText.textContent = `${tadWinners} / ${TAD_TARGET}`;
  const percent = (tadWinners / TAD_TARGET) * 100;
  elements.quotaFill.style.width = `${percent}%`;
  const remaining = TOTAL_WINNERS - state.results.length;
  elements.remainingText.textContent = `Sisa spin: ${remaining}`;
  elements.totalWinners.textContent = state.results.length;
}

function updateRoundList() {
  elements.roundList.innerHTML = "";
  ROUND_PLAN.forEach((round, index) => {
    const item = document.createElement("li");
    item.className = "round-item";
    if (index === state.currentRound) {
      item.classList.add("active");
    }
    if (index < state.currentRound) {
      item.classList.add("done");
    }
    item.innerHTML = `
      <span>${round.label}</span>
      <span class="round-count">${round.count} nama</span>
    `;
    elements.roundList.appendChild(item);
  });
}

function updateRoundUI() {
  const round = ROUND_PLAN[state.currentRound];
  if (!round) {
    elements.currentRoundLabel.textContent = "Semua hadiah selesai";
    elements.currentRoundSub.textContent = "Total 85 pemenang sudah keluar.";
    elements.roundHint.textContent = "";
    elements.spinnerName.textContent = "Selesai";
    elements.spinnerSub.textContent =
      "Terima kasih dan selamat untuk semua pemenang.";
    return;
  }
  elements.currentRoundLabel.textContent = round.label;
  elements.currentRoundSub.textContent = `Siap undi ${round.count} pemenang.`;
  elements.spinBtn.textContent = `Putar ${round.label}`;
  elements.roundHint.textContent = `Spin ini akan menghasilkan ${round.count} nama.`;
}

function updateSpinAvailability() {
  const roundAvailable = state.currentRound < ROUND_PLAN.length;
  const hasData = state.employees.length > 0;
  const isValid = !validateCounts();
  const isBusy = state.spinning;
  elements.spinBtn.disabled = !(
    roundAvailable &&
    hasData &&
    isValid &&
    !isBusy
  );
}

function getDisplayPool() {
  return [
    ...state.pool.regular.organik,
    ...state.pool.regular.tad,
    ...state.pool.guaranteed.organik,
    ...state.pool.guaranteed.tad,
  ];
}

function getRandomFromPool() {
  const pool = getDisplayPool();
  if (pool.length === 0) {
    return { name: "-", position: "", branch: "" };
  }
  const idx = Math.floor(Math.random() * pool.length);
  return pool[idx];
}

function getDistinctNames(count) {
  const pool = getDisplayPool();
  if (pool.length === 0) {
    return Array.from({ length: count }, () => "-");
  }
  if (pool.length <= count) {
    const names = pool.map((item) => item.name);
    while (names.length < count) {
      names.push(pool[Math.floor(Math.random() * pool.length)].name);
    }
    return names;
  }
  const picked = [];
  const used = new Set();
  while (picked.length < count) {
    const idx = Math.floor(Math.random() * pool.length);
    if (!used.has(idx)) {
      used.add(idx);
      picked.push(pool[idx].name);
    }
  }
  return picked;
}

function drawFromGroup(group, remainingDoorPrizeDraws) {
  if (!group || (group.tad.length === 0 && group.organik.length === 0)) {
    return null;
  }
  if (remainingDoorPrizeDraws <= 0) {
    return null;
  }
  const mustTad = state.remainingTad === remainingDoorPrizeDraws;
  const mustOrganik = state.remainingTad === 0;
  if (mustTad && group.tad.length === 0) {
    return null;
  }
  if (mustOrganik && group.organik.length === 0) {
    return null;
  }
  let pickTad = mustTad
    ? true
    : mustOrganik
      ? false
      : Math.random() < state.remainingTad / remainingDoorPrizeDraws;

  if (pickTad && group.tad.length === 0) {
    pickTad = false;
  }
  if (!pickTad && group.organik.length === 0) {
    pickTad = true;
  }

  const pool = pickTad ? group.tad : group.organik;
  if (pool.length === 0) {
    return null;
  }
  const idx = Math.floor(Math.random() * pool.length);
  const winner = pool.splice(idx, 1)[0];
  if (pickTad) {
    state.remainingTad -= 1;
  }
  return winner;
}

function drawFromRegularForDoorPrize(remainingDoorPrizeDraws) {
  const group = state.pool.regular;
  if (!group || (group.tad.length === 0 && group.organik.length === 0)) {
    return null;
  }
  if (remainingDoorPrizeDraws <= 0) {
    return null;
  }

  const remainingGrandDraws = GRAND_PRIZE_TOTAL - state.grandPrizeDraws;
  const mustTad = state.remainingTad === remainingDoorPrizeDraws;
  const mustOrganik = state.remainingTad === 0;
  if (mustTad && group.tad.length === 0) {
    return null;
  }
  if (mustOrganik && group.organik.length === 0) {
    return null;
  }
  let pickTad = mustTad
    ? true
    : mustOrganik
      ? false
      : Math.random() < state.remainingTad / remainingDoorPrizeDraws;

  if (
    !pickTad &&
    state.remainingTad + state.pool.guaranteed.organik.length >=
      remainingDoorPrizeDraws
  ) {
    pickTad = true;
  }

  if (!pickTad && group.organik.length - 1 < remainingGrandDraws) {
    pickTad = true;
  }

  if (pickTad && group.tad.length === 0) {
    return null;
  }
  if (!pickTad && group.organik.length === 0) {
    return null;
  }

  const pool = pickTad ? group.tad : group.organik;
  if (pool.length === 0) {
    return null;
  }
  const idx = Math.floor(Math.random() * pool.length);
  const winner = pool.splice(idx, 1)[0];
  if (pickTad) {
    state.remainingTad -= 1;
  }
  return winner;
}

function drawDoorPrize(remainingDraws) {
  const remainingDoorPrizeDraws = DOORPRIZE_TOTAL - state.doorPrizeDraws;
  if (remainingDoorPrizeDraws <= 0) {
    return null;
  }
  const mustGuaranteed = state.guaranteedRemaining === remainingDoorPrizeDraws;
  let pickGuaranteed = mustGuaranteed
    ? true
    : state.guaranteedRemaining === 0
      ? false
      : Math.random() < state.guaranteedRemaining / remainingDoorPrizeDraws;

  const group = pickGuaranteed ? state.pool.guaranteed : state.pool.regular;
  let usedGuaranteed = pickGuaranteed;
  let winner = pickGuaranteed
    ? drawFromGroup(group, remainingDoorPrizeDraws)
    : drawFromRegularForDoorPrize(remainingDoorPrizeDraws);

  if (!winner) {
    usedGuaranteed = !pickGuaranteed;
    winner = pickGuaranteed
      ? drawFromRegularForDoorPrize(remainingDoorPrizeDraws)
      : drawFromGroup(state.pool.guaranteed, remainingDoorPrizeDraws);
  }

  if (!winner) {
    return null;
  }

  if (usedGuaranteed) {
    state.guaranteedRemaining -= 1;
  }
  state.doorPrizeDraws += 1;
  return winner;
}

function drawGrandPrize(remainingDraws) {
  const pool = state.pool.regular.organik;
  if (pool.length === 0) {
    return null;
  }
  const idx = Math.floor(Math.random() * pool.length);
  const winner = pool.splice(idx, 1)[0];
  state.grandPrizeDraws += 1;
  return winner;
}

function renderRoundResults(round, winners) {
  const group = document.createElement("div");
  group.className = "result-group";
  if (round.type === "grand") {
    group.classList.add("grand");
  }
  group.innerHTML = `
    <div class="result-title">
      <h3>${round.label}</h3>
      <span class="round-count">${winners.length} pemenang</span>
    </div>
  `;

  const grid = document.createElement("div");
  grid.className = "result-grid";

  winners.forEach((winner, index) => {
    const card = document.createElement("div");
    card.className = "winner-card";
    card.style.animationDelay = `${index * 0.03}s`;
    card.innerHTML = `
      <div class="winner-name">${winner.name}</div>
      <div class="winner-meta">${winner.position} - ${winner.branch}</div>
    `;
    grid.appendChild(card);
  });

  group.appendChild(grid);
  elements.results.appendChild(group);
}

async function syncWinners(round, winners) {
  if (!round || !winners || winners.length === 0) {
    return;
  }

  const payload = {
    round_id: round.id,
    round_label: round.label,
    prize_type: round.type,
    winners: winners.map((winner) => ({
      employee_id: winner.id ? String(winner.id) : "",
      name: winner.name,
      position: winner.position,
      branch: winner.branch,
      employment_type: winner.type,
    })),
  };

  try {
    const response = await fetchWithKey("/api/winners", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });
    if (!response.ok) {
      throw new Error("sync failed");
    }
  } catch (error) {
    console.warn("Failed to sync winners", error);
  }
}

function updateSummary() {
  if (state.results.length === 0) {
    elements.resultsSummary.textContent =
      "Belum ada pemenang. Tekan spin untuk mulai.";
    return;
  }
  elements.resultsSummary.textContent = `Total ${state.results.length} pemenang.`;
}

function wait(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function buildSpinRows(rowCount) {
  if (!elements.spinReel) return;
  elements.spinReel.innerHTML = "";
  spinRows = [];
  for (let i = 0; i < rowCount; i += 1) {
    const row = document.createElement("div");
    row.className = "spin-row";
    row.style.setProperty("--delay", `${i * 0.08}s`);

    const left = document.createElement("div");
    left.className = "spin-side left";
    left.textContent = "-";

    const center = document.createElement("div");
    center.className = "spin-center";
    center.textContent = "-";

    const right = document.createElement("div");
    right.className = "spin-side right";
    right.textContent = "-";

    row.appendChild(left);
    row.appendChild(center);
    row.appendChild(right);
    spinRows.push({ row, left, center, right });
    elements.spinReel.appendChild(row);
  }
}

function setActiveSpinRow(activeIndex) {
  spinRows.forEach((entry, index) => {
    entry.row.classList.toggle("is-active", index === activeIndex);
  });
}

function fillSpinRow(entry) {
  const [left, center, right] = getDistinctNames(3);
  entry.left.textContent = left;
  entry.center.textContent = center;
  entry.right.textContent = right;
}

function fillAllSpinRows() {
  if (!spinRows.length) return;
  spinRows.forEach((entry) => {
    fillSpinRow(entry);
    entry.row.classList.remove("is-winner");
    entry.row.classList.remove("is-active");
  });
}

function applySpinWinner(entry, winner) {
  entry.row.classList.add("is-winner");
  entry.center.textContent = winner.name;
}

function setOverlayVisible(isVisible) {
  if (!elements.spinOverlay) return;
  if (isVisible) {
    elements.spinOverlay.classList.add("is-visible");
    elements.spinOverlay.setAttribute("aria-hidden", "false");
  } else {
    elements.spinOverlay.classList.remove("is-visible");
    elements.spinOverlay.setAttribute("aria-hidden", "true");
  }
}

function getChunkSizes(count) {
  if (count <= CHUNK_SIZE) {
    return [count];
  }
  const chunks = [];
  let remaining = count;
  while (remaining > CHUNK_SIZE) {
    chunks.push(CHUNK_SIZE);
    remaining -= CHUNK_SIZE;
  }
  if (remaining === 1 && chunks.length > 0) {
    chunks[chunks.length - 1] += 1;
  } else {
    chunks.push(remaining);
  }
  return chunks;
}

function waitForChunkTrigger(label, hint) {
  if (!elements.nextChunkBtn) {
    return Promise.resolve();
  }
  elements.nextChunkBtn.disabled = false;
  elements.nextChunkBtn.textContent = label;
  if (elements.overlayHint) {
    elements.overlayHint.textContent = hint;
  }
  return new Promise((resolve) => {
    elements.nextChunkBtn.addEventListener(
      "click",
      () => {
        elements.nextChunkBtn.disabled = true;
        elements.nextChunkBtn.textContent = "Spinning...";
        if (elements.overlayHint) {
          elements.overlayHint.textContent = "Sedang memilih nama...";
        }
        resolve();
      },
      { once: true },
    );
  });
}

async function waitForDone() {
  if (!elements.nextChunkBtn) {
    return Promise.resolve();
  }
  elements.nextChunkBtn.disabled = false;
  elements.nextChunkBtn.textContent = "Done";
  if (elements.overlayHint) {
    elements.overlayHint.textContent = "";
  }
  return new Promise((resolve) => {
    elements.nextChunkBtn.addEventListener(
      "click",
      () => {
        elements.nextChunkBtn.disabled = true;
        resolve();
      },
      { once: true },
    );
  });
}

async function runChunkSpin(
  round,
  chunkIndex,
  totalChunks,
  chunkSize,
  roundWinners,
  spinCounter,
) {
  buildSpinRows(chunkSize);
  fillAllSpinRows();

  const activeRows = Array.from({ length: chunkSize }, () => true);
  const spinInterval = setInterval(() => {
    spinRows.forEach((entry, index) => {
      if (activeRows[index]) {
        fillSpinRow(entry);
      }
    });
  }, SPIN_TICK_MS);

  elements.spinFrame.classList.add("is-spinning");
  await wait(CHUNK_DURATION);

  for (let i = 0; i < chunkSize; i += 1) {
    if (round.type === "grand") {
      await wait(GRAND_CHUNK_PAUSE);
    }

    spinCounter.value += 1;
    const spinNumber = spinCounter.value;
    const chunkLabel = `Chunk ${chunkIndex + 1}/${totalChunks} · Spin ${
      i + 1
    }/${chunkSize}`;
    elements.overlaySub.textContent = chunkLabel;

    setActiveSpinRow(i);
    const directionClass = spinNumber % 2 === 1 ? "spin-left" : "spin-right";
    elements.spinFrame.classList.remove("spin-left", "spin-right");
    elements.spinFrame.classList.add(directionClass);
    triggerConfetti(CHUNK_PAUSE);

    activeRows[i] = false;
    const entry = spinRows[i];

    const remainingDraws = TOTAL_WINNERS - state.results.length;
    const winner =
      round.type === "grand"
        ? drawGrandPrize(remainingDraws)
        : drawDoorPrize(remainingDraws);
    if (!winner) {
      clearInterval(spinInterval);
      elements.spinFrame.classList.remove(
        "is-spinning",
        "spin-left",
        "spin-right",
      );
      throw new Error("Gagal memilih pemenang. Periksa data karyawan.");
    }
    state.results.push(winner);
    roundWinners.push(winner);
    applySpinWinner(entry, winner);
    updateQuota();
    await wait(CHUNK_PAUSE);
  }

  clearInterval(spinInterval);
  elements.spinFrame.classList.remove("is-spinning", "spin-left", "spin-right");
}

function triggerConfetti(durationMs) {
  if (!elements.confetti) return;
  elements.confetti.classList.add("is-active");
  if (confettiTimer) {
    clearTimeout(confettiTimer);
  }
  const activeFor = Math.min(4200, durationMs + 1200);
  confettiTimer = setTimeout(() => {
    elements.confetti.classList.remove("is-active");
  }, activeFor);
}

async function spinRound() {
  if (state.spinning) return;
  const round = ROUND_PLAN[state.currentRound];
  elements.overlayTitle.textContent = round.label;
  if (!round) return;

  const validation = validateCounts();
  if (validation) {
    setAlert(validation);
    return;
  }

  setAlert("");
  state.spinning = true;
  updateSpinAvailability();
  elements.spinnerName.textContent = "Memutar...";
  elements.spinnerSub.textContent = `Menentukan ${round.count} nama.`;

  const chunks = getChunkSizes(round.count);
  setOverlayVisible(true);

  const roundWinners = [];
  const spinCounter = { value: 0 };
  try {
    for (let i = 0; i < chunks.length; i += 1) {
      const label = i === 0 ? "Start spin" : "Next spin";
      const hint = `Tekan untuk mulai chunk ${i + 1} dari ${chunks.length}.`;
      elements.overlaySub.textContent = `Chunk ${i + 1}/${chunks.length} · ${
        chunks[i]
      } nama`;
      await waitForChunkTrigger(label, hint);
      await runChunkSpin(
        round,
        i,
        chunks.length,
        chunks[i],
        roundWinners,
        spinCounter,
      );
    }
  } catch (error) {
    setAlert(
      error instanceof Error ? error.message : "Gagal memilih pemenang.",
    );
    setOverlayVisible(false);
    state.spinning = false;
    updateSpinAvailability();
    return;
  }

  await waitForDone();
  setOverlayVisible(false);
  if (elements.overlayHint) {
    elements.overlayHint.textContent = "";
  }
  elements.spinnerName.textContent = "Pemenang Muncul!";
  elements.spinnerSub.textContent = `${round.count} nama telah dipilih.`;

  renderRoundResults(round, roundWinners);
  syncWinners(round, roundWinners);
  state.currentRound += 1;
  state.spinning = false;
  updateQuota();
  updateRoundList();
  updateRoundUI();
  updateSummary();
  updateSpinAvailability();
}

function resetAll() {
  if (state.spinning) return;
  setOverlayVisible(false);
  elements.spinnerName.textContent = "Ready";
  elements.spinnerSub.textContent = "Doorprize dimulai setelah data siap.";
  if (elements.nextChunkBtn) {
    elements.nextChunkBtn.disabled = true;
    elements.nextChunkBtn.textContent = "Next spin";
  }
  if (elements.overlayHint) {
    elements.overlayHint.textContent = "";
  }
  if (elements.confetti) {
    elements.confetti.classList.remove("is-active");
  }
  if (confettiTimer) {
    clearTimeout(confettiTimer);
    confettiTimer = null;
  }
  setAlert("");
  loadEmployees();
}

async function loadExistingWinners() {
  try {
    const response = await fetchWithKey("/api/winners");
    if (!response.ok) {
      throw new Error("bad_response");
    }
    const data = await response.json();
    return Array.isArray(data) ? data : [];
  } catch (error) {
    return null;
  }
}

async function loadEmployees() {
  setStatus("Mengambil data...");
  setAlert("");
  elements.spinBtn.disabled = true;
  elements.resetBtn.disabled = true;

  try {
    const response = await fetchWithKey("/api/employees/all");
    if (!response.ok) {
      throw new Error("bad_response");
    }
    const data = await response.json();

    const rawList = Array.isArray(data) ? data : data.data;
    if (!Array.isArray(rawList)) {
      setStatus("Format data tidak dikenali.");
      setAlert("Gunakan array JSON dengan kolom sesuai database.");
      return;
    }

    const normalized = [];
    let skipped = 0;
    rawList.forEach((row) => {
      const item = normalizeEmployee(row);
      if (item) {
        normalized.push(item);
      } else {
        skipped += 1;
      }
    });

    if (skipped > 0) {
      setAlert(
        `Ada ${skipped} data dilewati karena jenis kepegawaian tidak valid.`,
      );
    }

    initPool(normalized);
    setExistingWinnersStatus("Memuat pemenang tersimpan...");
    const existingWinners = await loadExistingWinners();
    if (existingWinners === null) {
      setExistingWinnersStatus("Gagal mengambil data pemenang.");
    } else {
      applyExistingWinners(existingWinners);
    }

    const validation = validateCounts();
    if (validation) {
      setAlert(validation);
      setStatus("Data tidak cukup untuk menjalankan spin.");
      updateSpinAvailability();
      return;
    }

    setStatus("Data siap dari /api/employees/all.");
    updateSpinAvailability();
  } catch (error) {
    setAlert("Gagal mengambil data dari /api/employees/all.");
    if (state.employees.length) {
      setStatus("Gagal memuat data baru. Menggunakan data sebelumnya.");
      elements.resetBtn.disabled = false;
      updateSpinAvailability();
    } else {
      setStatus("Data belum tersedia.");
    }
  }
}

function init() {
  updateRoundList();
  elements.spinBtn.addEventListener("click", spinRound);
  elements.resetBtn.addEventListener("click", resetAll);
  elements.reloadBtn.addEventListener("click", loadEmployees);
  loadEmployees();
}

init();
