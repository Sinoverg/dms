function Delete(){
var selectedRadio = document.querySelector('#table_set input[name="table"]:checked');
var table_name = selectedRadio ? selectedRadio.nextSibling.textContent.trim() : null;
const id_field = document.getElementById('delete_id');
fetch('/api/delete/' + table_name + '/' + id_field.value, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json'
    },
  })
    .then(response => response.text())
    .then(data => {
      var errorMessageContainer = document.getElementById('delete-message');
      errorMessageContainer.innerHTML = data;
      // Handle successful response
      console.log('User added successfully:', data);
      // You can display a success message here
      // document.getElementById('add_user_dialog').close();
      // Clear form fields (optional)
      // uadd_form.reset();
    })
    .catch(error => {
      var errorMessageContainer = document.getElementById('delete-message');
      errorMessageContainer.innerHTML = error;
      console.error('Error adding user:', error);
      // Handle server errors here (e.g., display an error message)
    });
}

