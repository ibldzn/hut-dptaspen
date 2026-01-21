const API_KEY = "Dptaspen@25!";

const elements = {
  totalEmployees: document.getElementById("totalEmployees"),
  totalPresent: document.getElementById("totalPresent"),
  totalAbsent: document.getElementById("totalAbsent"),
  presentRate: document.getElementById("presentRate"),
  attendanceTimestamp: document.getElementById("attendanceTimestamp"),
  attendanceSummary: document.getElementById("attendanceSummary"),
  attendanceStatus: document.getElementById("attendanceStatus"),
  attendanceRows: document.getElementById("attendanceRows"),
  searchInput: document.getElementById("searchInput"),
  branchFilter: document.getElementById("branchFilter"),
  statusFilter: document.getElementById("statusFilter"),
  refreshAttendance: document.getElementById("refreshAttendance"),
  tabButtons: Array.from(document.querySelectorAll(".tab-btn")),
  attendancePanes: Array.from(document.querySelectorAll(".attendance-pane")),
  guestSearchInput: document.getElementById("guestSearchInput"),
  guestNameInput: document.getElementById("guestNameInput"),
  addGuestBtn: document.getElementById("addGuestBtn"),
  guestRows: document.getElementById("guestRows"),
  guestStatus: document.getElementById("guestStatus"),
  guestSummary: document.getElementById("guestSummary"),
  exportWinnersBtn: document.getElementById("exportWinnersBtn"),
  exportAttendanceBtn: document.getElementById("exportAttendanceBtn"),
  resetWinnersBtn: document.getElementById("resetWinnersBtn"),
  resetAttendanceBtn: document.getElementById("resetAttendanceBtn"),
  exportStatus: document.getElementById("exportStatus"),
  winnerCount: document.getElementById("winnerCount"),
  attendanceCount: document.getElementById("attendanceCount"),
};

const state = {
  employees: [],
  guests: [],
  activeTab: "employees",
};

function withAPIKey(options = {}) {
  const headers = new Headers(options.headers || {});
  headers.set("X-API-Key", API_KEY);
  return { ...options, headers };
}

async function fetchJSON(url, options) {
  const response = await fetch(url, withAPIKey(options));
  if (!response.ok) {
    throw new Error(`Request gagal (${response.status})`);
  }
  return response.json();
}

function parseDate(raw) {
  if (!raw) return null;
  const date = new Date(raw);
  if (Number.isNaN(date.getTime())) {
    return null;
  }
  return date;
}

function normalizeEmployee(raw) {
  if (!raw) return null;
  const presentAt = parseDate(
    raw.PRESENT_AT ?? raw.present_at ?? raw.presentAt,
  );
  return {
    name: raw.NAMA_KARYAWAN ?? raw.nama_karyawan ?? raw.name ?? "-",
    position: raw.JABATAN ?? raw.jabatan ?? raw.position ?? "-",
    branch: raw.KANTOR_CABANG ?? raw.kantor_cabang ?? raw.branch ?? "-",
    type: raw.JENIS_KEPEGAWAIAN ?? raw.jenis_kepegawaian ?? raw.type ?? "-",
    presentAt,
    present: Boolean(presentAt),
  };
}

