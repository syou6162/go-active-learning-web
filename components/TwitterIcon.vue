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
import { Component, Vue } from 'vue-property-decorator';

export default Vue.extend({
  computed: {
    twitterId(): string {
      return '@' + this.tweet.ScreenName;
    }
  },
  props: ['tweet']
});
</script>
