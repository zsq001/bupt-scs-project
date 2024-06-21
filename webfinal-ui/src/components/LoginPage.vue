<template>
  <div>
    <h2>Login</h2>
    <input v-model="username" type="text" placeholder="Username">
    <input v-model="password" type="password" placeholder="Password">
    <button @click="login">Login</button>
    <p class="error" v-if="errorMessage">{{ errorMessage }}</p>
  </div>
  <div>
    <h2>Register</h2>

    <router-link to="/register"><button>Register</button></router-link>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: '',
      password: '',
      errorMessage: ''
    }
  },
  methods: {
    async login() {
      const formData = new FormData();
      formData.append('user', this.username);
      formData.append('pass', this.password);

      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/user/login`, {
        method: 'POST',
        body: formData
      });

      const data = await response.json();

      if (data && data.data) {
        localStorage.setItem('token', data.data);
        location.reload(true);
        this.$router.push('/'); // Redirect to main page
      }
      else {
        this.errorMessage = "error:" + data.error;
      }
    }
  }
}
</script>