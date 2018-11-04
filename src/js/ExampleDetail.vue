<template>
  <b-card
    v-if="example"
    v-bind:title="example | getTitle(1000, '...')" 
    tag="article"
    class="mx-auto"
    style="max-width: 40rem;">
    <img v-if="example.OgImage" class="img-thumbnail img-responsive" style="width: 128px; height: 96px; margin: 3px; float: right;" v-lazy="example.OgImage" onerror="this.style.display='none'" />
    <p class="card-text">
      {{ example | getDescription(1000, '...') }}
    </p>
    <b-card-footer>
      <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
      <a v-bind:href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
    </b-card-footer>
  </b-card>
</template>

<script>
import axios from 'axios';
import NewExample from './util';

export default {
  data () {
    return {
      url: this.$route.params.url,
      example: null
    }
  },
  mounted() {
    this.fetchExample(this.$route.params.url)
  },
  methods: {
    fetchExample(url) {
      this.examples = [];
      axios.get("/api/example?url=" + encodeURIComponent(url))
      .then(response => {
        this.example = NewExample(response.data);
      });
    }
  },
}
</script>
