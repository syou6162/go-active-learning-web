<template>
  <div>
    <b-card no-body>
      <b-card-body v-bind:title="example | getTitle(75, '...')" title-tag="h5" class="m-1 p-2">
        <div v-if="example.HatenaBookmark.bookmarks.length > 0">
          <a v-bind:href="example.HatenaBookmark.entry_url" style="color: #ff4166;">{{ example.HatenaBookmark.bookmarks.length}} users</a>:
          <hatena-bookmark-icon 
            v-for="b in example.HatenaBookmark.bookmarks.slice(0, 9)"
            v-bind:bookmark="b"
            ></hatena-bookmark-icon>
        </div>
        <b-card-footer>
          <b-button v-bind:href="'/example/' + encodeURIComponent(example.FinalUrl)" class="float-right" size="sm">Read more</b-button>
          <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
          <a v-bind:href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
        </b-card-footer>
      </b-card-body>
    </b-card>
  </div>
</template>

<script>
import HatenaBookmarkIcon from './HatenaBookmarkIcon.vue';
export default {
  data () {
    return {
      modalShow: false
    }
  },
  props: ['example'],
  components: {
    "hatena-bookmark-icon": HatenaBookmarkIcon 
  }
}
</script>
