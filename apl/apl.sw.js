this.addEventListener("install", (event) => {
 console.log("sw install");
 event.waitUntil(
  caches
   .open("v1")
   .then((cache) =>
    cache.addAll([
     "/apl.html",
     "/k.js",
     "/k.wasm",
     "/apl.woff2",
     "/icon.svg",
     "/div.png"
    ])
   )
  )
})


self.addEventListener('fetch', function(event) {
 event.respondWith(
  caches.match(event.request).then(function(response) {
   return response || fetch(event.request);
  })
 );
});

self.addEventListener("activate",event=>{
 console.log("Service worker activated");
})
self.addEventListener("beforeinstallprompt",event=>{
 console.log("Service worker before install");
})
