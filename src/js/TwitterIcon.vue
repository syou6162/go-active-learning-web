<template>
  <img 
    v-b-popover.hover.bottomright="tweet.FullText" 
    v-bind:alt="twitterId" 
    v-bind:title="twitterId"
    style="width: 24px; height: 24px; margin: 2px" 
    v-lazy="tweet.ProfileImageUrl"
    />
</template>

<script>
import { getDomain, getTwitterId, getTweetTitle } from './util';
export default {
  computed: {
    twitterId: function() {
      return '@' + this.tweet.ScreenName;
    },
  },
  props: ['tweet'],
}
</script>
