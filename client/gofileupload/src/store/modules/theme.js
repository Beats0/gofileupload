import * as types from '../types'
import Storage from '../../utils/storage'

const state = {
  theme: Storage.get('theme') || 'light',
  drawer: null,
}

const mutations = {
  [types.LIGHT](state) {
    state.theme = 'light'
  },
  [types.DARK](state) {
    state.theme = 'dark'
  },
  [types.SETDRAWER](state, visible) {
    state.drawer = visible
  },
}

const actions = {
  changeTheme({ commit }) {
    if (state.theme === 'light') {
      Storage.set('theme', 'dark')
      commit(types.DARK)
    } else {
      Storage.set('theme', 'light')
      commit(types.LIGHT)
    }
  },
  changeDrawer({ commit }) {
    commit(types.SETDRAWER, !state.drawer)
  },
}

const getters = {
  isDark() {
    return state.theme === 'dark'
  },
  drawerVisible() {
    return state.drawer
  },
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
}
