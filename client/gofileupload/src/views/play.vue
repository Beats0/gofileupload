<template>
  <v-app id="inspire">
  <Drawer />
  <AppBar />
  <div class="playerContainer">
    <div class="playerTitle">{{ videoInfo.file_namel }}</div>
    <div id="dplayer"></div>
    <div class="videoInfoContainer">
      <div class="videoInfo">
        <span class="videoLabel videoTitle">{{ videoInfo.file_name }}</span>
        <span class="videoLabel videoSize">大小: {{ bytesToSizeFormat(videoInfo.file_size) }}</span>
        <span class="videoLabel videoTime">上传时间: {{ videoInfo.last_modified * 1000 | dayjsFormat }}</span>
      </div>
      <div class="videoAction">
        <v-btn color="primary" outlined small style="margin-right: 10px" @click="handleDownload"> 下载 <v-icon right dark>cloud_download</v-icon></v-btn>
        <v-btn color="primary" outlined small style="margin-right: 10px" @click="handleShare"> 分享 <v-icon right dark>share</v-icon></v-btn>
      </div>
    </div>
  </div>
</v-app>
</template>

<script>
import DPlayer from 'dplayer';
import SparkMD5 from 'spark-md5'
import api from '../../api'
import Drawer from '../components/drawer.vue'
import AppBar from '../components/appBar.vue'
import { backEnd, checkSecretKey } from '../../config';
import 'dplayer/dist/DPlayer.min.css';
import { bytesToSize } from '../utils/format';
import { saveToDisk } from '../utils/util';

export default {
  name: 'play',
  data() {
    return {
      backEnd: '',
      v: '',
      videoInfo: {
        id: 0,
        date: 0,
        file_name: '1.mp4',
        file_size: '332791844',
        file_ext: '.mp4',
        file_type: 'media',
        md5: '86fa76e13bf6401053cd19c43cedfb8e',
        is_dir: 0,
        last_modified: '1577775426',
      },
    }
  },
  components: {
    AppBar,
    Drawer,
  },
  mounted() {
    this.backEnd = backEnd
    const { v } = this.$route.query
    this.v = v
    this.getVideoInfo(v)
    const exptime = parseInt((Date.now() / 1000), 10)
    const checkStr = `${checkSecretKey}${backEnd}/video/${v}${exptime}`
    const secret = SparkMD5.hash(checkStr)
    const requestUrl = `${backEnd}/video/${v}?exptime=${exptime}&secret=${secret}`
    this.initPlayer(requestUrl)
  },
  methods: {
    getVideoInfo(v) {
      const data = {
        v,
      }
      api.user.videoInfo(data).then((resData) => {
        this.videoInfo = resData.data
      })
    },
    bytesToSizeFormat(bytes) {
      return bytesToSize(bytes)
    },
    initPlayer(url) {
      const dp = new DPlayer({
        container: document.getElementById('dplayer'),
        // TODO: add thumbnails
        video: {
          url,
          // pic: '',
          // thumbnails: '',
        },
      })
    },
    handleDownload() {
      const { md5 } = this.videoInfo
      const fileExt = this.videoInfo.file_ext
      const fileName = this.videoInfo.file_name
      const serverFileName = `${md5}${fileExt}`
      const fileUrl = `${backEnd}/download/${serverFileName}`
      saveToDisk(fileUrl, fileName)
    },
    handleShare() {
    },
  },
}
</script>
