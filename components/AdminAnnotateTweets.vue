<template>
  <div>
    <div
      v-for="tweet in tweets"
      :key="tweet.IdStr"
      :tweet="tweet"
      :isAdmin="isAdmin"
    >
      <b-card no-body>
        <b-card-body
          :title="example | getTitle(75, '...')"
          title-tag="h5"
          class="m-1 p-2"
        >
          <img
            v-if="example.OgImage"
            v-lazy="example.OgImage"
            class="img-thumbnail img-responsive ogimage"
            onerror="this.style.display='none'"
          >
          <img
            v-if="tweet.ProfileImageUrl"
            v-lazy="tweet.ProfileImageUrl"
            class="tweet-profile-img"
            onerror="this.style.display='none'"
          >
          <a :href="'https://twitter.com/' + tweet.ScreenName + '/status/' + tweet.IdStr">@{{ tweet.ScreenName }}</a>
          <p class="card-text">
            {{ tweet.FullText }}
          </p>
          Fav: {{ tweet.FavoriteCount }}, RT: {{ tweet.RetweetCount }}, Score: {{ tweet.Score.toFixed(1) }}
          <tweet-annotate-buttons
            v-if="isAdmin"
            :tweet="tweet"
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
              class="example-favicon-img"
              onerror="this.style.display='none'"
            >
            <a :href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
          </b-card-footer>
        </b-card-body>
      </b-card>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Example, getTitle, getDomain, getUserName, getExampleUrl } from '~/models/Example'
import TweetAnnotateButtons from '~/components/TweetAnnotateButtons.vue'
import Tweet from '~/models/Tweet'

@Component({
  components: {
    TweetAnnotateButtons,
  },
  filters: {
    getTitle(example: Example, length: number, omission: string): string {
      return getTitle(example, length, omission)
    },
    getDomain(example: Example): string {
      return getDomain(example)
    },
    getUserName(example: Example): string {
      return getUserName(example)
    },
    getExampleUrl(example: Example): string {
      return getExampleUrl(example)
    },
  }  
})

export default class AdminAnnotateTweets extends Vue {
  modalShow: boolean = false

  @Prop({required: true})
  example!: Example 
  
  @Prop({required: true})
  tweets!: Tweet[]

  @Prop({required: true, default: false})
  isAdmin!: boolean
}
</script>

<style scoped>
.example-favicon-img {
  width: 16px;
  height: 16px;
}
.tweet-profile-img {
  width: 24px;
  height: 24px;
  margin: 2px;
}
.ogimage {
  width: 128px;
  height: 96px;
  margin: 3px;
  float: right;
}
</style>
