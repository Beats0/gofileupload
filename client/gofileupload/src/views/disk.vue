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
              <v-menu offset-y open-on-hover v-show="diskMode">
                <template v-slot:activator="{ on }">
                  <v-btn color="primary"
                         v-on="on"
                         v-show="diskMode"
                         outlined small
                         @click="$refs.fileInput.click()"
                         style="margin-right: 10px"> 上传 <v-icon right dark>cloud_upload</v-icon></v-btn>
                </template>
                <v-list>
                  <v-list-item @click="$refs.fileInput.click()">
                    <v-list-item-title>上传文件</v-list-item-title>
                  </v-list-item>
                  <v-list-item @click="$refs.folderInput.click()">
                    <v-list-item-title>上传文件夹</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
              <v-btn color="primary" outlined small style="margin-right: 10px" @click="showCreateFolderDialog = true" v-show="diskMode"> 新建文件夹 <v-icon right dark>create_new_folder</v-icon></v-btn>
              <v-btn color="primary" outlined small style="margin-right: 10px" @click="handleDownload"> 下载 <v-icon right dark>cloud_download</v-icon></v-btn>
              <v-btn color="primary" outlined small style="margin-right: 10px" @click="handleDelete"> 删除 <v-icon right dark>delete_forever</v-icon></v-btn>
              <v-btn color="primary" outlined small style="margin-right: 10px" @click="handleShare"> 分享 <v-icon right dark>share</v-icon></v-btn>
              <v-btn color="primary" outlined small style="margin-right: 10px" @click="handleCopy"> 复制到 <v-icon right dark>file_copy</v-icon></v-btn>
              <v-btn color="primary" outlined small style="margin-right: 10px" @click="handleMove"> 移动到 <v-icon right dark>redo</v-icon></v-btn>
            </div>
            <div class="tool-header-viewer">
              <v-btn small fab :elevation="0"><v-icon>storage</v-icon></v-btn>
            </div>
          </div>
          <div class="table-container">
            <v-data-table
              v-model="selected"
              :headers="tableHeaders"
              fixed-header
              :footer-props="{ 'items-per-page-options': [50, 100] }"
              :loading="loading"
              :server-items-length="total"
              :items="lists"
              show-select
              @update:sort-desc="sortDescFn"
              @update:sort-by="sortByFn"
              @pagination="paginationFn"
              class="elevation-1 full-width">
              <template v-slot:item.image="{ item }">
              <!-- 预览图  -->
               <FileElem :file="item" @vewLarger="showThumbnail" />
              </template>
              <template v-slot:item.file_name="{ item }">
                <span v-if="item.is_dir">
                  <router-link :to="{name: 'disk', query: {dirId: item.id ,order }}"> {{ item.file_name }} </router-link>
                </span>
                <span v-else>{{ item.file_name }}</span>
              </template>
              <template v-slot:item.file_size="{ item }">
                {{  bytesToSizeFormat(item.file_size) }}
              </template>
              <template v-slot:item.date="{ item }">
                {{ item.last_modified * 1000 | dayjsFormat }}
              </template>
              <template v-slot:item.action="{ item }">
                <v-icon small class="mr-2" @click="downloadItem(item)" >
                  vertical_align_bottom
                </v-icon>
                <v-menu offset-y>
                  <template v-slot:activator="{ on }">
                    <v-icon small v-on="on">
                      more_horiz
                    </v-icon>
                  </template>
                  <v-list class="file_action_list">
                    <v-list-item @click="downloadItem(item)">
                      <v-list-item-title>下载</v-list-item-title>
                    </v-list-item>
                    <v-list-item @click="shareItem(item)">
                      <v-list-item-title>分享</v-list-item-title>
                    </v-list-item>
                    <v-list-item @click="copyItem(item)">
                      <v-list-item-title>复制到</v-list-item-title>
                    </v-list-item>
                    <v-list-item @click="moveItem(item)">
                      <v-list-item-title>移动到</v-list-item-title>
                    </v-list-item>
                    <v-list-item @click="renameItem(item)">
                      <v-list-item-title>重命名</v-list-item-title>
                    </v-list-item>
                    <v-list-item @click="goToItem(item)">
                      <v-list-item-title>查看目录</v-list-item-title>
                    </v-list-item>
                    <v-list-item @click="deleteItem(item)">
                      <v-list-item-title>删除</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </template>
            </v-data-table>
            <v-spacer style="height: 30px;background: #fff;" />
          </div>
        </div>
      </v-container>
    </v-content>
    <input type="file" accept="*" multiple ref="fileInput" style="display: none;" @change="uploadFile($event)" />
    <input type="file" accept="*" webkitdirectory directory multiple ref="folderInput" style="display: none;" @change="uploadFolder($event)" />
    <div class="upload-progress-container" v-show="showUploadPanel">
      <v-expansion-panels>
        <v-expansion-panel>
          <v-expansion-panel-header>
            {{ panelScanProgress > 0 ? `正在扫描 ${panelScanProgress}%` : `` }}
            {{ panelUploadProgress > 0 ? `正在上传 ${panelUploadProgress}%` : `` }}
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <div class="v-expansion-action-container">
              <div>当前{{ panelUploadProgress }}%</div>
              <div>
                <v-btn text color="primary" small @click="handleCancel">取消</v-btn>
                <v-btn text color="primary" small @click="showUploadPanel=false">隐藏</v-btn>
              </div>
            </div>
            <div class="upload-list-container">
              <div class="upload-item" v-show="panelScanProgress > 0 || panelUploadProgress > 0">
                <span v-show="currentFile.name">{{ currentFile.name }}</span>
                <v-progress-circular
                  :rotate="-90"
                  :size="20"
                  :width="3"
                  :value="panelScanProgress"
                  color="primary"/>
              </div>
            </div>
          </v-expansion-panel-content>
        </v-expansion-panel>
      </v-expansion-panels>
    </div>
    <!-- 相册预览  -->
    <img-viewer ref="viewer" />
    <!-- dialog-新建文件夹 -->
    <v-dialog v-model="showCreateFolderDialog" max-width="400">
      <v-card>
        <div class="newFolderContainer">
          <span class="newFolderTitle">新文件夹</span>
          <v-text-field v-model="createFolderName" label="新建文件夹" single-line @keyup.enter="createFolderSubmit"></v-text-field>
        </div>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="showCreateFolderDialog = false">取消</v-btn>
          <v-btn text @click="createFolderSubmit">确认</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!-- dialog-重命名 -->
    <v-dialog v-model="showNewFileNameDialog" max-width="400">
      <v-card>
        <div class="newFolderContainer">
          <span class="newFolderTitle">重命名</span>
          <v-text-field v-model="newFileName" single-line @keyup.enter="renameItemSubmit"></v-text-field>
        </div>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="showNewFileNameDialog = false">取消</v-btn>
          <v-btn text @click="renameItemSubmit">确认</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script>
