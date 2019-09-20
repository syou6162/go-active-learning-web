<template>
  <div>
    <p>ログイン</p>
    <p>{{ status }}</p>
    <p>{{ messageText }} </p>
    <label>ユーザー名</label>
    <input
      v-model="userInfo.username"
      type="text"
    >
    <label>パスワード</label>
    <input
      v-model="userInfo.password"
      type="password"
      data-private="lipsum"
    >
    <button
      class="btn btn-primary"
      @click="signIn()"
    >
      ログイン
    </button>
    <button
      class="btn btn-primary"
      @click="signOut()"
    >
      ログアウト
    </button>
    <div id="annotation-links">
      <ul>
        <li><a href="admin/annotation/examples">Exampleのアノテーション</a></li>
        <li><a href="admin/annotation/tweets">Tweetのアノテーション</a></li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { IsAdmin, signIn, signOut } from '~/plugins/amplify';

@Component
export default class AdminIndexPage extends Vue {
  title: string = "ML-News"
  error: string | null = null
  loading: boolean = true
  status: string = ''
  userInfo: { [key: string]: string } = {
    username: '',
    password: ''
  }
  messageText: string = ''

  created() {
    IsAdmin().then(isAdmin => {
      if (isAdmin) {
        this.status = 'ログインしています';
      } else {
        this.status = 'ログインしていません';
      }
    });
  }
  signIn() {
    signIn(this.userInfo.username, this.userInfo.password)
    .then((data) => {
      this.messageText = 'ログインしました';
      this.status = 'こんにちは、' + data.username + 'さん';
    }).catch((err) => {
      this.messageText = 'ログインできませんでした';
    });
  }
  signOut() {
    signOut()
    .then((data) => {
      this.messageText = 'ログアウトしました';
      this.status = 'またきてね';
    }).catch((err) => {
      this.messageText = 'ログインできませんでした';
    });
  }
}
</script>
