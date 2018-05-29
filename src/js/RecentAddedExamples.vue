<template>
  <div>
    <select v-model="filter_label">
      <option value=1>Positive</option>
      <option value=-1>Negative</option>
    </select>
    <span>Selected: {{ filter_label }}</span>
    <ul>
      <list-example 
        v-for="example in examplesFilterByLabel(filter_label)"
        v-bind:key="example.Url"
        v-bind:example="example"
        ></list-example>
    </ul>
  </div>
</template>

<script>
import ListExample from './ListExample.vue';

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
    "list-example": ListExample
  }
}
</script>
