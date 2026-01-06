/* KeyFlip site JS: tiny niceties, no framework, no drama. */
(function () {
  const $ = (sel, el = document) => el.querySelector(sel);
  const $$ = (sel, el = document) => Array.from(el.querySelectorAll(sel));

  // Active nav link (no brittle pathname assumptions)
  try {
    const page = document.body.getAttribute("data-page");
    if (page) {
      $$(".navlinks a[data-nav]").forEach(a => {
        a.classList.toggle("active", a.getAttribute("data-nav") === page);
      });
    }
  } catch {}

  // Reveal on view
  const reveals = $$(".reveal");
  if (reveals.length) {
    const io = new IntersectionObserver((entries) => {
      entries.forEach(e => {
        if (e.isIntersecting) e.target.classList.add("show");
      });
    }, { threshold: 0.12 });
    reveals.forEach(el => io.observe(el));
  }

  // Docs: highlight sidebar section on scroll
  const docNav = $(".sidebar nav");
  if (docNav) {
    const links = $$("a[href^='#']", docNav);
    const sections = links
      .map(a => document.getElementById(a.getAttribute("href").slice(1)))
      .filter(Boolean);

    const setActive = (id) => {
      links.forEach(a => a.classList.toggle("active", a.getAttribute("href") === "#" + id));
    };

    const io2 = new IntersectionObserver((entries) => {
      const visible = entries.filter(e => e.isIntersecting)
        .sort((a,b) => b.intersectionRatio - a.intersectionRatio)[0];
      if (visible) setActive(visible.target.id);
    }, { rootMargin: "-20% 0px -65% 0px", threshold: [0.1, 0.2, 0.3, 0.4, 0.5] });

    sections.forEach(s => io2.observe(s));
    links.forEach(a => a.addEventListener("click", (e) => {
      const id = a.getAttribute("href").slice(1);
      const target = document.getElementById(id);
      if (!target) return;
      e.preventDefault();
      target.scrollIntoView({ behavior: "smooth", block: "start" });
      history.replaceState(null, "", "#" + id);
    }));
  }

  // Download page: OS tabs
  const tabs = $$(".tab");
  const panels = $$(".download-card[data-os]");
  if (tabs.length && panels.length) {
    const key = "keyflip_os";
    const pick = (os) => {
      tabs.forEach(t => t.setAttribute("aria-selected", String(t.dataset.os === os)));
      panels.forEach(p => p.style.display = (p.dataset.os === os ? "block" : "none"));
      try { localStorage.setItem(key, os); } catch {}
    };
    const initial =
      (location.hash && location.hash.replace("#", "")) ||
      (() => { try { return localStorage.getItem(key); } catch { return null; } })() ||
      "macos";

    pick(["macos","windows","linux"].includes(initial) ? initial : "macos");

    tabs.forEach(t => t.addEventListener("click", () => {
      const os = t.dataset.os;
      pick(os);
      history.replaceState(null, "", "#" + os);
    }));
  }

  // Copy buttons (optional)
  $$(".copy").forEach(btn => {
    btn.addEventListener("click", async () => {
      const target = btn.getAttribute("data-copy");
      const el = target ? document.querySelector(target) : null;
      const text = el ? el.textContent.trim() : "";
      if (!text) return;
      try {
        await navigator.clipboard.writeText(text);
        const old = btn.textContent;
        btn.textContent = "Copied";
        setTimeout(() => (btn.textContent = old), 900);
      } catch {}
    });
  });
})();
