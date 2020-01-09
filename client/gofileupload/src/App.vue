<template>
  <div id="app">
    <Message />
    <router-view :key="key" />
  </div>
</template>

<script>
import Message from '@/components/message.vue'

export default {
  created() {
  },
  components: {
    Message,
  },
  computed: {
    // router-view 不刷新问题
    // https://blog.csdn.net/w390058785/article/details/82813032
    key() {
      return this.$route.path + Math.random();
    },
  },
}
</script>
