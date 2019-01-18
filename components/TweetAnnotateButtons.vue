<template>
  <div>
    Annotate as: 
    <b-button size="sm" v-bind:variant="tweet.Label == 1 ? 'primary' : 'outline-primary'" :pressed="false" v-on:click="updateLabel(tweet, 1)">Positive</b-button>
    <b-button size="sm" v-bind:variant="tweet.Label == -1 ? 'danger' : 'outline-danger'" v-on:click="updateLabel(tweet, -1)">Negative</b-button>
    <b-button size="sm" v-bind:variant="tweet.Label == 0 ? 'secondary' : 'outline-secondary'" v-on:click="updateLabel(tweet, 0)">Unlabeled</b-button>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  props: ['tweet'],
  methods: {
    updateLabel(tweet, label) {
      let idToken = localStorage.getItem("CognitoIdentityServiceProvider.4ia5ifrn456rqg4vr6dqfh7e68.yasuhisa.idToken");
      let headers = { headers: { 'Authorization': idToken } };
      axios.post(
        "https://3ojd2wnlpg.execute-api.us-east-1.amazonaws.com/Prod/update_tweet_label", 
        {
          example_id: tweet.ExampleId,
          id_str: tweet.IdStr,
          label: label,
        },
        headers
      ).then(response => {
        tweet.Label = label;
      }).catch(function (error) {
        alert(`${error}: Failed to annotate "${tweet.ScreenName}"`);
      })
    },
  }
}
</script>
