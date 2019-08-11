<template>
  <div>
    Annotate as: 
    <b-button size="sm" v-bind:variant="example.Label == 1 ? 'primary' : 'outline-primary'" :pressed="false" v-on:click="updateLabel(example, 1)">Positive</b-button>
    <b-button size="sm" v-bind:variant="example.Label == -1 ? 'danger' : 'outline-danger'" v-on:click="updateLabel(example, -1)">Negative</b-button>
    <b-button size="sm" v-bind:variant="example.Label == 0 ? 'secondary' : 'outline-secondary'" v-on:click="updateLabel(example, 0)">Unlabeled</b-button>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import axios from 'axios';
import Example from '~/models/Example'

@Component
export default class AnnotateButtons extends Vue {
  @Prop({required: true})
  example!: Example 
  
  updateLabel(example: Example, label: number) {
    let idToken = localStorage.getItem("CognitoIdentityServiceProvider.4ia5ifrn456rqg4vr6dqfh7e68.yasuhisa.idToken");
    let headers = { headers: { 'Authorization': idToken } };
    axios.post(
      "https://3ojd2wnlpg.execute-api.us-east-1.amazonaws.com/Prod/update_example_label", 
      {
        id: example.Id,
        label: label,
      },
      headers
    ).then(response => {
      example.Label = response.data.Label;
    }).catch(function (error) {
      alert(`${error}: Failed to annotate "${example.Title}"`);
    })
  }
}
</script>