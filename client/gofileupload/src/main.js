import Vue from 'vue'
import dayjs from 'dayjs'
import Viewer from 'v-viewer'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify'
import store from './store'


Vue.config.productionTip = false;

Vue.use(Viewer)
Vue.filter('dayjsFormat', (date, template = 'YYYY-MM-DD hh:mm:ss') => dayjs(date).format(template))


new Vue({
  router,
  store,
  vuetify,
  render: h => h(App),
}).$mount('#app')
