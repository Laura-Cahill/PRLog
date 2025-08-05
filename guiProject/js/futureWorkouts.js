document.addEventListener("DOMContentLoaded", () => {
  // Highlight selected calendar date
  const calendarCells = document.querySelectorAll(".calendar-cell");
  calendarCells.forEach(cell => {
    cell.addEventListener("click", () => {
      calendarCells.forEach(c => c.classList.remove("bg-purple-300"));
      cell.classList.add("bg-purple-300");
      document.getElementById("selectedDate").textContent = cell.dataset.date || cell.textContent;
    });
  });

  // Handle save button click
  const form = document.getElementById("futureWorkoutsForm");
  const message = document.getElementById("confirmationMessage");

  form.addEventListener("submit", (event) => {
    event.preventDefault();

    // Grab input values
    const date = document.getElementById("workoutDate").value;
    const workout = document.getElementById("exercise").value;
    const sets = document.getElementById("sets").value;
    const reps = document.getElementById("reps").value;
    const weight = document.getElementById("weight").value;

    // For now: just log to console
    console.log({ date, workout, sets, reps, weight });

    // Show confirmation
    message.classList.remove("hidden");
    setTimeout(() => message.classList.add("hidden"), 3000);
  });
});
