const API_KEY = "Dptaspen@25!";

const contentEl = document.getElementById("scan-card-content");
const videoEl = document.getElementById("qr-video");

let currentData = "";
let clearTimer = null;
let scannerId = null;

function isValidScanner(value) {
  const number = Number(value);
  return Number.isInteger(number) && number >= 1 && number <= 3;
}

function resolveScannerId() {
  const params = new URLSearchParams(window.location.search);
  const fromQuery = params.get("scanner");
  if (isValidScanner(fromQuery)) {
    return Number(fromQuery);
  }

  let input = "";
  while (!isValidScanner(input)) {
    input = window.prompt("Masukkan scanner ID (1-3):") || "";
  }
  return Number(input);
}

function setScanData(text) {
  if (text === currentData) {
    return;
  }

  currentData = text;
  contentEl.textContent = text;

  if (clearTimer) {
    clearTimeout(clearTimer);
    clearTimer = null;
  }

  if (text) {
    clearTimer = setTimeout(() => setScanData(""), 10000);
    storeAttendance(text);
  }
}

async function storeAttendance(name) {
  try {
    await fetch("/api/scans", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-API-Key": API_KEY,
      },
      body: JSON.stringify({ name, scanner_id: scannerId }),
    });
  } catch (error) {
    console.error("Failed to store attendance:", error);
  }
}

function startScanner() {
  if (!scannerId) {
    scannerId = resolveScannerId();
  }

  if (!videoEl) {
    console.error("Video element not found.");
    return;
  }

  if (!window.ZXing || !ZXing.BrowserMultiFormatReader) {
    console.error("ZXing library not available.");
    return;
  }

  const reader = new ZXing.BrowserMultiFormatReader();
  reader.decodeFromVideoDevice(null, videoEl, (result, err) => {
    if (result) {
      setScanData(result.getText());
    } else if (err && err.name !== "NotFoundException") {
      console.error("QR scan error:", err);
    }
  });
}

startScanner();
