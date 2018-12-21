<template>
  <div>
    <b-card no-body>
      <b-card-body v-bind:title="example | getTitle(75, '...')" title-tag="h5" class="m-1 p-2">
        <div v-if="example.HatenaBookmark.count > 0">
          <a v-bind:href="example.HatenaBookmark.entry_url" style="color: #ff4166;">{{ example.HatenaBookmark.count }} users</a>:
          <hatena-bookmark-icon 
            v-for="b in example.HatenaBookmark.bookmarks.slice(0, 9)"
            v-bind:bookmark="b"
            v-bind:key="b.user"
            ></hatena-bookmark-icon>
        </div>
        <div v-if="tweets !== undefined && tweets !== null && tweets.length > 0">
          <span style="color: #ff4166;">{{ tweets.length }} mentions</span>:
          <twitter-icon
             v-for="tweet in tweets.slice(0, 8)"
             v-bind:tweet="tweet"
             v-bind:key="tweet.ScreenName"
             ></twitter-icon>
        </div>
        Date: 
        {{ example.CreatedAt.tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm") }}
        <annotate-buttons
          v-if="isAdmin"
          v-bind:example="example"
          ></annotate-buttons>
        <b-card-footer>
          <b-button v-bind:href="'/example/' + encodeURIComponent(example.FinalUrl)" class="float-right" size="sm" v-bind:variant="example | getButtonStyle">Read more</b-button>
          <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
          <a v-bind:href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
        </b-card-footer>
      </b-card-body>
    </b-card>
  </div>
</template>

<script>
import HatenaBookmarkIcon from './HatenaBookmarkIcon.vue';
import TwitterIcon from './TwitterIcon.vue';
import AnnotateButtons from './AnnotateButtons.vue';

export default {
  data () {
    return {
      modalShow: false
    }
  },
  props: ['example', 'tweets', 'isAdmin'],
  components: {
    "hatena-bookmark-icon": HatenaBookmarkIcon,
    "twitter-icon": TwitterIcon,
    "annotate-buttons": AnnotateButtons 
  }
}
</script>
