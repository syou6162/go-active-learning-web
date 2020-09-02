<template>
  <b-list-group-item class="tweet-item">
    <img
      v-if="tweet.ProfileImageUrl"
      :src="tweet.ProfileImageUrl"
      class="tweet-full-text-icon-img"
      onerror="this.src='/img/twitter_icon.png'"
    >
    <div class="tweet-screen-name-and-full-text-container">
      <a
        :href="'https://twitter.com/' + tweet.ScreenName + '/status/' + tweet.IdStr"
        target="_blank"
        rel="noopener"
      >@{{ tweet.ScreenName }}</a>
      <span v-html="fullTextWithLinks(tweet.FullText)" />
      <div class="tweet-footer">
        <span class="tweet-retweet-count">{{ tweet.RetweetCount }} RT</span>, <span class="tweet-favorite-count"> {{ tweet.FavoriteCount }} Fav</span>
        <span class="tweet-created-at">{{ tweet.CreatedAt }}</span>
      </div>
    </div>
    <tweet-annotate-buttons
      v-if="isAdmin"
      :tweet="tweet"
    />
  </b-list-group-item>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Autolinker, AutolinkerConfig } from 'autolinker';
import Tweet from '~/models/Tweet'
import TweetAnnotateButtons from '~/components/TweetAnnotateButtons.vue'

@Component({
  components: {
    TweetAnnotateButtons
  }
})
export default class ReferringTweet extends Vue {
  @Prop()
  tweet!: Tweet

  @Prop({required: true, default: false})
  isAdmin!: Boolean
  
  fullTextWithLinks(fullText: string): string {
    const opts: AutolinkerConfig = { 
      mention: 'twitter',
      hashtag: 'twitter'
    };
    return Autolinker.link(fullText, opts);
  }
}
</script>

<style scoped>
.tweet-item {
  padding: 10px 10px;
}
.tweet-screen-name-and-full-text-container {
  overflow: hidden;
}
.tweets-count, .tweet-retweet-count, .tweet-favorite-count {
  color: #ff4166;
}
.tweet-full-text-icon-img {
  float: left;
  width: 32px;
  height: 32px;
  margin: 0 10px 0 0;
}
.tweet-created-at {
  color: #999;
  float: right;
}
</style>