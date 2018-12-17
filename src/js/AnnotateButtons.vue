<template>
  <div>
    Annotate as: 
    <b-button size="sm" variant="outline-primary" :pressed="false" v-on:click="updateLabel(example, 1)">Positive</b-button>
    <b-button size="sm" variant="outline-danger" v-on:click="updateLabel(example, -1)">Negative</b-button>
    <b-button size="sm" variant="outline-secondary" v-on:click="updateLabel(example, 0)">Unlabeled</b-button>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  props: ['example'],
  methods: {
    updateLabel(example, label) {
      axios.post(
        "https://3ojd2wnlpg.execute-api.us-east-1.amazonaws.com/Prod/update_example_label", 
        {
          url: example.Url,
          label: label,
        },
      ).then(response => {
        example.Label = label;
      }).catch(function (error) {
        alert(`Failed to annotate "${example.Title}"`);
      })
    },
  }
}
</script>
