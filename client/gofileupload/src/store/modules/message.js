import * as types from '../types'

const state = {
  msg: '',
  msgType: '',
}

const mutations = {
  [types.SHOWMSG](state, data) {
    state.msg = data.msg
    state.msgType = data.msgType
  },
  [types.CLOSEMSG](state) {
    state.msg = ''
  },
}

const actions = {
  showMsg({ commit }, data) {
    commit(types.SHOWMSG, data)
  },
  closeMsg({ commit }) {
    commit(types.CLOSEMSG)
  },
}

const getters = {
  hasMsg() {
    return state.msg !== ''
  },
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
}
