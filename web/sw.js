// Service worker mínimo: por enquanto só existe para o app ser instalável.
// O cache offline (vias/ginásios) entra aqui depois.
self.addEventListener("install", () => self.skipWaiting());
self.addEventListener("activate", (event) => event.waitUntil(self.clients.claim()));
