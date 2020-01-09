<template>
  <v-navigation-drawer v-model="drawer" :clipped="$vuetify.breakpoint.lgAndUp" app>
    <v-list dense flat>
      <template v-for="item in sideMenuItems">
        <v-list-item :key="item.text" @click="sideMenuGoTo(item)">
          <v-list-item-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>
              {{ item.text }}
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </template>
      <div style="padding-left: 10px;margin-top: 20px;">
        <span class="font-weight-light pa-2" v-text="`${bytesToSizeFormat(totalSize)}/${bytesToSizeFormat(maxDiskSize)}`"></span>
        <v-slider v-model="totalSize"
                  :color="color"
                  readonly
                  track-color="grey"
                  always-dirty
                  min="0"
                  :max="maxDiskSize">
        </v-slider>
      </div>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import api from '../../api'
import { bytesToSize } from '../utils/format';

export default {
  name: 'drawer',
  mounted() {
    this.serverStatic()
  },
  data() {
    return {
      // 侧边栏菜单
      sideMenuItems: [
        {
          icon: 'file_copy',
          text: '全部文件',
          router: 'disk',
        },
        {
          icon: 'share',
          text: '分享',
          router: 'share',
        },
        {
          icon: 'edit',
          text: '文档',
          type: 'doc',
        },
        {
          icon: 'image',
          text: '图片',
          type: 'image',
        },
        {
          icon: 'videocam',
          text: '视频',
          type: 'media',
        },
        {
          icon: 'library_music',
          text: '音乐',
          type: 'audio',
        },
        {
          icon: 'insert_link',
          text: '种子',
          type: 'torrent',
        },
        {
          icon: 'radio_button_checked',
          text: '其他',
          type: 'other',
        },
        {
          icon: 'delete_forever',
          text: '回收站',
          router: 'trashcan',
        },
        {
          icon: 'settings',
          text: '设置',
          router: 'setting',
        },
        {
          icon: 'storage',
          text: '存储空间',
          router: 'storage',
        },
      ],
      maxDiskSize: 100,
      totalSize: 0,
    }
  },
  methods: {
    ...mapActions({
      changeDrawer: 'theme/changeDrawer',
    }),
    // 获取服务器static大小
    serverStatic() {
      api.user.serverStatic().then((resData) => {
        if (resData.code === 0) {
          this.totalSize = resData.data.totalSize
          this.maxDiskSize = resData.data.maxDiskSize
        }
      })
    },
    // 菜单路由
    sideMenuGoTo(item) {
      if (item.router) {
        this.$router.push({
          name: item.router,
        })
      } else {
        this.$router.push({
          name: 'disk',
          query: { fileType: item.type },
        })
      }
    },
    bytesToSizeFormat(bytes) {
      return bytesToSize(bytes)
    },
  },
  computed: {
    ...mapGetters({
      drawerVisible: 'theme/drawerVisible',
    }),
    drawer: {
      get() {
        return this.drawerVisible
      },
      set() {
      },
    },
    // 文件磁盘配额色值
    color() {
      const { maxDiskSize, totalSize } = this
      if (totalSize < (maxDiskSize / 3)) return 'indigo'
      if (totalSize < (maxDiskSize / 3) * 2) return 'orange'
      if (totalSize < maxDiskSize) return 'red'
      return 'red'
    },
  },
}
</script>
