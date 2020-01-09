<template>
  <v-app-bar :clipped-left="$vuetify.breakpoint.lgAndUp" app dark color="blue darken-3">
    <v-toolbar-title style="width: 300px" class="ml-0 pl-3">
      <v-app-bar-nav-icon @click.stop="changeDrawer"></v-app-bar-nav-icon>
      <span class="hidden-sm-and-down">Go FileUpload</span>
    </v-toolbar-title>
    <v-text-field v-model="q"
                  flat solo-inverted hide-details
                  prepend-inner-icon="search"
                  label="Search"
                  @keyup.enter="handleSearch"
                  class="hidden-sm-and-down"/>
    <v-spacer/>
    <div class="v-btn v-btn--flat v-btn--icon v-btn--round theme--dark v-size--default">
      <a href="https://github.com/Beats0/gofileupload" target="_blank">
        <svg viewBox="64 64 896 896" fill="currentColor" style="color: #fff; width: 26px;"><path d="M511.6 76.3C264.3 76.2 64 276.4 64 523.5 64 718.9 189.3 885 363.8 946c23.5 5.9 19.9-10.8 19.9-22.2v-77.5c-135.7 15.9-141.2-73.9-150.3-88.9C215 726 171.5 718 184.5 703c30.9-15.9 62.4 4 98.9 57.9 26.4 39.1 77.9 32.5 104 26 5.7-23.5 17.9-44.5 34.7-60.8-140.6-25.2-199.2-111-199.2-213 0-49.5 16.3-95 48.3-131.7-20.4-60.5 1.9-112.3 4.9-120 58.1-5.2 118.5 41.6 123.2 45.3 33-8.9 70.7-13.6 112.9-13.6 42.4 0 80.2 4.9 113.5 13.9 11.3-8.6 67.3-48.8 121.3-43.9 2.9 7.7 24.7 58.3 5.5 118 32.4 36.8 48.9 82.7 48.9 132.3 0 102.2-59 188.1-200 212.9a127.5 127.5 0 0 1 38.1 91v112.5c.8 9 0 17.9 15 17.9 177.1-59.7 304.6-227 304.6-424.1 0-247.2-200.4-447.3-447.5-447.3z"></path></svg>
      </a>
    </div>
    <v-btn icon>
      <v-icon>notifications</v-icon>
    </v-btn>
    <div class="text-center">
      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <v-btn icon large v-on="on">
            <v-avatar size="32px" item>
              <v-img src="https://avatar-static.segmentfault.com/236/630/2366300268-5a9e50c785bce_huge256" alt="Vuetify"/>
            </v-avatar>
          </v-btn>
        </template>
        <v-list>
          <v-list-item @click="{}">
            <v-list-item-title> 添加其他账号 </v-list-item-title>
          </v-list-item>
          <v-list-item @click="logOut">
            <v-list-item-title> 退出 </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </div>
  </v-app-bar>
</template>

<script>
import { mapActions } from 'vuex'
import storage from '../utils/storage'

export default {
  name: 'appBar',
  data() {
    return {
      q: '',
    }
  },
  methods: {
    ...mapActions({
      changeDrawer: 'theme/changeDrawer',
    }),
    handleSearch() {
      const q = this.q
      this.$router.push({
        name: 'disk',
        query: { q },
      })
    },
    // 退出登录
    logOut() {
      storage.delete('gin_token')
      setTimeout(() => {
        window.location.reload()
      }, 300)
    },
  },
}
</script>
