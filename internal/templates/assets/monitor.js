const API_KEY = "Dptaspen@25!";

const cards = new Map();
const state = {
  1: [],
  2: [],
  3: [],
};

function formatTime(value) {
  if (!value) return "--:--:--";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "--:--:--";
  }
  return date.toLocaleTimeString("id-ID", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
}

function hydrateCards() {
  document.querySelectorAll(".scan-card").forEach((card) => {
    const scannerId = Number(card.dataset.scanner);
    const entries = Array.from(card.querySelectorAll(".scan-entry")).map(
      (entry) => ({
        nameEl: entry.querySelector(".scan-name"),
        timeEl: entry.querySelector(".scan-time"),
      }),
    );
    cards.set(scannerId, { card, entries });
  });
}

function renderScanner(scannerId) {
  const card = cards.get(scannerId);
  if (!card) return;
  const events = state[scannerId] || [];
  card.entries.forEach((entry, index) => {
    const data = events[index];
    entry.nameEl.textContent = data ? data.name : "-";
    entry.timeEl.textContent = data ? formatTime(data.scanned_at) : "--:--:--";
  });
}

function applySnapshot(snapshot) {
  [1, 2, 3].forEach((scannerId) => {
    const key = String(scannerId);
    const events = Array.isArray(snapshot?.[key]) ? snapshot[key] : [];
    state[scannerId] = events.slice(0, 3);
    renderScanner(scannerId);
  });
}

async function loadSnapshot() {
  try {
    const response = await fetch("/api/scans/recent", {
      headers: {
        "X-API-Key": API_KEY,
      },
    });
    if (!response.ok) {
      throw new Error("snapshot_failed");
    }
    const data = await response.json();
    applySnapshot(data);
  } catch (error) {
    // keep placeholders
  }
}

function handleIncoming(event) {
  if (!event || !event.scanner_id) return;
  const scannerId = Number(event.scanner_id);
  if (![1, 2, 3].includes(scannerId)) return;
  const list = state[scannerId] || [];
  list.unshift(event);
  state[scannerId] = list.slice(0, 3);
  renderScanner(scannerId);
}

function connectSocket() {
  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const socket = new WebSocket(
    `${protocol}://${window.location.host}/ws/attendance?key=${API_KEY}`,
  );

  socket.addEventListener("message", (event) => {
    try {
      const payload = JSON.parse(event.data);
      handleIncoming(payload);
    } catch (error) {
      // ignore bad payload
    }
  });

  socket.addEventListener("close", () => {
    setTimeout(connectSocket, 2000);
  });

  socket.addEventListener("error", () => {
    socket.close();
  });
}

hydrateCards();
loadSnapshot();
connectSocket();
