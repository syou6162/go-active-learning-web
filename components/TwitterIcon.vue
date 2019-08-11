<template>
  <img 
    v-b-popover.hover.bottomright="tweet.FullText" 
    v-bind:alt="twitterId" 
    v-bind:title="twitterId"
    style="width: 24px; height: 24px; margin: 2px" 
    v-lazy="tweet.ProfileImageUrl"
    />
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