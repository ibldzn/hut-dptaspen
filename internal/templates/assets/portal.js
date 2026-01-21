const API_KEY = "Dptaspen@25!";

const nipForm = document.getElementById("nipForm");
const nipInput = document.getElementById("nipInput");
const modal = document.getElementById("errorModal");
const closeModal = document.getElementById("closeModal");
const dismissBtn = document.getElementById("dismissBtn");

const normalizeNip = (value) => value.trim().replace(/\s+/g, "");

const showModal = () => {
  modal.classList.add("is-visible");
  modal.setAttribute("aria-hidden", "false");
  closeModal.focus();
};

const hideModal = () => {
  modal.classList.remove("is-visible");
  modal.setAttribute("aria-hidden", "true");
  nipInput.focus();
};

const revealItems = document.querySelectorAll(".reveal");
const observer = new IntersectionObserver(
  (entries, obs) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add("is-visible");
        obs.unobserve(entry.target);
      }
    });
  },
  {
    threshold: 0.2,
    rootMargin: "0px 0px -10% 0px",
  },
);

revealItems.forEach((item) => observer.observe(item));

async function lookupInvitation(nip) {
  const response = await fetch(
    `/api/invitations/lookup?nip=${encodeURIComponent(nip)}`,
    {
      headers: {
        "X-API-Key": API_KEY,
      },
    },
  );
  if (!response.ok) {
    const error = new Error("lookup_failed");
    error.status = response.status;
    throw error;
  }
  return response.json();
}

nipForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const normalized = normalizeNip(nipInput.value);

  if (!normalized) {
    nipInput.focus();
    return;
  }

  nipInput.value = normalized;

  try {
    const data = await lookupInvitation(normalized);
    if (data && data.url) {
      window.location.href = data.url;
      return;
    }
    showModal();
  } catch (error) {
    showModal();
  }
});

closeModal.addEventListener("click", hideModal);
dismissBtn.addEventListener("click", hideModal);
modal.addEventListener("click", (event) => {
  if (event.target === modal) {
    hideModal();
  }
});
window.addEventListener("keydown", (event) => {
  if (event.key === "Escape" && modal.classList.contains("is-visible")) {
    hideModal();
  }
});
