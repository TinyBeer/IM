<template>
  <div>
    <b-navbar
      toggleable="lg"
      type="dark"
      variant="info"
    >
      <b-container>
        <b-navbar-brand @click="$router.replace({name:'Home'})">TinyBeer</b-navbar-brand>

        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

        <b-collapse
          id="nav-collapse"
          is-nav
        >

          <!-- Right aligned nav items -->
          <b-navbar-nav class="ml-auto">
            <b-nav-item-dropdown
              left
              v-if="userInfo"
            >
              <template #button-content>
                <em> {{ userInfo.name }} </em>
              </template>
              <b-dropdown-item href="#">个人主页</b-dropdown-item>
              <b-dropdown-item @click="logout">登出</b-dropdown-item>
            </b-nav-item-dropdown>

          </b-navbar-nav>
          <b-navbar-nav v-if="!userInfo">
            <b-nav-item
              v-if="$route.name != 'login'"
              @click="$router.replace({name:'login'})"
            >登录</b-nav-item>
            <b-nav-item
              v-if="$route.name != 'register'"
              @click="$router.replace({name:'register'})"
            >注册</b-nav-item>

          </b-navbar-nav>
        </b-collapse>
      </b-container>
    </b-navbar>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';

export default {
  computed: mapState({
    userInfo: (state) => state.userModule.userInfo,
  }),
  methods: {
    ...mapActions('userModule', { userLogout: 'logout' }),
    logout() {
      this.userLogout();
      this.$router.replace({ name: 'login' });
    },
  },
};
</script>

<style lang="scss" scoped >
</style>
