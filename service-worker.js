const CACHE_NAME = 'my-htmx-go-app-v1';
const urlsToCache = [
  '/',
  '/styles.css',
  '/htmx.min.js'
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => cache.addAll(urlsToCache))
  );
});

self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request)
      .then((response) => response || fetch(event.request))
  );
});


let startY;
let pullDistance = 0;
const pullThreshold = 80;
const pullToRefreshElement = document.getElementById('pullToRefresh');

document.addEventListener('touchstart', (e) => {
  startY = e.touches[0].pageY;
});

document.addEventListener('touchmove', (e) => {
  const y = e.touches[0].pageY;
  pullDistance = y - startY;
  
  if (pullDistance > 0 && window.scrollY === 0) {
    pullToRefreshElement.style.transform = `translateY(${Math.min(pullDistance / 2, pullThreshold)}px)`;
    e.preventDefault();
  }
});

document.addEventListener('touchend', () => {
  if (pullDistance > pullThreshold) {
    refresh();
  }
  
  pullToRefreshElement.style.transform = 'translateY(-50px)';
  pullDistance = 0;
});

function refresh() {
  // Implement your refresh logic here
  console.log('Refreshing...');
  
  // For example, you could reload the page:
  location.reload();
  
  // Or fetch new data and update the DOM:
  // fetchNewData().then(updateDOM);
}