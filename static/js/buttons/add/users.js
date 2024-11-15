// add user
const uadd_form = document.getElementById('uadd-form');

uadd_form.addEventListener('submit', (event) => {
  event.preventDefault(); // Prevent default form submission

  // Get form data
  const first_name = document.getElementById('add_u_first_name').value.trim();
  const second_name = document.getElementById('add_u_second_name').value.trim();
  const nickname = document.getElementById('add_u_nickname').value.trim();
  const login = document.getElementById('add_u_login').value.trim();
  const email = document.getElementById('add_u_email').value.trim();
  const password = document.getElementById('add_u_password').value;
  const repassword = document.getElementById('add_u_repassword').value;

  // Simple validation (can be extended with more checks)
  let isValid = true;
  const errorMessages = [];

  if (!/^[a-zA-Z]+$/.test(first_name)) {
    isValid = false;
    errorMessages.push('First name can only contain letters');
    document.getElementById('add_u_first_name').setAttribute("aria-invalid", "true");
  } else {

    document.getElementById('add_u_first_name').setAttribute("aria-invalid", "false");
  }

  if (!/^[a-zA-Z]+$/.test(second_name)) {
    isValid = false;
    errorMessages.push('Second name can only contain letters');
    document.getElementById('add_u_second_name').setAttribute("aria-invalid", "true");
  }
  else {

    document.getElementById('add_u_second_name').setAttribute("aria-invalid", "false");
  }


  if (nickname.length > 15 || nickname.length < 4) {
    isValid = false;
    errorMessages.push('Nickname can not be longer than 15 characters');
    document.getElementById('add_u_nickname').setAttribute("aria-invalid", "true");
  } else {
    document.getElementById('add_u_nickname').setAttribute("aria-invalid", "false");
  }

  if (/\s/.test(login) || !/^[a-zA-Z0-9]+$/.test(login)) {
    isValid = false;
    errorMessages.push('Login can not contain spaces or special characters');
    document.getElementById('add_u_login').setAttribute("aria-invalid", "true");
  }
  else {

    document.getElementById('add_u_login').setAttribute("aria-invalid", "false");
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    isValid = false;
    errorMessages.push('Invalid email format');
    document.getElementById('add_u_email').setAttribute("aria-invalid", "true");
  } else {

    document.getElementById('add_u_email').setAttribute("aria-invalid", "false");
  }

  if (password.length < 8) {
    isValid = false;
    errorMessages.push('Password must be at least 8 characters long');
    document.getElementById('add_u_password').setAttribute("aria-invalid", "true");
  } else {
    document.getElementById('add_u_password').setAttribute("aria-invalid", "false");
  }
  if (repassword != password || repassword.length < 8) {
    isValid = false;
    errorMessages.push('Password must be at least 8 characters long');
    document.getElementById('add_u_repassword').setAttribute("aria-invalid", "true");
  } else {
    document.getElementById('add_u_repassword').setAttribute("aria-invalid", "false");
  }
  // // If validation fails, display error messages
  if (!isValid) {
    return;
  }

  // Send data to server (replace with your actual logic)
  fetch('/api/add/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      first_name,
      second_name,
      nickname,
      login,
      email,
      password
    })
  })
    .then(response => response.text())
    .then(data => {
      var errorMessageContainer = document.getElementById('add-message');
      errorMessageContainer.innerHTML = data;
      // Handle successful response
      console.log('User added successfully:', data);
      // You can display a success message here
      // document.getElementById('add_user_dialog').close();
      // Clear form fields (optional)
      // uadd_form.reset();
    })
    .catch(error => {
      var errorMessageContainer = document.getElementById('add-message');
      errorMessageContainer.innerHTML = error;
      console.error('Error adding user:', error);
      // Handle server errors here (e.g., display an error message)
    });
});
// add user end