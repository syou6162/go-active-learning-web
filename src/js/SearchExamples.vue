<template>
  <div>
    <input v-model="query" type="text" placeholder="Input search query here"></input>
    <div v-if="loading && query">
      Now loading...
    </div>
    <div v-else-if="error">
      Fail to retrieve from API server. Error: {{ error }}
    </div>
    <div v-else-if="query && results.length == 0">
      No search result for '{{ query }}'
    </div>
    <div v-else>
      <vue-headful
         v-bind:title="getSearchResultsTitle"
         v-bind:description="getSearchResultsDescription"
         v-bind:url="getSearchResultUrl"
        />
      <b-card-group columns>
        <example 
          v-for="example in results"
          v-bind:key="example.Url"
          v-bind:example="example"
          ></example>
      </b-card-group>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import _ from 'lodash';
import Example from './Example.vue';
import { NewExample } from './util';

export default {
  mounted() {
    this.query = this.$route.query.query;
    this.debouncedGetSearchResult();
  },
  data () {
    return {
      query: this.$route.query.query,
      results: [],
      error: null,
      loading: true,
    }
  },
  watch: {
    query: function(newSearchQuery, oldSearchQuery) {
      // https://router.vuejs.org/ja/guide/essentials/navigation.html
      this.$router.push({ query: { query: newSearchQuery }})
      this.debouncedGetSearchResult();
    },
  },
  created: function() {
    this.debouncedGetSearchResult = _.debounce(this.searchExamples, 500)
  },
  methods: {
    searchExamples: function() {
      let self = this;
      // undefinedのときは検索しない
      if (this.query === void 0) {
        return;
      }

      this.loading = true;
      this.error = null;

      let params = new URLSearchParams();
      params.append('query', this.query);
      axios.post("/api/search", params)
        .then(response => {
          this.results = response.data.Examples.map(e => NewExample(e));
          this.loading = false;
        }).catch(function (error) {
          if (error.response) {
            self.loading = false;
            self.error = error.response.statusText;
          }
        });
    },
  },
  computed: {
    getSearchResultsTitle() {
      return "ML News - 「" + this.query + "」に関する検索結果";
    },
    getSearchResultsDescription() {
      return this.results.map(e => e.Title).join("\n");
    },
    getSearchResultUrl: function() {
      return "https://www.machine-learning.news/search?query=" + this.query;
    },
  },
  components: {
    "example": Example
  }
}
</script>
