// add broadcast files 
const buadd_form = document.getElementById('buadd-form');

buadd_form.addEventListener('submit', (event) => {
  event.preventDefault();

  // Get form data
  const broadcastId = document.getElementById('add_bu_broadcast_id');
  const videofileId = document.getElementById('add_bu_user_id');

  // Validate broadcast ID and videofile ID (positive integers)
  fetch('/api/add/broadcasts_users', { // Replace with your actual API endpoint
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      broadcast_id: broadcastId.value,
      user_id: videofileId.value,
    })
  })
    .then(response => response.text())
    .then(data => {
      // Handle successful response
      console.log('Broadcast added successfully:', data);
      const messageDiv = document.getElementById('buadd-message');
      messageDiv.innerHTML = data; // Assuming the API returns a message
    })
    .catch(error => {
      const messageDiv = document.getElementById('buadd-message');
      messageDiv.innerHTML = error; // Assuming the API returns a message
      console.error('Error adding broadcast:', error);
      // Handle errors here (e.g., display an error message)
    });
});
// add broadcast files end


