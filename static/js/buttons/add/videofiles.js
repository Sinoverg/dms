// add videofile
var vadd_form = document.getElementById('add_videofile_dialog');

vadd_form.addEventListener('submit', (event) => {
  event.preventDefault(); // Prevent default form submission

  // Get form data
  const filename = document.getElementById('add_v_videofilename').value.trim();
  const uploader = document.getElementById('add_v_uploader').value;
  const size = document.getElementById('add_v_size').value;
  const duration = document.getElementById('add_v_duration').value;

  // Validate filename (only letters and spaces)
  if (!/^[a-zA-Z\s]+$/.test(filename)) {
    document.getElementById('add_v_videofilename').setAttribute("aria-invalid", "true");
    return;
  } else {
    document.getElementById('add_v_videofilename').setAttribute("aria-invalid", "false");
  }

  // Validate uploader ID (positive integer)
  if (!/^\d+$/.test(uploader) || uploader <= 0) {
    document.getElementById('add_v_uploader').setAttribute("aria-invalid", "true");
    return;
  } else {
    document.getElementById('add_v_uploader').setAttribute("aria-invalid", "false");
  }

  // Validate file size (positive number)
  if (size <= 0) {
    document.getElementById('add_v_size').setAttribute("aria-invalid", "true");
    return;
  } else {
    document.getElementById('add_v_size').setAttribute("aria-invalid", "false");
  }

  // Validate duration (positive integer)
  if (!/^\d+$/.test(duration) || duration <= 0) {
    document.getElementById('add_v_duration').setAttribute("aria-invalid", "true");
    return;
  } else {
    document.getElementById('add_v_duration').setAttribute("aria-invalid", "false");
  }


  // Send data to server (replace with your actual logic)
  fetch('/api/add/videofiles', { // Replace with your actual API endpoint
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      filename,
      uploader,
      size,
      duration
    })
  })
    .then(response => response.text())
    .then(data => {
      const messageDiv = document.getElementById('vadd-message');
      messageDiv.innerHTML = data; // Assuming the API returns a message

      // Handle successful response
      // vadd_form.reset();
    })
    .catch(error => {
      const messageDiv = document.getElementById('vadd-message');
      messageDiv.innerHTML = 'Error adding videofile: ' + error;
      console.error('Error adding videofile:', error);
    });
});

// add videofile end