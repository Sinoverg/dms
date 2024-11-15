// add broadcast files 
const bfadd_form = document.getElementById('add_broadcasts_files_dialog');

bfadd_form.addEventListener('submit', (event) => {
  event.preventDefault();

  // Get form data
  const broadcastId = document.getElementById('add_bf_broadcast_id');
  const videofileId = document.getElementById('add_bf_videofile_id');

  // Validate broadcast ID and videofile ID (positive integers)
  if (!/^\d+$/.test(broadcastId.value) || broadcastId.value <= 0) {
    broadcastId.setAttribute("aria-invalid", "true");
    return;
  } else {
    videofileId.setAttribute("aria-invalid", "false");
  }

  if (!/^\d+$/.test(videofileId) || videofileId <= 0) {
    broadcastId.setAttribute("aria-invalid", "true");
    return;
  } else {
    videofileId.setAttribute("aria-invalid", "false");
  }

  fetch('/api/add/broadcasts_files', { // Replace with your actual API endpoint
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      broadcast_start_time: broadcastId.value,
      broadcast_end_time: videofileId.value,
    })
  })
    .then(response => response.text())
    .then(data => {
      // Handle successful response
      console.log('Broadcast added successfully:', data);
      const messageDiv = document.getElementById('badd-message');
      messageDiv.textContent = data.message; // Assuming the API returns a message
      badd_form.reset();
    })
    .catch(error => {
      const messageDiv = document.getElementById('badd-message');
      messageDiv.textContent = error.message; // Assuming the API returns a message
      console.error('Error adding broadcast:', error);
      // Handle errors here (e.g., display an error message)
    });
});
// add broadcast files end


