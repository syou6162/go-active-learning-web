<template>
  <img 
    v-b-popover.hover.bottomright="tweet.FullText" 
    v-lazy="tweet.ProfileImageUrl" 
    :alt="twitterId"
    :title="twitterId" 
    class="tweet-icon-img"
  >
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Tweet from '~/models/Tweet'

@Component
export default class TwitterIcon extends Vue {
  @Prop()
  tweet!: Tweet;

  get twitterId(): string {
    return '@' + this.tweet.ScreenName;
  }
}
</script>

<style scoped>
.tweet-icon-img {
  width: 24px;
  height: 24px;
  margin: 2px;
}
</style>