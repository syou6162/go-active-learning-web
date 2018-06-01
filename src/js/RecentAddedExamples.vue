<template>
  <div>
    <select v-model="filter_label">
      <option value=1>Positive</option>
      <option value=-1>Negative</option>
    </select>
    <span>Selected: {{ filter_label }}</span>
    <b-card-group columns>
      <example 
        v-for="example in examplesFilterByLabel(filter_label)"
        v-bind:key="example.Url"
        v-bind:example="example"
        ></example>
    </b-card-group>
  </div>
</template>

<script>
import axios from 'axios';
import Example from './Example.vue';

export default {
  data () {
    return {
      filter_label: 1,
      results: []
    }
  },
  mounted() {
   axios.get("/api/recent_added_examples")
      .then(response => {
        this.results = response.data;
      });
  },
  methods: {
    examplesFilterByLabel: function(label) {
      return this.results.filter(function(e) {
        return e.Label == label;
      })
    }
  },
  components: {
    "example": Example
  }
}
</script>
