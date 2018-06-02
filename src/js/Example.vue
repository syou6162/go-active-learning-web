<template>
  <a v-bind:href="example.Url">
    <b-card v-bind:title="example | getTitle(100, '...')">
      <b-card-footer>
        {{ example | getDomain }}
      </b-card-footer>
    </b-card>
  </a>
</template>

<script>
export default {
  props: ['example'],
  filters: {
    getTitle: function(example, length, omission) {
      var title = example.Title ? example.Title : example.Url;
      var length = length ? parseInt(length, 10) : 20;
      var ommision = omission ? omission.toString() : '...';

      if (title.length <= length) {
        return title;
      }
      else {
        return title.substring(0, length) + ommision;
      }
    },
    getDomain: function(example) {
      var url = example.Url;
      return url.replace('http://','').replace('https://','').split(/[/?#]/)[0];
    }
  }
}
</script>
