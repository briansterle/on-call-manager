let startY;
let pullDistance = 0;
const pullThreshold = 20;
const indicator = document.getElementById('pull-to-refresh-indicator');
const pullMessage = document.getElementById('pull-message');
const refreshSpinner = document.getElementById('refresh-spinner');

// Get the safe area inset top value
const safeAreaInsetTop = parseInt(getComputedStyle(document.documentElement).getPropertyValue('--sat') || '0');

document.addEventListener('touchstart', (e) => {
    startY = e.touches[0].pageY;
});

document.addEventListener('touchmove', (e) => {
    if (window.scrollY === 0) {
        const y = e.touches[0].pageY;
        pullDistance = y - startY;
        
        if (pullDistance > 0) {
            indicator.style.transform = `translateY(calc(${Math.min(pullDistance / 2, 50)}px - 100% + ${safeAreaInsetTop}px))`;
            e.preventDefault();
        }
    }
});

document.addEventListener('touchend', () => {
    if (pullDistance > pullThreshold) {
        pullMessage.style.display = 'none';
        refreshSpinner.style.display = 'inline-block';
        indicator.style.transform = `translateY(calc(-50% + ${safeAreaInsetTop}px))`;
        // indicator.style.transform = `translateY(calc(-100% + ${safeAreaInsetTop}px))`;
        
        // Delay the refresh by 0.5 seconds
        setTimeout(() => {
            location.reload();
        }, 500);
    } else {
        indicator.style.transform = `translateY(calc(-50% + ${safeAreaInsetTop}px))`;
        // indicator.style.transform = `translateY(calc(-100% - ${safeAreaInsetTop}px))`;
    }
    pullDistance = 0;
});

// Function to set CSS safe area variables
function setCSSSafeAreaVariables() {
    const safeAreaTop = getComputedStyle(document.documentElement).getPropertyValue('--safe-area-inset-top').slice(0, -2);
    document.documentElement.style.setProperty('--sat', safeAreaTop + 'px');
}

// Call the function initially and add event listener for resize
setCSSSafeAreaVariables();
window.addEventListener('resize', setCSSSafeAreaVariables);