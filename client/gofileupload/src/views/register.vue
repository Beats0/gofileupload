<template>
  <v-app id="inspire">
    <v-content>
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md4>
            <v-card class="elevation-12">
              <v-toolbar color="primary" dark flat>
                <v-toolbar-title>Register</v-toolbar-title>
                <v-spacer></v-spacer>
              </v-toolbar>
              <v-card-text>
                <v-form>
                  <v-text-field
                    label="UserName"
                    name="username"
                    prepend-icon="person"
                    type="text"
                    v-model="uname"
                  />
                  <v-text-field
                    label="Mail"
                    name="login"
                    prepend-icon="mail_outline"
                    type="text"
                    v-model="mail"
                  />
                  <v-text-field
                    id="password"
                    label="Password"
                    name="password"
                    prepend-icon="lock"
                    type="password"
                    v-model="pwd"
                  />
                </v-form>
              </v-card-text>
              <v-card-actions style="padding-bottom: 20px;">
                <v-spacer></v-spacer>
                <v-btn color="primary" small style="margin-right: 10px" @click="handleRegister">Register</v-btn>
                <v-btn color="primary" small outlined style="margin-right: 10px" @click="goToLogin">Login</v-btn>
              </v-card-actions>
            </v-card>
          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import api from '../../api'

export default {
  name: 'register',
  data() {
    return {
      uname: '',
      mail: '',
      pwd: '',
    }
  },
  methods: {
    goToLogin() {
      this.$router.push({ path: 'login' })
    },
    handleRegister() {
      const { uname, mail, pwd } = this
      const data = {
        uname,
        mail,
        pwd,
      }
      api.user.register(data)
        .then((resData) => {
          if (resData.code === 0 && resData.data.token) {
            window.location.reload()
          }
        })
    },
  },
};
</script>
<style>
  #inspire {
    background-image: url("https://steamuserimages-a.akamaihd.net/ugc/775101441378770069/72E8A2AFB6BCB120F82BD4725D8872074007E41A/");
    background-repeat: no-repeat;
    background-position: center;
    background-size: cover;
  }
</style>
