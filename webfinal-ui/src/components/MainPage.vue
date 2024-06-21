<template>
  <div>
    <h1>Welcome to the Main Page! {{this.$attrs.nowUser}}</h1>
    Want to Update Info?
    <mdui-text-field v-model="newusername" label="New User Name"></mdui-text-field>
    <mdui-text-field v-model="newuserpass" type="password" toggle-password label="New Password"></mdui-text-field>
    <mdui-button slot="end-icon" @click="updateUserInfo(this.$attrs.nowUserId, newusername, newuserpass)">Update</mdui-button>
    <div>Want to Logout?
    <mdui-button v-if="check_token" @click="logout">Logout</mdui-button></div>
  </div>
</template>

<script>
export default {
  computed: {
    token() {
      return localStorage.getItem('token');
    },
    check_token() {
      console.log('Token:', localStorage.getItem('token'));
      return localStorage.getItem('token') != null;

    },
  },
  methods: {
    logout() {
      // Remove the token from local storage
      localStorage.removeItem('token');

      // Redirect to login page
      this.$router.push('/login');
    },
    async updateUserInfo(userId, username, password) {
      console.log(userId, username, password)
      if (!username || !password) {
        alert('Please enter both username and password!')
        return;
      }
      try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/user/update/${userId}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('token')
          },
          body: JSON.stringify({ name:username, password })
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (data && data.data) {
          // The user info was successfully updated
          // You can update the local user info here if needed
          location.reload(true);

          alert('User info updated successfully!');
          this.changeUserName(username);


        } else {
          // The API returned an error
          throw new Error(data.error);
        }
      } catch (error) {
        // Handle the error
        console.error('An error occurred:', error.message);
      }
    }
  }
}
</script>