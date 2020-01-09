<template>
  <v-app id="inspire">
    <Drawer />
    <AppBar />
    <v-content>
      <v-container fluid fill-height style="padding-top: 0">
        <div class="full-width">
          <div class="tool-header-box">
            <!-- 文件列表操作栏  -->
            <div class="tool-header-actions">
              <v-btn color="primary" outlined small style="margin-right: 10px" @click="handleCancel"> 取消分享 <v-icon right dark>cloud_download</v-icon></v-btn>
            </div>
            <div class="tool-header-viewer">
              <v-btn small fab :elevation="0"><v-icon>storage</v-icon></v-btn>
            </div>
          </div>
          <div class="table-container" v-viewer>
            <v-data-table
              v-model="selected"
              :headers="tableHeaders"
              fixed-header
              :footer-props="{ 'items-per-page-options': [100] }"
              :loading="loading"
              :server-items-length="total"
              :items="lists"
              show-select
              @update:sort-desc="sortDescFn"
              @update:sort-by="sortByFn"
              @pagination="paginationFn"
              class="elevation-1 full-width">
              <template v-slot:item.image="{ item }">
                {{ item.file_type }}
                <!-- 预览图  -->
                <!-- <FileElem :file="item" @vewLarger="showThumbnail" />-->
              </template>
              <template v-slot:item.file_name="{ item }">
                <span v-if="item.is_dir">
                  {{ item.file_name }}
                </span>
                <span v-else>{{ item.file_name }}</span>
              </template>
              <template v-slot:item.sdate="{ item }">
                {{ item.sdate * 1000 | dayjsFormat }}
              </template>
              <template v-slot:item.action="{ item }">
                <v-icon class="mr-5" title="复制链接" color="blue lighten-1" @click="copyItem(item)">
                  link
                </v-icon>
                <v-icon size="18" class="mr-2" title="删除" color="deep-orange lighten-1" @click="deleteItem(item)">
                  block
                </v-icon>
              </template>
            </v-data-table>
            <v-spacer style="height: 30px;background: #fff;" />
          </div>
        </div>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import api from '../../api'
import Drawer from '../components/drawer.vue'
import AppBar from '../components/appBar.vue'

export default {
  name: 'share',
  components: {
    AppBar,
    Drawer,
  },
  data() {
    return {
      selected: [],
      tableHeaders: [
        {
          text: '#',
          value: 'image',
          sortable: false,
          width: '50px',
        },
        {
          text: '文件名',
          align: 'left',
          value: 'file_name',
          width: '250px',
        },
        {
          text: '浏览次数',
          value: 'view',
        },
        {
          text: '保存次数',
          value: 'save',
        },
        {
          text: '下载次数',
          value: 'download',
        },
        {
          text: '分享时间',
          value: 'sdate',
        },
        {
          text: '操作',
          value: 'action',
          sortable: false,
        },
      ],
      loading: true,
      // 分页
      page: 1,
      pages: 1,
      limit: 100,
      total: 100,
      order: 'date',
      desc: 1,
      fileType: '',
      lists: [],
    }
  },
  mounted() {
    this.refresh()
  },
  methods: {
    refresh() {
      this.loading = true
      this.lists = [
        {
          id: 2,
          file_name: '新建文本文档.txt',
          file_size: 716,
          file_ext: '.txt',
          file_type: 'doc',
          md5: 'e2753424c937bc4023008d87f3e989c1',
          is_dir: 0,
          view: 100,
          save: 100,
          download: 100,
          sdate: 1577256037,
        },
      ]
      this.loading = false
    },
    async goToItem(item) {
      console.log(item)
      const data = {
        fs_id: item.id,
      }
      const resData = await api.user.findFilePath(data)
      if (resData.code === 0) {
        // 在新窗口打开目录
        const { href } = this.$router.resolve({
          name: 'disk',
          query: {
            dirId: resData.data.dirId,
          },
        })
        window.open(href, '_blank')
      }
    },
    handleCancel() {
    },
    copyItem(item) {
      console.log(item)
    },
    deleteItem(item) {
      console.log(item)
    },
    sortDescFn(desc) {
      this.desc = desc[0] === true ? 1 : 0
      this.refresh()
    },
    sortByFn(order) {
      if (order[0]) {
        this.order = order[0]
        this.refresh()
      }
    },
    paginationFn({ page = 1, itemsPerPage = 50 }) {
      // 防止page, limit改变导致重新请求刷新
      if (this.page !== page || this.limit !== itemsPerPage) {
        this.page = page
        this.limit = itemsPerPage
        this.refresh()
      }
    },
  },
}
</script>
