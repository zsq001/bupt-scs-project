<template>
  <mdui-list>
    <mdui-button @click="backtoChat">Back</mdui-button>
    <mdui-list-item>
      <mdui-text-field v-model="newMessage" clearable label="Enter Text"></mdui-text-field>
      <mdui-button slot="end-icon" @click="sendMessage">Send</mdui-button>
    </mdui-list-item>
    <mdui-list-subheader>Messages With {{ userName }}</mdui-list-subheader>
    <div class="chat-content">
      <mdui-list-item v-for="(message, index) in messages" :key="index">
        <div v-if="message.SendUserID === this.$attrs.nowUserId">{{ this.$attrs.nowUser }}:</div>
        <div v-else>{{ userName }}:</div>
        <div>{{ message.Content }}</div>
        <mdui-button v-if="isableRecall(message)" variant="elevated" slot="end-icon" @click="recall(message.ID)">Recall
        </mdui-button>
        <mdui-button variant="elevated" slot="end-icon" @click="deleteMessage(message.ID)">Delete
        </mdui-button>
      </mdui-list-item>
    </div>
  </mdui-list>
</template>

<script>
export default {
  data() {
    return {
      messages: [],
      newMessage: '',
      errorMessage: '',
      userName: '',
      myUserName: '',
      targetUid: '',
    };
  },
  async created() {
    if (typeof this.$route.params.id !== 'number') {
      await this.GetUserID(this.$route.params.id);
      this.$route.params.id = this.targetUid;
    }
    await this.fetchMessages();
    await this.GetUserName(this.$route.params.id);
    this.intervalID = setInterval(this.fetchMessages, 1000);
  },
  methods: {
    async recall(messageId) {
      try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/msg/v1/recall/${messageId}`, {
          method: 'GET',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        await this.fetchMessages();
        if (data && data.data) {
          this.messages = data.data;
        } else {
          this.errorMessage = data.error;
        }
      } catch (error) {
        this.errorMessage = error.message;
      }
    },
    isableRecall(message) {
      return parseInt(Date.now() / 1000) - message.SendTime < 300 && (message.ReceiveUserID === parseInt(this.$route.params.id));
    },
    backtoChat() {
      this.$router.push('/chat');
      clearInterval(this.intervalID);
    },
    async fetchMessages() {
      try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/msg/v1/get/${this.$route.params.id}`, {
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
          this.messages = data.data;
        } else {
          this.errorMessage = data.error;
        }
      } catch (error) {
        this.errorMessage = error.message;
      }
    },
    async sendMessage() {
      try {
        if (this.newMessage === '') {
          return;
        }
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/msg/v1/send/${this.$route.params.id}`, {
          method: 'POST',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token'),
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({Content: this.newMessage})
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (data && data.data) {
          this.messages.push(data.data);
          this.newMessage = '';
        } else {
          this.errorMessage = data.error;
        }
      } catch (error) {
        this.errorMessage = error.message;
      }
    },
    async deleteMessage(messageId) {
      try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/msg/v1/delete/${messageId}`, {
          method: 'GET',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        await this.fetchMessages();
        if (data && data.data) {
          this.messages = this.messages.filter(message => message.id !== messageId);

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
        await this.fetchMessages();
        if (data && data.data) {
          this.userName = data.data.name;

        } else {
          this.errorMessage = data.error;
        }
      } catch (error) {
        this.errorMessage = error.message;
      }
    },
    async GetUserID(name) {
      try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/user/info/${name}`, {
          method: 'GET',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        this.targetUid = data.data.id;
        if (data && data.data) {
          this.myUserName = data.data;

        } else {
          this.errorMessage = data.error;
        }
      } catch (error) {
        this.errorMessage = error.message;
      }
    }
  },
  beforeRouteLeave(to, from, next) {
    clearInterval(this.intervalID);
    next();
  }
};
</script>

<style scoped>
/* Add your styles here */
.chat-content {
  max-height: 60%; /* Adjust this value as needed */
  overflow-y: auto; /* Enable vertical scrolling */
}
</style>