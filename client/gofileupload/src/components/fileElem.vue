<template>
  <div class="table-item-preview">
    <img :src="`${backEnd}/image/${file.md5}${file.file_ext}?w=200`"
         :alt="file.file_name"
         :title="file.file_name"
         class="image-preview"
         v-if="file.file_type === 'image'"
         @dblclick="viewDetail(file.id, 'image', file.md5)"/>
    <v-icon @dblclick="viewDetail(file.id, file.file_type, file.md5)" v-else>{{ renderFileIcon(file.is_dir, file.file_type) }}</v-icon>
  </div>
</template>
6
<script>
import { backEnd } from '../../config';

export default {
  name: 'FileElem',
  props: {
    file: {
      // file_ext:".png"
      // file_name:"redux.png"
      // file_size:591962
      // file_type:"image"
      // id:1
      // is_dir: 0
      // md5:"818ab680f36528e25ebcb8433f23e0c3"
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      backEnd: '',
    }
  },
  mounted() {
    this.backEnd = backEnd
  },
  methods: {
    renderFileIcon(isDir, fileType) {
      let iconName = ''
      if (isDir) return 'folder_open'
      switch (fileType) {
        case 'audio':
          iconName = 'music_note'
          break
        case 'media':
          iconName = 'videocam'
          break
        case 'doc':
          iconName = 'edit'
          break
        case 'seed':
          iconName = 'insert_link'
          break
        default:
          iconName = 'cloud_download'
          break
      }
      return iconName
    },
    viewDetail(fid, fileType, md5) {
      switch (fileType) {
        case 'image':
          this.$emit('vewLarger', fid);
          break
        case 'media':
          const url = this.$router.resolve({
            path: '/play',
            query: { v: md5 },
          });
          window.open(url.href, '_blank')
          break
        default:
          break
      }
    },
  },
}
</script>
