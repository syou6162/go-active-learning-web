<template>
  <div class="mx-auto" style="max-width: 40rem;">
    <b-card v-if="example" v-bind:title="example | getTitle(1000, '...')" tag="article">
      <img v-if="example.OgImage" class="img-thumbnail img-responsive" style="width: 128px; height: 96px; margin: 3px; float: right;" v-lazy="example.OgImage" onerror="this.style.display='none'" />
      <p class="card-text">
        {{ example | getDescription(1000, '...') }}
      </p>
      <b-card-footer>
        <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
        <a v-bind:href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
      </b-card-footer>
    </b-card>
    <h4 v-if="similarExamples.length > 0">Related Entries</h4>
    <b-list-group>
      <b-list-group-item v-for="example in similarExamples">
        <b-button v-bind:href="'/example/' + encodeURIComponent(example.FinalUrl)" class="float-right" size="sm">Read more</b-button>
        <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
        {{ example | getTitle(100, '...') }}
      </b-list-group-item>
    </b-list-group>
  </div>
</template>

<script>
import axios from 'axios';
import Example from './Example.vue';
import NewExample from './util';

export default {
  data () {
    return {
      url: this.$route.params.url,
      example: null,
      similarExamples: []
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
        this.example = NewExample(response.data.Example);
        this.similarExamples = response.data.SimilarExamples.filter(function(e) {
          return e.Score > 0.0;
        });
      });
    }
  },
  components: {
    "example": Example
  }
}
</script>
