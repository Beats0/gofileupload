<template>
    <viewer :images="images"
            :options="options"
            class="viewer"
            ref="viewer"
            @inited="inited"
            v-if="images && images.length">
      <img v-for="{source, thumbnail} in images"
           :src="thumbnail"
           :data-source="source"
           :key="source"
           class="image">
    </viewer>
</template>

<script>
import Vue from 'vue';
import Viewer from 'v-viewer';
import 'viewerjs/dist/viewer.css';

Vue.use(Viewer);

export default {
  data() {
    return {
      options: {
        url: 'data-source',
      },
      index: 0,
      images: [],
    };
  },
  methods: {
    inited(viewer) {
      this.$viewer = viewer;
      this.$viewer.view(this.index);
    },
    view(index) {
      this.index = index
      this.$viewer.view(this.index);
    },
    show(images, index = 0) {
      if (this.images === images) {
        this.view(index)
        return
      }
      this.images = images
      this.index = index
    },
  },
};
</script>

<style scoped>
  .image {
    display: none;
  }
  .viewer-loading > img {
    /* hide big images when it is loading */
    display: none;
  }
</style>
