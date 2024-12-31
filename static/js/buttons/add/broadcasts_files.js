// add broadcast files 
const bfadd_form = document.getElementById('add_broadcasts_files_dialog');

bfadd_form.addEventListener('submit', (event) => {
  event.preventDefault();

  // Get form data
  const broadcastId = document.getElementById('add_bf_broadcast_id');
  const videofileId = document.getElementById('add_bf_videofile_id');

  // Validate broadcast ID and videofile ID (positive integers)
  fetch('/api/add/broadcasts_files', { // Replace with your actual API endpoint
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      broadcast_id: broadcastId.value,
      videofile_id: videofileId.value,
    })
  })
    .then(response => response.text())
    .then(data => {
      // Handle successful response
      console.log('Broadcast added successfully:', data);
      const messageDiv = document.getElementById('bfadd-message');
      messageDiv.innerHTML = data; // Assuming the API returns a message
    })
    .catch(error => {
      const messageDiv = document.getElementById('bfadd-message');
      messageDiv.innerHTML = error; // Assuming the API returns a message
      console.error('Error adding broadcast:', error);
      // Handle errors here (e.g., display an error message)
    });
});
// add broadcast files end