function formatDate(date) {
  if (!date) return "-";
  return date.toLocaleString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function setAttendanceStatus(message) {
  elements.attendanceStatus.textContent = message;
}

function setGuestStatus(message) {
  if (!elements.guestStatus) return;
  elements.guestStatus.textContent = message;
}

function setActiveTab(tab) {
  state.activeTab = tab;
  elements.tabButtons.forEach((btn) => {
    btn.classList.toggle("is-active", btn.dataset.tab === tab);
  });
  elements.attendancePanes.forEach((pane) => {
    pane.classList.toggle("is-active", pane.dataset.tab === tab);
  });
}

function updateSummary(list) {
  const total = list.length;
  const present = list.filter((item) => item.present).length;
  const absent = total - present;
  const rate = total ? Math.round((present / total) * 100) : 0;

  elements.totalEmployees.textContent = total;
  elements.totalPresent.textContent = present;
  elements.totalAbsent.textContent = absent;
  elements.presentRate.textContent = `${rate}%`;
}

function updateAttendanceCount() {
  const total = state.employees.length + state.guests.length;
  elements.attendanceCount.textContent = total;
}

function updateAttendanceTimestamp() {
  const now = new Date();
  elements.attendanceTimestamp.textContent = `Terakhir sinkron: ${formatDate(now)}`;
}

function renderRows(list) {
  elements.attendanceRows.innerHTML = "";
  const fragment = document.createDocumentFragment();

  list.forEach((item) => {
    const row = document.createElement("tr");
    const cells = [
      item.name,
      item.position,
      item.branch,
      item.type,
      null,
      formatDate(item.presentAt),
    ];
    cells.forEach((value, index) => {
      const cell = document.createElement("td");
      if (index === 4) {
        const pill = document.createElement("span");
        pill.className = `status-pill ${item.present ? "present" : "absent"}`;
        pill.textContent = item.present ? "Hadir" : "Belum hadir";
        cell.appendChild(pill);
      } else {
        cell.textContent = value;
      }
      row.appendChild(cell);
    });
    fragment.appendChild(row);
  });

  elements.attendanceRows.appendChild(fragment);
}

function applyFilters() {
  const search = elements.searchInput.value.trim().toLowerCase();
  const branch = elements.branchFilter.value;
  const status = elements.statusFilter.value;

  let filtered = [...state.employees];
  if (search) {
    filtered = filtered.filter((item) =>
      [item.name, item.position, item.branch]
        .join(" ")
        .toLowerCase()
        .includes(search),
    );
  }
  if (branch !== "all") {
    filtered = filtered.filter((item) => item.branch === branch);
  }
  if (status !== "all") {
    filtered = filtered.filter((item) =>
      status === "present" ? item.present : !item.present,
    );
  }

  renderRows(filtered);
  elements.attendanceSummary.textContent = `Menampilkan ${filtered.length} dari ${state.employees.length} karyawan.`;
}

function updateBranchFilter(list) {
  const branches = Array.from(new Set(list.map((item) => item.branch))).sort();
  elements.branchFilter.innerHTML = `<option value="all">Semua cabang</option>`;
  branches.forEach((branch) => {
    const option = document.createElement("option");
    option.value = branch;
    option.textContent = branch;
    elements.branchFilter.appendChild(option);
  });
}

function normalizeGuest(raw) {
  if (!raw) return null;
  const presentAt = parseDate(raw.PRESENT_AT ?? raw.present_at ?? raw.presentAt);
  return {
    id: raw.ID ?? raw.id ?? 0,
    name: raw.NAMA_TAMU ?? raw.nama_tamu ?? raw.name ?? "-",
    presentAt,
  };
}

function renderGuestRows(list) {
  if (!elements.guestRows) return;
  elements.guestRows.innerHTML = "";
  const fragment = document.createDocumentFragment();

  list.forEach((item) => {
    const row = document.createElement("tr");
    const nameCell = document.createElement("td");
    nameCell.textContent = item.name;
    const timeCell = document.createElement("td");
    timeCell.textContent = formatDate(item.presentAt);
    row.appendChild(nameCell);
    row.appendChild(timeCell);
    fragment.appendChild(row);
  });

  elements.guestRows.appendChild(fragment);
}

function applyGuestFilter() {
  if (!elements.guestSearchInput) return;
  const search = elements.guestSearchInput.value.trim().toLowerCase();
  let filtered = [...state.guests];
  if (search) {
    filtered = filtered.filter((item) =>
      item.name.toLowerCase().includes(search),
    );
  }

  renderGuestRows(filtered);
  if (elements.guestSummary) {
    elements.guestSummary.textContent = `Menampilkan ${filtered.length} dari ${state.guests.length} tamu.`;
  }
}

async function loadAttendance() {
  setAttendanceStatus("Memuat data kehadiran...");
  try {
    const data = await fetchJSON("/api/employees/all");
    const list = Array.isArray(data) ? data : data.data;
    if (!Array.isArray(list)) {
      throw new Error("Format data tidak valid");
    }
    const normalized = list.map(normalizeEmployee).filter(Boolean);
    state.employees = normalized;
    updateSummary(normalized);
    updateBranchFilter(normalized);
    applyFilters();
    updateAttendanceTimestamp();
    updateAttendanceCount();
    setAttendanceStatus("Data kehadiran siap.");
  } catch (error) {
    setAttendanceStatus(
      error instanceof Error ? error.message : "Gagal memuat data kehadiran.",
    );
  }
}

async function loadGuests() {
  setGuestStatus("Memuat data tamu...");
  try {
    const data = await fetchJSON("/api/guests");
    const list = Array.isArray(data) ? data : data.data;
    if (!Array.isArray(list)) {
      throw new Error("Format data tamu tidak valid");
    }
    const normalized = list.map(normalizeGuest).filter(Boolean);
    state.guests = normalized;
    applyGuestFilter();
    updateAttendanceCount();
    setGuestStatus("Data tamu siap.");
  } catch (error) {
    setGuestStatus(
      error instanceof Error ? error.message : "Gagal memuat data tamu.",
    );
  }
}

async function loadWinnerCounts() {
  try {
    const allWinners = await fetchJSON("/api/winners");
    elements.winnerCount.textContent = Array.isArray(allWinners)
      ? allWinners.length
      : 0;
  } catch (error) {
    elements.winnerCount.textContent = "-";
  }
}

async function downloadWinners() {
  elements.exportStatus.textContent = "Menyiapkan export pemenang...";
  try {
    const response = await fetch("/api/winners/export", withAPIKey());
    if (!response.ok) {
      throw new Error(`Export gagal (${response.status})`);
    }
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "winners-all.csv";
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
    elements.exportStatus.textContent = "Export pemenang selesai.";
  } catch (error) {
    elements.exportStatus.textContent =
      error instanceof Error ? error.message : "Export gagal.";
  }
}

async function downloadAttendance() {
  elements.exportStatus.textContent = "Menyiapkan export presensi...";
  try {
    const response = await fetch("/api/employees/export", withAPIKey());
    if (!response.ok) {
      throw new Error(`Export gagal (${response.status})`);
    }
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "attendance.csv";
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
    elements.exportStatus.textContent = "Export presensi selesai.";
  } catch (error) {
    elements.exportStatus.textContent =
      error instanceof Error ? error.message : "Export gagal.";
  }
}

async function resetWinners() {
  const confirmed = window.confirm(
    "Yakin ingin menghapus semua data pemenang? Tindakan ini tidak bisa dibatalkan.",
  );
  if (!confirmed) {
    return;
  }

  elements.exportStatus.textContent = "Menghapus data pemenang...";
  try {
    const response = await fetch(
      "/api/winners",
      withAPIKey({
        method: "DELETE",
      }),
    );
    if (!response.ok) {
      throw new Error(`Reset gagal (${response.status})`);
    }
    elements.exportStatus.textContent = "Semua pemenang berhasil dihapus.";
    await loadWinnerCounts();
  } catch (error) {
    elements.exportStatus.textContent =
      error instanceof Error ? error.message : "Reset pemenang gagal.";
  }
}

async function resetAttendance() {
  const confirmed = window.confirm(
    "Yakin ingin menghapus semua data kehadiran? Tindakan ini tidak bisa dibatalkan.",
  );
  if (!confirmed) {
    return;
  }

  setAttendanceStatus("Menghapus data kehadiran...");
  try {
    const response = await fetch(
      "/api/employees/present",
      withAPIKey({
        method: "DELETE",
      }),
    );
    if (!response.ok) {
      throw new Error(`Reset gagal (${response.status})`);
    }
    await loadAttendance();
    setAttendanceStatus("Semua status kehadiran berhasil dihapus.");
  } catch (error) {
    setAttendanceStatus(
      error instanceof Error ? error.message : "Reset kehadiran gagal.",
    );
  }
}

async function resetGuests() {
  const confirmed = window.confirm(
    "Yakin ingin menghapus semua data tamu? Tindakan ini tidak bisa dibatalkan.",
  );
  if (!confirmed) {
    return;
  }

  setGuestStatus("Menghapus data tamu...");
  try {
    const response = await fetch(
      "/api/guests/present",
      withAPIKey({
        method: "DELETE",
      }),
    );
    if (!response.ok) {
      throw new Error(`Reset gagal (${response.status})`);
    }
    await loadGuests();
    setGuestStatus("Semua data tamu berhasil dihapus.");
  } catch (error) {
    setGuestStatus(
      error instanceof Error ? error.message : "Reset tamu gagal.",
    );
  }
}

async function addGuest() {
  if (!elements.guestNameInput) return;
  const name = elements.guestNameInput.value.trim();
  if (!name) {
    setGuestStatus("Nama tamu wajib diisi.");
    return;
  }

  setGuestStatus("Menyimpan data tamu...");
  try {
    const response = await fetch(
      "/api/guests/mark_present",
      withAPIKey({
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name }),
      }),
    );
    if (!response.ok) {
      throw new Error(`Simpan gagal (${response.status})`);
    }
    elements.guestNameInput.value = "";
    await loadGuests();
    setGuestStatus("Tamu berhasil ditambahkan.");
  } catch (error) {
    setGuestStatus(
      error instanceof Error ? error.message : "Gagal menambah tamu.",
    );
  }
}

