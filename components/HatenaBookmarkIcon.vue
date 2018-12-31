<template>
  <img 
    v-b-popover.hover.bottomright="tags_and_comment" 
    v-bind:alt="user_id" 
    v-bind:title="user_id"
    style="width: 24px; height: 24px; margin: 2px" 
    v-lazy="'https://cdn.profile-image.st-hatena.com/users/' + bookmark.user+ '/profile.png'" 
    />
</template>

<script>
export default {
  computed: {
    user_id: function() {
      return 'id:' + this.bookmark.user;
    },
    tags_and_comment: function() {
      const tags = this.bookmark.tags.map(t => '[' + t + ']').join("")
      return tags + ' '+ this.bookmark.comment;
    },
  },
  props: ['bookmark'],
}
</script>
