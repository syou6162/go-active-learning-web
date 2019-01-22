<template>
  <div>
    <div
      v-for="tweet in tweets"
      v-bind:key="tweets.IdStr"
      v-bind:tweet="tweet"
      v-bind:isAdmin="isAdmin"
      >
      <b-card no-body>
        <b-card-body v-bind:title="example | getTitle(75, '...')" title-tag="h5" class="m-1 p-2">
          <img v-if="example.OgImage" class="img-thumbnail img-responsive" style="width: 128px; height: 96px; margin: 3px; float: right;" v-lazy="example.OgImage" onerror="this.style.display='none'" />
          <img v-if="tweet.ProfileImageUrl" style="width: 24px; height: 24px;" v-lazy="tweet.ProfileImageUrl" onerror="this.style.display='none'" />
          <a v-bind:href="'https://twitter.com/' + tweet.ScreenName + '/status/' + tweet.IdStr">@{{ tweet.ScreenName }}</a>
          <p class="card-text">
            {{ tweet.FullText }}
          </p>
          Fav: {{ tweet.FavoriteCount }}, RT: {{ tweet.RetweetCount }}, Score: {{ tweet.Score.toFixed(1) }}
          <tweet-annotate-buttons v-if="isAdmin" v-bind:tweet="tweet"></tweet-annotate-buttons>
          <b-card-footer>
            <b-button v-bind:to="example | getExampleUrl" class="float-right" size="sm">Read more</b-button>
            <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
            <a v-bind:href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
          </b-card-footer>
        </b-card-body>
      </b-card>
    </div>
  </div>
</template>

<script>
import TwitterIcon from './TwitterIcon.vue';
import TweetAnnotateButtons from '~/components/TweetAnnotateButtons.vue';

export default {
  data () {
    return {
      modalShow: false
    }
  },
  props: ['example', 'tweets', 'isAdmin'],
  components: {
    "twitter-icon": TwitterIcon,
    "tweet-annotate-buttons": TweetAnnotateButtons
  }
}
</script>