function resetActiveAttendance() {
  if (state.activeTab === "guests") {
    resetGuests();
    return;
  }
  resetAttendance();
}

function refreshActiveAttendance() {
  if (state.activeTab === "guests") {
    loadGuests();
    return;
  }
  loadAttendance();
}

function init() {
  elements.searchInput.addEventListener("input", applyFilters);
  elements.branchFilter.addEventListener("change", applyFilters);
  elements.statusFilter.addEventListener("change", applyFilters);
  if (elements.guestSearchInput) {
    elements.guestSearchInput.addEventListener("input", applyGuestFilter);
  }
  if (elements.addGuestBtn) {
    elements.addGuestBtn.addEventListener("click", addGuest);
  }
  if (elements.guestNameInput) {
    elements.guestNameInput.addEventListener("keydown", (event) => {
      if (event.key === "Enter") {
        addGuest();
      }
    });
  }
  elements.refreshAttendance.addEventListener("click", refreshActiveAttendance);
  elements.exportWinnersBtn.addEventListener("click", downloadWinners);
  elements.exportAttendanceBtn.addEventListener("click", downloadAttendance);
  elements.resetWinnersBtn.addEventListener("click", resetWinners);
  elements.resetAttendanceBtn.addEventListener("click", resetActiveAttendance);
  elements.tabButtons.forEach((btn) => {
    btn.addEventListener("click", () => {
      setActiveTab(btn.dataset.tab);
    });
  });

  setActiveTab(state.activeTab);
  loadAttendance();
  loadGuests();
  loadWinnerCounts();
}

init();
