<template>
  <b-col md=4>
    <b-card v-bind:href="example.Url">
      {{ example | getTitle(100, '...') }}
    </b-card>
  </b-col>
</template>

<script>
export default {
  props: ['example'],
  filters: {
    getTitle: function(example, length, omission) {
      var title = example.Title ? example.Title : example.Url;
      var length = length ? parseInt(length, 10) : 20;
      var ommision = omission ? omission.toString() : '...';

      if(title.length <= length) {
        return title;
      }
      else {
        return title.substring(0, length) + ommision;
      }
    }
  }
}
</script>
