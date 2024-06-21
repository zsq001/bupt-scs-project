<template>
  <!--  <div>-->
  <!--    <div v-for="(chat, index) in chatList" :key="index">-->
  <!--      <p>{{ chat.name }}</p>-->
  <!--      <p>{{ chat.latest_msg }}</p>-->
  <!--      <p>{{ chat.latest_msg_time }}</p>-->
  <!--    </div>-->
  <!--  </div>-->


  <mdui-list-item>
    <mdui-text-field v-model="tuserid" label="Enter User Id"></mdui-text-field>
    <mdui-button slot="end-icon" @click="goToChat(tuserid)">Chat</mdui-button>
  </mdui-list-item>
  <mdui-list class="chat-list">
    <mdui-list-item v-for="(chat, index) in chatList" :key="index" :description="chat.latest_msg"
                    :headline="chat.name">
      <mdui-button variant="elevated" slot="end-icon" @click="goToChat(chat.ID)">Chat
      </mdui-button>
    </mdui-list-item>
  </mdui-list>
</template>

<script>
export default {
  props: {
    name: String
  },
  data() {
    return {
      chatList: [],
      errorMessage: '',
      userId: '',
      tuserid: ''
    };
  },
  methods: {
    async get_id() {
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
    async goToChat(id) {
      if (parseInt(id) === parseInt(this.userId)) {
        alert("You can't chat with yourself!");
      } else {
        try {
          const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/user/info/${id}`, {
            method: 'GET',
            headers: {
              'Authorization': 'Bearer ' + localStorage.getItem('token')
            }
          });

          if (!response.ok) {
            alert("User not found!")
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          this.$router.push(`/chat/${id}`);
        } catch (error) {
          this.errorMessage = error.message;
        }
      }
    },
    async fetchList() {
      try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/msg/v1/list`, {
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
          this.chatList = data.data;
        } else {
          this.errorMessage = data.error;
        }
      } catch (error) {
        this.errorMessage = error.message;
      }
    }
  },
  created() {
    this.fetchList();
    this.userId = this.get_id();
  },
};
</script>

<style scoped>
/* Add your styles here */
.chat-list {
  max-height: 60%; /* Adjust this value as needed */
  overflow-y: auto; /* Enable vertical scrolling */
}
</style>