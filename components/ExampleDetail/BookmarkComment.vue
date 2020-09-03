<template>
  <b-list-group-item
    class="hatena-bookmark-item"
  >
    <img
      v-lazy="'https://cdn.profile-image.st-hatena.com/users/' + bookmark.user+ '/profile.png'"
      class="hatena-bookmark-comment-user-icon-img"
    >
    <div class="hatena-bookmark-user-link-and-comment-container">
      <a 
        :href="'http://b.hatena.ne.jp/' + bookmark.user"
        target="_blank"
        rel="noopener"
        class="hatena-bookmark-user-link"
      >
        id:{{ bookmark.user }}
      </a>
      <span class="hatena-bookmark-comment">
        {{ bookmark.comment }}
      </span>
      <span class="hatena-bookmark-timestamp">
        {{ bookmark.timestamp }}
      </span>
    </div>
  </b-list-group-item>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Bookmark from '~/models/Bookmark'

@Component
export default class BookmarkComment extends Vue {
  @Prop()
  bookmark!: Bookmark
}
</script>

<style scoped>
.hatena-bookmark-item {
  padding: 10px 10px;
}
.hatena-bookmark-comment-user-icon-img {
  float: left;
  width: 32px;
  height: 32px;
  margin: 0 10px 0 0;
}
.hatena-bookmark-user-link-and-comment-container {
  overflow: hidden;
}
.hatena-bookmark-user-link {
  margin: 0 0 0 0;
}
.hatena-bookmark-comment {
  margin: 0 4px 0 0;
}
.hatena-bookmark-timestamp {
  display: block;
  color: #999;
  float: right;
}
</style>