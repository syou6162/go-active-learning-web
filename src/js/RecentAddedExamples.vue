<template>
  <div>
    <div v-if="loading">
      Now loading...
    </div>
    <div v-else-if="error">
      Fail to retrieve from API server. Error: {{ error }}
    </div>
    <div v-else>
      <b-form-group label="Please select a label">
        <b-form-radio-group
          buttons
          v-model="label"
          button-variant="outline-primary"
          :options="options" />
      </b-form-group>
      <b-card-group columns>
        <example 
          v-for="example in searchExamplesByLabel(label)"
          v-bind:key="example.Url"
          v-bind:example="example"
          v-bind:isAdmin="isAdmin"
          ></example>
      </b-card-group>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Example from './Example.vue';
import { NewExample, IsAdmin } from './util';

export default {
  data () {
    return {
      label: 0,
      options: [
        { text: 'Unlabeled', value: 0 },
        { text: 'Positive', value: 1 },
        { text: 'Negative', value: -1 },
      ],
      positive: [],
      negative: [],
      unlabeled: [],
      error: null,
      loading: true,
      isAdmin: false,
    }
  },
  mounted() {
    let self = this;
    this.loading = true;
    this.error = null;
    
    axios.get("/api/recent_added_examples")
      .then(response => {
        this.positive = response.data.PositiveExamples.map(e => NewExample(e));
        this.negative = response.data.NegativeExamples.map(e => NewExample(e));
        this.unlabeled = response.data.UnlabeledExamples.map(e => NewExample(e));
        this.loading = false;
      }).catch(function (error) {
        if (error.response) {
          self.loading = false;
          self.error = error.response.statusText;
        }
      });
    this.isAdmin = IsAdmin();
  },
  methods: {
    searchExamplesByLabel: function(label) {
      if (label == 1) {
        return this.positive;
      } else if (label == -1) {
        return this.negative;
      } else if (label == 0) {
        return this.unlabeled;
      } else {
        return [];
      }
    }
  },
  components: {
    "example": Example
  }
}
</script>
