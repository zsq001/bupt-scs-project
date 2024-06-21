<template>
  <AdminPicturePage v-if="currentComponent === 'AdminPicturePage'" />
  <AdminUserPage v-if="currentComponent === 'AdminUserPage'" />

  <mdui-navigation-bar class="bottom-bar">
    <mdui-navigation-bar-item icon="photo" @click="currentComponent = 'AdminPicturePage'">Pictures</mdui-navigation-bar-item>
    <mdui-navigation-bar-item icon="people" value="item-3" @click="currentComponent = 'AdminUserPage'">Users</mdui-navigation-bar-item>
  </mdui-navigation-bar>
</template>

<script>
import AdminPicturePage from "@/components/AdminPicturePage.vue";
import AdminUserPage from "@/components/AdminUserPage.vue";

export default {
  components: {
    AdminUserPage,
    AdminPicturePage
  },
  data() {
    return {
      currentComponent: 'AdminPicturePage'
    }
  }
}
</script>
<style scoped>
.bottom-bar {
  position: fixed;
  height: 81px;
}
</style>