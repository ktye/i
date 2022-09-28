this.addEventListener("install", (event) => {
 console.log("sw install");
 event.waitUntil(
  caches
   .open("v1")
   .then((cache) =>
    cache.addAll([
     "/kweb/a.html",
     "/k.js",
     "/kweb/kweb.js",
     "/kweb/plot.js",
     "/kweb/a.k"
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
