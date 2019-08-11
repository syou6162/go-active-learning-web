<template>
  <img 
    v-b-popover.hover.bottomright="tags_and_comment" 
    v-bind:alt="user_id" 
    v-bind:title="user_id"
    style="width: 24px; height: 24px; margin: 2px" 
    v-lazy="'https://cdn.profile-image.st-hatena.com/users/' + bookmark.user+ '/profile.png'" 
    />
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Bookmark from '~/models/Bookmark'

@Component
export default class HatenaBookmarkIcon extends Vue {
  get user_id(): string {
    return 'id:' + this.bookmark.user;
  }
  get tags_and_comment(): string {
    const tags = this.bookmark.tags.map(t => '[' + t + ']').join("")
    return tags + ' '+ this.bookmark.comment;
  }
  @Prop()
  bookmark!: Bookmark;
}
</script>