import axios from 'axios'
import SparkMD5 from 'spark-md5'
import { mapActions } from 'vuex'
import FileElem from '../components/fileElem.vue'
import Drawer from '../components/drawer.vue'
import AppBar from '../components/appBar.vue'
import storage from '../utils/storage'
import api from '../../api'
import { backEnd } from '../../config';
import { bytesToSize } from '../utils/format';
import ImgViewer from '../components/imgViewer.vue';
import { saveToDisk } from '../utils/util';
// import { makeCancelable } from '../utils/util';

export default {
  name: 'disk',
  components: {
    FileElem,
    AppBar,
    Drawer,
    ImgViewer,
  },
  data: () => ({
    dirId: 0,
    backEnd: '',
    // 选中文件列表数据
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
        width: '50%',
      },
      {
        text: '大小',
        value: 'file_size',
      },
      {
        text: '修改日期',
        value: 'date',
      },
      {
        text: '操作',
        value: 'action',
        sortable: false,
      },
    ],
    // 上传分片
    // 当前文件, 为file对象
    currentFile: {
      name: '',
    },
    fileMD5: '',
    chunks: [],
    // 2m
    chunkSize: 1024 * 1024 * 2,
    chunkCount: 0,
    sendChunkCount: 0,
    loading: true,
    // 分页
    page: 1,
    pages: 1,
    limit: 50,
    total: 0,
    order: 'date',
    desc: 1,
    fileType: '',
    q: '',
    diskMode: false,
    lists: [],
    // 上传面板
    // 扫描进度
    panelScanProgress: 0,
    // 所有文件上传进度
    panelUploadProgress: 0,
    panelChunkTotal: 0,
    showUploadPanel: false,
    uploadPanelList: [],
    // dialog
    showCreateFolderDialog: false,
    createFolderName: '',
    showNewFileNameDialog: false,
    newFileName: '',
    newFileNameFsId: 0,
    newFileNameDirId: 0,
  }),
  mounted() {
    this.backEnd = backEnd
    const dirId = parseInt(this.$route.query.dirId, 10) || 0
    this.dirId = dirId
    this.order = this.$route.query.order || 'date'
    const { fileType = '', q = '' } = this.$route.query
    if (!fileType && !q) {
      this.diskMode = true
    }
    this.fileType = fileType
    this.q = q
    this.refresh(dirId)
  },
  methods: {
    ...mapActions({
      showMsg: 'message/showMsg',
    }),
    // 刷新用户文件列表数据api
    async refresh(dirId = this.dirId || 0) {
      this.loading = true
      const {
        desc, order, page, pages, limit, fileType, q, diskMode,
      } = this
      if (page > pages + 1) return
      const queryData = {
        page,
        limit,
        desc,
        order,
        dirId,
      }
      let res = {}
      if (fileType) {
        queryData.fileType = fileType
        res = await api.user.fileListByFileType(queryData)
      } else if (q) {
        queryData.q = q
        res = await api.user.search(queryData)
      } else if (diskMode) {
        res = await api.user.fileList(queryData)
      }
      const { data } = res
      this.loading = false
      this.page = data.page
      this.total = data.total
      this.pages = data.pages
      this.lists = data.lists
    },
    // 相册预览
    showThumbnail(id) {
      let { lists } = this
      lists = lists.filter(i => i.file_type === 'image')
      const thumbnails = [];
      for (const item of lists) {
        thumbnails.push({
          id: item.id,
          thumbnail: `${backEnd}/image/${item.md5}${item.file_ext}?w=200`,
          source: `${backEnd}/image/${item.md5}${item.file_ext}`,
        });
      }
      const index = thumbnails.findIndex(i => i.id === id)
      this.$refs.viewer.show(
        thumbnails,
        index,
      );
    },
    // 选择文件并分片
    uploadFile(ev) {
      const _this = this
      const el = ev.target
      let allFileSize = 0
      // 文件列表
      for (let i = 0; i < el.files.length; i++) {
        const file = el.files[i]
        allFileSize += file.size
      }
      console.log('allFileSize', allFileSize)
      _this.panelChunkTotal = Math.ceil(allFileSize / _this.chunkSize)
      console.log(_this.panelChunkTotal)
      // TODO:
      for (let i = 0; i < el.files.length; i++) {
        _this.currentFile = el.files[i]
        this.sendFileData()
      }
    },
    // 上传文件夹
    uploadFolder(ev) {
      const el = ev.target
      console.log(el.files)
      this.showMsg({ type: 'warning', msg: 'TODO' })
    },
    // 上传文件数据
    async sendFileData() {
      const token = storage.get('gin_token')
      const _this = this
      const file = _this.currentFile
      const fileName = file.name
      const fileMD5 = await _this.getFileMD5(file)
      const data = {
        fileMD5,
      }
      // 秒传: 服务器检验是否有对应MD5值,如果有返回最后分片文件索引值
      const checkFileResult = await api.user.checkFile(data)
      console.log(checkFileResult)
      // //为异步操作添加可取消的功能
      // const cancelable = makeCancelable(_this.getFileMD5(file));
      // cancelable
      //   .promise
      //   .then((fileMD5) => console.log('resolved', fileMD5))
      //   .catch(({ isCanceled, ...error }) => console.log('isCanceled', isCanceled, error));
      // // 取消异步操作
      // console.log('cancel')
      // cancelable.cancel();

      const { chunkSize } = _this
      const chunks = Math.ceil(file.size / chunkSize)

      const blobSlice = File.prototype.slice || File.prototype.mozSlice || File.prototype.webkitSlice;
      // 设置当前分片索引值
      let currentChunk = checkFileResult.index || 0
      console.log(currentChunk, chunks)
      const fileReader = new FileReader();
      // 文件分片
      function loadChunk() {
        const start = currentChunk * chunkSize;
        const end = ((start + chunkSize) >= file.size) ? file.size : start + chunkSize;
        fileReader.readAsArrayBuffer(blobSlice.call(file, start, end));
      }
      // 分片上传
      async function uploadSliceFile(currentChunkFile, index) {
        const formData = new FormData()
        formData.append('file', new Blob([currentChunkFile]))
        formData.append('fileMD5', fileMD5)
        formData.append('index', index)
        formData.append('fileName', fileName)
        formData.append('fileSize', file.size)
        if (index === chunks) {
          formData.append('uploadType', 'merge')
          formData.append('dirId', _this.dirId)
          console.log('上传完成')
        } else {
          formData.append('uploadType', 'slice')
        }
        const resData = await axios.post(`${backEnd}/api/v1/sliceUpload`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
            Authorization: `${token}`,
          },
          transformRequest: [function (data) {
            return data
          }],
          onUploadProgress: (progressEvent) => {
            const complete = `${progressEvent.loaded / progressEvent.total * 100 || 0}%`
            console.log(`上传分片 ${index} - ${complete}`)
          },
        })
        console.log(resData)
        return resData
      }
      fileReader.onload = async (e) => {
        const complete = Math.floor(((_this.sendChunkCount) / _this.panelChunkTotal) * 10000) / 100;
        console.log('上传总进度==>', complete)
        _this.panelUploadProgress = complete
        const resData = await uploadSliceFile(e.target.result, currentChunk)
        console.log(resData)
        currentChunk++;
        _this.sendChunkCount++
        if (currentChunk <= chunks) {
          loadChunk();
        } else {
          console.log('上传完毕')
          _this.sendChunkCount = 0
          _this.panelChunkTotal = 0
          _this.panelUploadProgress = 0
          _this.refresh()
        }
      };
      fileReader.onerror = (e) => {
        console.log('读取文件时错误', e);
      };
      loadChunk();
    },
    // 扫描文件MD5
    getFileMD5(file) {
      const _this = this
      _this.showUploadPanel = true
      const { chunkSize } = _this
      const chunks = Math.ceil(file.size / chunkSize)
      return new Promise((resolve, reject) => {
        const blobSlice = File.prototype.slice || File.prototype.mozSlice || File.prototype.webkitSlice;
        let currentChunk = 0;
        const spark = new SparkMD5.ArrayBuffer();
        const fileReader = new FileReader();

        function loadChunk() {
          const start = currentChunk * chunkSize;
          const end = ((start + chunkSize) >= file.size) ? file.size : start + chunkSize;
          fileReader.readAsArrayBuffer(blobSlice.call(file, start, end));
        }

        fileReader.onload = (e) => {
          const panelScanProgress = Math.floor(((currentChunk + 1) / chunks) * 10000) / 100;
          // TODO: 延时渲染
          _this.panelScanProgress = panelScanProgress
          console.log(`扫描文件: ${panelScanProgress}`)
          // t.checkProgress = Math.floor(((currentChunk+1)/chunks) * 10000) / 100;
          spark.append(e.target.result);
          currentChunk++;
          if (currentChunk < chunks) {
            loadChunk();
          } else {
            console.log('文件扫描完成')
            const res = spark.end();
            _this.panelScanProgress = 0
            resolve(res);
          }
        };
        fileReader.onerror = (e) => {
          reject(new Error('读取文件时错误', e));
        };
        loadChunk();
      })
    },
    // 取消操作
    handleCancel() {
    },
    // 格式化文件大小
    bytesToSizeFormat(bytes) {
      return bytesToSize(bytes)
    },
    // 创建文件夹
    createFolderSubmit() {
      console.log('create')
      const data = {
        dirId: this.dirId,
        dirName: this.createFolderName,
      }
      api.user.createFolder(data)
        .then((resData) => {
          if (resData.code === 0) {
            this.showCreateFolderDialog = false
            this.lists.push(resData.data)
          }
        })
    },
    // 下载文件
    handleDownload() {
      const { selected } = this
      // date: 1575345141
      // file_ext: ".png"
      // file_name: "steam.png"
      // file_size: 601675
      // file_type: "image"
      // id: 2
      // md5: "fa4a971d7d1e56d685285a26ef22de38"
      for (let i = 0; i < selected.length; i++) {
        const file = selected[i]
        const serverFileName = `${file.md5}${file.file_ext}`
        const fileUrl = `${backEnd}/download/${serverFileName}`
        // download/818ab680f36528e25ebcb8433f23e0c3.jpg
        saveToDisk(fileUrl, file.file_name)
      }
    },
    handleDelete() {
      const { selected } = this
      let deleteFileArr = selected.map(i => i.id)
      deleteFileArr = Array.from(new Set(deleteFileArr))
      if (deleteFileArr.length === 0) return
      const data = {
        deleteFileArr,
      }
      api.user.deleteFileArr(data)
        .then((resData) => {
          if (resData.code === 0) {
            console.log(resData)
            this.refresh()
          }
        })
    },
    handleShare() {
      const { selected } = this
      console.log(selected)
      this.showMsg({ type: 'warning', msg: 'TODO' })
    },
    handleCopy() {
      const { selected } = this
      console.log(selected)
      this.showMsg({ type: 'warning', msg: 'TODO' })
    },
    handleMove() {
      const { selected } = this
      console.log(selected)
      this.showMsg({ type: 'warning', msg: 'TODO' })
    },
    // 文件项目操作列表
    downloadItem(item) {
      const serverFileName = `${item.md5}${item.file_ext}`
      const fileUrl = `${backEnd}/download/${serverFileName}`
      // download/818ab680f36528e25ebcb8433f23e0c3.jpg
      saveToDisk(fileUrl, item.file_name)
    },
    shareItem(item) {
      console.log(item)
      this.showMsg({ type: 'warning', msg: 'TODO' })
    },
    copyItem(item) {
      console.log(item)
      this.showMsg({ type: 'warning', msg: 'TODO' })
    },
    moveItem(item) {
      console.log(item)
      this.showMsg({ type: 'warning', msg: 'TODO' })
    },
    renameItem(item) {
      this.showNewFileNameDialog = true
      this.newFileName = item.file_name
      this.newFileNameFsId = item.id
      this.newFileNameDirId = item.parent
    },
    renameItemSubmit() {
      const data = {
        fsId: this.newFileNameFsId,
        dirId: this.newFileNameDirId,
        fileName: this.newFileName,
      }
      console.log(data)
      api.user.rename(data)
        .then((resData) => {
          if (resData.code === 0) {
            console.log(resData)
            this.showNewFileNameDialog = false
            this.refresh()
          }
        })
    },
    async goToItem(item) {
      console.log(item)
      const data = {
        fs_id: item.id,
      }
      const resData = await api.user.findFilePath(data)
      console.log(resData)
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
    async deleteItem(item) {
      const data = {
        fs_id: item.id,
      }
      api.user.deleteFile(data)
        .then((resData) => {
          if (resData.code === 0) {
            this.refresh()
          }
        })
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
  // watch: {
  //   $route(to, from) {
  //     this.refresh(to.query.dirId)
  //   },
  // },
}
</script>
