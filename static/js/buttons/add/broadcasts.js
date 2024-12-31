// add broadcast 
const badd_form = document.getElementById('add_broadcast_dialog');

badd_form.addEventListener('submit', (event) => {
  event.preventDefault();

  // Get form data
  const b_start_time = document.getElementById('add_b_broadcast_start_time');
  const b_end_time = document.getElementById('add_b_broadcast_end_time');

  // Validate start time is before end time
  if (new Date(b_start_time.value) >= new Date(b_end_time.value)) {
    b_start_time.setAttribute("aria-invalid", "true");
    b_end_time.setAttribute("aria-invalid", "true");
    return;
  } else {
    b_start_time.setAttribute("aria-invalid", "false");
    b_end_time.setAttribute("aria-invalid", "false");
  }

    const formatISO = (dateString) => {
      const date = new Date(dateString);
      return date.toISOString().slice(0, 16); // Format as "YYYY-MM-DDTHH:MM"
  };
  // Send data to server (replace with your actual logic)
  fetch('/api/add/broadcasts', { // Replace with your actual API endpoint
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
  
  body: JSON.stringify({
      b_start_time: formatISO(b_start_time.value),
      b_end_time: formatISO(b_end_time.value)
  })

  })
    .then(response => response.text())
    .then(data => {
      // Handle successful response
      const messageDiv = document.getElementById('badd-message');
      messageDiv.innerHTML = data; // Assuming the API returns a message
      console.log('Broadcast added successfully:', data);
    })
    .catch(error => {
      const messageDiv = document.getElementById('badd-message');
      messageDiv.innerHTML = error; // Assuming the API returns a message
      console.error('Error adding broadcast:', error);
      // Handle errors here (e.g., display an error message)
    });
});
// add broadcast end