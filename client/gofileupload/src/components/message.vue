<template>
  <v-snackbar v-model="show"
              top
              :color="msgType"
              :timeout="timeout">
    {{ msg }}
    <v-btn fab x-small @click="closeMsg">
      <v-icon color="#F44336">close</v-icon>
    </v-btn>
  </v-snackbar>
</template>

<script>
import { mapState, mapActions, mapGetters } from 'vuex'

export default {
  data() {
    return {
      timeout: 5000,
    }
  },
  computed: {
    ...mapState({
      msg: state => state.message.msg,
      msgType: state => state.message.msgType,
    }),
    ...mapGetters({
      hasMsg: 'message/hasMsg',
    }),
    show: {
      get() {
        return this.hasMsg
      },
      set(value) {
        if (!value) {
          this.closeMsg()
        }
      },
    },
  },
  methods: {
    ...mapActions({
      closeMsg: 'message/closeMsg',
    }),
  },
}
</script>
