<template>
  <div>
    <p>ログイン</p>
    <p>{{ status }}</p>
    <p>{{ message_text }} </p>
    <label>ユーザー名</label>
    <input type="text" v-model="userInfo.username"/>
    <label>パスワード</label>
    <input type="password" v-model="userInfo.password"/>
    <button class="btn btn-primary" @click="signIn()">ログイン</button>
    <button class="btn btn-primary" @click="signOut()">ログアウト</button>
  </div>
</template>

<script>
import { IsAdmin, signIn, signOut } from './util';

export default {
  data () {
    return {
      title: "ML News",
      listname: 'general',
      examples: [],
      isNew: 0,
      options: [
        { text: 'All', value: 0 },
        { text: 'Recent', value: 1 },
      ],
      error: null,
      loading: true,
      status: '',
      userInfo: {
        username: '',
        password: '',
      },
      message_text: '',
    }
  },
  created() {
    IsAdmin().then(isAdmin => {
      if (isAdmin) {
        this.status = 'ログインしています';
      } else {
        this.status = 'ログインしていません';
      }
    });
  },
  methods: {
    signIn: function () {
      signIn(this.userInfo.username, this.userInfo.password)
      .then((data) => {
        this.message_text = 'ログインしました';
        this.status = 'こんにちは、'+data.username+'さん';
      }).catch((err) => {
        this.message_text = 'ログインできませんでした';
      });
    },
    signOut: function () {
      signOut()
      .then((data) => {
        this.message_text = 'ログアウトしました';
        this.status = 'またきてね';
      }).catch((err) => {
        this.message_text = 'ログインできませんでした';
      });
    },
  }
}
</script>
