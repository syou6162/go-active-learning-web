<template>
  <div>
    <b-form-group label="Please select a label">
      <b-form-radio-group
        buttons
        v-model="filter_label"
        button-variant="outline-primary"
        :options="options" />
    </b-form-group>
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
import NewExample from './util';

export default {
  data () {
    return {
      filter_label: 1,
      options: [
        { text: 'Positive', value: 1 },
        { text: 'Negative', value: -1 },
        { text: 'Unlabeled', value: 0 },
      ],
      results: []
    }
  },
  mounted() {
   axios.get("/api/recent_added_examples")
      .then(response => {
        this.results = response.data.map(e => NewExample(e));
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
