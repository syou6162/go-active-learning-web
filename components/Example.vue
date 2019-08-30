<template>
  <div>
    <b-card no-body>
      <b-card-body
        :title="example | getTitle(75, '...')"
        title-tag="h5"
        class="m-1 p-2"
      >
        <p class="example-description">
          {{ example | getDescription(100, '...') }}
        </p>
        <div 
          v-if="example.HatenaBookmark.count > 0"
          class="hatena-bookmarks"
          >
          <a
            :href="example.HatenaBookmark.entry_url"
            class="hatena-bookmark-link"
          >{{ example.HatenaBookmark.count }} users</a>:
          <hatena-bookmark-icon 
            v-for="b in example.HatenaBookmark.bookmarks.slice(0, 9)"
            :key="b.user"
            :bookmark="b"
          />
        </div>
        <div 
          v-if="tweets !== undefined && tweets !== null && tweets.length > 0"
          class="tweets"
          >
          <span class="tweets-count">{{ tweets.length }} mentions</span>:
          <twitter-icon
            v-for="tweet in tweets.slice(0, 8)"
            :key="tweet.ScreenName"
            :tweet="tweet"
          />
        </div>
        <div class="example-created-at">{{ example.CreatedAt }}</div>
        <annotate-buttons
          v-if="isAdmin"
          :example="example"
        />
        <b-card-footer>
          <b-button
            :to="example | getExampleUrl"
            class="float-right"
            size="sm"
          >
            Read more
          </b-button>
          <img
            v-if="example.Favicon"
            v-lazy="example.Favicon"
            :alt="example.Title"
            class="example-favicon-img"
            onerror="this.style.display='none'"
          >
          <a :href="example.FinalUrl" target="_blank" rel="noopener">{{ example | getDomain }} {{ example | getUserName }}</a>
        </b-card-footer>
      </b-card-body>
    </b-card>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Example from '~/models/Example'
import Tweet from '~/models/Tweet'

@Component({
  components: {
    HatenaBookmarkIcon: () => import('./HatenaBookmarkIcon.vue'),
    TwitterIcon: () => import('./TwitterIcon.vue'),
    AnnotateButtons: () => import('./AnnotateButtons.vue'),
  }
})

export default class ExampleComponent extends Vue {
  modalShow: boolean = false

  @Prop({required: true})
  example!: Example 

  @Prop({required: false})
  tweets!: Tweet[]

  @Prop({required: true, default: false})
  isAdmin!: Boolean
}
</script>

<style scoped>
.example-favicon-img {
  width: 16px; 
  height: 16px;
}
.example-description {
  font-size: 16px;
  line-height: 18px;
  margin: 0 0 4px;
  color: #55606a;
}
.example-created-at {
  color: #999;
  text-align: right;
  margin: 0 0 4px;
  line-height: 16px;
}
.hatena-bookmarks {
  margin: 0 0 4px;
}
.hatena-bookmark-link {
  color: #ff4166;
}
.tweets {
  margin: 0 0 4px;
}
.tweets-count {
  color: #ff4166;
}
</style>