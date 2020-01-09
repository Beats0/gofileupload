import Vue from 'vue';
import Router from 'vue-router';
import store from './store'
import storage from './utils/storage';

const originalPush = Router.prototype.push
Router.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}

const Login = () => import(/* webpackChunkName: "login" */ './views/login.vue');
const Register = () => import(/* webpackChunkName: "register" */ './views/register.vue');
const Disk = () => import(/* webpackChunkName: "disk" */ './views/disk.vue');
const Share = () => import(/* webpackChunkName: "share" */ './views/share.vue');
const Play = () => import(/* webpackChunkName: "play" */ './views/play.vue');
const Notfound = () => import(/* webpackChunkName: "notfound" */ './views/notfound.vue');


Vue.use(Router);

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Disk,
      meta: { checkLogin: true },
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { loginRedirect: true },
    },
    {
      path: '/register',
      name: 'register',
      component: Register,
      meta: { loginRedirect: true },
    },
    {
      path: '/disk',
      name: 'disk',
      component: Disk,
      meta: { checkLogin: true },
    },
    {
      path: '/share',
      name: 'share',
      component: Share,
      meta: { checkLogin: true },
    },
    {
      path: '/play',
      name: 'play',
      component: Play,
      meta: { checkLogin: true },
    },
    {
      path: '*',
      name: 'notfound',
      component: Notfound,
    },
  ],
});


// 全局身份确认
router.beforeEach((to, from, next) => {
  if (to.meta.checkLogin) {
    const token = storage.get('gin_token')
    if (token) {
      next()
    } else {
      store.commit('session/SHOWLOGIN')
      next({
        path: '/login',
      })
    }
  } else if (to.meta.loginRedirect) {
    const token = storage.get('gin_token')
    if (token) {
      next({
        path: '/',
      })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
