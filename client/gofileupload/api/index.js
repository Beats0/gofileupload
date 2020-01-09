import axios from 'axios'
import { backEnd } from '../config';
import storage from '../src/utils/storage'
import store from '../src/store'

axios.defaults.baseURL = `${backEnd}/api`

// 跨域
axios.defaults.withCredentials = process.env.NODE_ENV === 'production'
// 最大响应超时
axios.defaults.timeout = 100000
axios.defaults.headers.post['Content-Type'] = 'multipart/form-data'


// jwt
axios.interceptors.request.use(
  async (config) => {
    const token = await storage.get('gin_token')
    if (token) {
      // TODO: token过期校验
      // const { token } = session
      // const expired = Number(session.expired)
      // console.log('expired', expired)
      // // 时间戳过期检验
      // if (token && expired > Date.now()) {
      //   // Bearer是JWT的认证头部信息
      //   config.headers.common.Authorization = `Bearer ${token}`;
      // }
      config.headers.common.Authorization = `${token}`;
    }
    return config;
  },
  error => Promise.reject(error),
);

const instance = {};
['get', 'post', 'put', 'delete'].forEach((key) => {
  instance[key] = function (...args) {
    console.log([key], axios.defaults.baseURL, ...args)
    return axios[key](...args)
      .then((res) => {
        if (res.data && res.data.data && res.data.data.token) {
          storage.set('gin_token', res.data.data.token)
        }
        if (res.data.code !== 0 && res.data.msg) {
          res.data.msgType = 'warning'
          store.commit('message/SHOWMSG', res.data)
        }
        return res.data
      })
      .catch((err) => {
        // 全局信息提示
        if (err.response.data.code !== 200 && err.response.data.code !== 0 && err.response.data.msg) {
          err.response.data.msgType = 'error'
          store.commit('message/SHOWMSG', err.response.data)
        }
        console.log(err)
      })
  }
})

const api = {
  user: {
    login: data => instance.post('/v2/user/login', data),
    register: data => instance.post('/v2/user/register', data),
    porfile: () => instance.get('/v1/user/porfile'),
    checkFile: data => instance.get('/v1/checkFile', { params: data }),
    fileList: data => instance.get('/v1/user/fileList', { params: data }),
    fileListByFileType: data => instance.get('/v1/user/fileListByFileType', { params: data }),
    createFolder: data => instance.post('/v1/createFolder', data),
    findFilePath: data => instance.post('/v1/findFilePath', data),
    deleteFile: data => instance.post('/v1/deleteFile', data),
    search: data => instance.post('/v1/search', data),
    rename: data => instance.post('/v1/rename', data),
    deleteFileArr: data => instance.post('/v1/deleteFileArr', data),
    serverStatic: () => instance.get('/v1/serverStatic'),
    videoInfo: (data) => instance.get('/v1/videoInfo', { params: data }),
  },
  upload: () => instance.post('/v1/upload'),
  multiUpload: () => instance.post('/v1/multiUpload'),
}

export default api
