<template>
  <mdui-layout>
    <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
    <link href="https://fonts.googleapis.com/icon?family=Material+Icons+Outlined" rel="stylesheet">
    <div style="position: relative" class="main-body">
      <mdui-navigation-rail contained class="side-bar">
        <mdui-navigation-rail-item v-if="isAuthenticated" icon="insert_photo--outlined" value="Picture"
                                   @click="picture">
          Picture
        </mdui-navigation-rail-item>
        <mdui-navigation-rail-item v-if="isAuthenticated" icon="message--outlined" value="Chat" @click="chat">Chat
        </mdui-navigation-rail-item>
        <mdui-navigation-rail-item icon="assignment_ind" value="Account" @click="account">Account
        </mdui-navigation-rail-item>
        <mdui-navigation-rail-item icon="fact_check" value="Admin" @click="admin">Admin</mdui-navigation-rail-item>
      </mdui-navigation-rail>
      <div class="page-view">
        <router-view :nowUser="userName" :nowUserId="userId" class="page-view"></router-view>
      </div>
    </div>
  </mdui-layout>
</template>

<style scoped>
.main-body {
  display: flex;
  //align-items: center;
  justify-content: left;
  //margin-left: 10%;
  padding: 0;
  position: fixed;
  height: 100vh; /* Ensure full viewport height */
  margin: 0; /* Remove default body margin */
}

.page-view {
  //display: flow;
  //align-items: center;
  //justify-content: left;
  position: fixed;
  width: 80%;
  margin: 0 auto;
  height: 100vh; /* Ensure full viewport height */
  //margin: 0; /* Remove default body margin */
  top:0;
}

.side-bar {
  height: 100%;
  width: 80px;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1;
}
</style>

<script>
export default {
  data() {
    return {
      errorMessage: '',
      userId: '',
      userName: ''
    };
  },
  methods: {
    isAuthenticated() {
      return localStorage.getItem('token') != null;
    },
    changeUserName(newName) {
      this.userName = newName;
    },
    admin(){
      clearInterval(this.intervalID);
      return this.$router.push('/admin');
    },
    picture() {
      clearInterval(this.intervalID);
      return this.$router.push('/pic');
    },
    chat() {
      clearInterval(this.intervalID);
      return this.$router.push('/chat');
    },
    account() {
      clearInterval(this.intervalID);
      return this.$router.push('/');
    },
    async getUserId() {
      try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/user/whoami`, {
          method: 'GET',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (data && data.data) {
          this.userId = data.data;
        } else {
          this.errorMessage = data.error;
        }
      } catch (error) {
        this.errorMessage = error.message;
      }
    },
    async GetUserName(userId) {
      try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/user/info/${userId}`, {
          method: 'GET',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        if (data && data.data) {
          this.userName = data.data.name;

        } else {
          this.errorMessage = data.error;
        }
      } catch (error) {
        this.errorMessage = error.message;
      }
    }
  },
  async created() {
    await this.getUserId();
    await this.GetUserName(this.userId);
  }
}
</script>