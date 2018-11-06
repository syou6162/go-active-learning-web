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
      No search result for {{ query }}
    </div>
    <div v-else>
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
import NewExample from './util';

export default {
  data () {
    return {
      query: "",
      results: [],
      error: null,
      loading: true,
    }
  },
  watch: {
    query: function(newSearchQuery, oldSearchQuery) {
      this.debouncedGetSearchResult()
    },
  },
  created: function() {
    this.debouncedGetSearchResult = _.debounce(this.searchExamples, 500)
  },
  methods: {
    searchExamples: function() {
      let self = this;
      this.loading = true;
      this.error = null;

      let params = new URLSearchParams();
      params.append('query', this.query);
      axios.post("/api/search", params)
        .then(response => {
          this.results = response.data.map(e => NewExample(e));
          this.loading = false;
        }).catch(function (error) {
          if (error.response) {
            self.error = error.response.statusText;
          }
        });
    },
  },
  components: {
    "example": Example
  }
}
</script>
