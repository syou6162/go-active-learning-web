<template>
  <img 
    v-b-popover.hover="'id:' + bookmark.user" 
    v-bind:alt="bookmark.user" 
    style="width: 24px; height: 24px; margin: 2px" 
    v-lazy="'https://cdn.profile-image.st-hatena.com/users/' + bookmark.user+ '/profile.png'" 
    />
</template>

<script>
export default {
  props: ['bookmark'],
}
</script>
