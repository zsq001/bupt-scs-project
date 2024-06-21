<script setup>
import { ref, onMounted } from 'vue';

const users = ref([]);
async function fetchUsers() {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/user/list`,{
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    }
  });
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  const data = await response.json();
  if (data && data.data){
    users.value = data.data;
  }
}

async function deleteUser(id) {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/user/delete/${id}`,{
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    }
  });
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  // Remove the user from the list
  users.value = users.value.filter(user => user.ID !== id);
}

function permission(id){
  if (id==1){
    return "Admin";
  }
  else return "User";
}
onMounted(fetchUsers);
</script>

<template>
  <mdui-navigation-bar class="bottom-bar">
    <mdui-navigation-bar-item icon="photo">Pictures</mdui-navigation-bar-item>
    <mdui-navigation-bar-item icon="people" value="item-3">Users</mdui-navigation-bar-item>
  </mdui-navigation-bar>
  <mdui-list class="user-list">
    <mdui-list-item v-for="(user, index) in users" :key="index" :description="'UserID: '+user.ID +' ('+permission(user.Role)+')'"
                    :headline="user.Name">
      <mdui-button variant="elevated" slot="end-icon" @click="deleteUser(user.ID)">Delete
      </mdui-button>
    </mdui-list-item>
  </mdui-list>
</template>

<style scoped>
.bottom-bar {
  position: fixed;
  height: 81px;
}
.user-list{
  top:0;
}

</style>