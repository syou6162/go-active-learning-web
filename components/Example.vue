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
        Date: {{ example.CreatedAt }}
        <annotate-buttons
          v-if="isAdmin"
          v-bind:example="example"
          ></annotate-buttons>
        <b-card-footer>
          <b-button v-bind:to="example | getExampleUrl" class="float-right" size="sm">Read more</b-button>
          <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
          <a v-bind:href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
        </b-card-footer>
      </b-card-body>
    </b-card>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { IsAdmin } from '../plugins/amplify';
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
  example!: Object

  @Prop({required: false})
  tweets!: Array<Object>

  @Prop({required: true, default: false})
  isAdmin!: Boolean
}
</script>