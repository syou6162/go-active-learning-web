<template>
  <div>
    <input v-model="query" type="text" placeholder="Input search query here"></input>
    <b-card-group columns>
      <example 
        v-for="example in results"
        v-bind:key="example.Url"
        v-bind:example="example"
        ></example>
    </b-card-group>
  </div>
</template>

<script>
import axios from 'axios';
import _ from 'lodash';
import Example from './Example.vue';

export default {
  data () {
    return {
      query: "",
      results: []
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
      let params = new URLSearchParams();
      params.append('query', this.query);
      axios.post("/api/search", params)
        .then(response => {
          this.results = response.data;
        });
    },
  },
  components: {
    "example": Example
  }
}
</script>
