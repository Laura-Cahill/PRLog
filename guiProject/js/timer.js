// Switch between Countdown and Stopwatch UI
const countdownTab = document.getElementById("countdownTab");
const stopwatchTab = document.getElementById("stopwatchTab");
const countdownUI = document.getElementById("countdownUI");
const stopwatchUI = document.getElementById("stopwatchUI");

countdownTab.addEventListener("click", () => {
  countdownUI.classList.remove("hidden");
  stopwatchUI.classList.add("hidden");
  countdownTab.classList.add("bg-purple-400", "text-white");
  stopwatchTab.classList.remove("bg-purple-400", "text-white");
  stopwatchTab.classList.add("text-purple-400");
});

stopwatchTab.addEventListener("click", () => {
  stopwatchUI.classList.remove("hidden");
  countdownUI.classList.add("hidden");
  stopwatchTab.classList.add("bg-purple-400", "text-white");
  countdownTab.classList.remove("bg-purple-400", "text-white");
  countdownTab.classList.add("text-purple-400");
});

// Countdown Timer
let countdownInterval;
const countdownDisplay = document.getElementById("countdownDisplay");
const countdownMinutes = document.getElementById("countdownMinutes");
const countdownSeconds = document.getElementById("countdownSeconds");

document.getElementById("startCountdown").addEventListener("click", () => {
  let totalSeconds = parseInt(countdownMinutes.value) * 60 + parseInt(countdownSeconds.value);
  if (isNaN(totalSeconds) || totalSeconds <= 0) return;

  clearInterval(countdownInterval);

  countdownInterval = setInterval(() => {
    if (totalSeconds <= 0) {
      clearInterval(countdownInterval);
      return;
    }

    totalSeconds--;
    const min = String(Math.floor(totalSeconds / 60)).padStart(2, "0");
    const sec = String(totalSeconds % 60).padStart(2, "0");
    countdownDisplay.textContent = `${min}:${sec}`;
  }, 1000);
});

document.getElementById("pauseCountdown").addEventListener("click", () => {
  clearInterval(countdownInterval);
});

document.getElementById("resetCountdown").addEventListener("click", () => {
  clearInterval(countdownInterval);
  countdownDisplay.textContent = "00:00";
  countdownMinutes.value = "";
  countdownSeconds.value = "";
});

// Stopwatch
let stopwatchInterval;
let stopwatchTime = 0;
const stopwatchDisplay = document.getElementById("stopwatchDisplay");

function updateStopwatchDisplay() {
  const min = String(Math.floor(stopwatchTime / 60)).padStart(2, "0");
  const sec = String(stopwatchTime % 60).padStart(2, "0");
  stopwatchDisplay.textContent = `${min}:${sec}`;
}

document.getElementById("startStopwatch").addEventListener("click", () => {
  clearInterval(stopwatchInterval);
  stopwatchInterval = setInterval(() => {
    stopwatchTime++;
    updateStopwatchDisplay();
  }, 1000);
});

document.getElementById("pauseStopwatch").addEventListener("click", () => {
  clearInterval(stopwatchInterval);
});

document.getElementById("resetStopwatch").addEventListener("click", () => {
  clearInterval(stopwatchInterval);
  stopwatchTime = 0;
  updateStopwatchDisplay();
});
