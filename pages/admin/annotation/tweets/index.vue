<template>
  <div>
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
        v-bind:key="getKey(example, label)"
        v-bind:example="example"
        v-bind:tweet="example.ReferringTweets[0]"
        v-bind:isAdmin="isAdmin"
        ></example>
    </b-card-group>
  </div>
</template>

<script>
import axios from 'axios';
import AdminAnnotateTweet from '~/components/AdminAnnotateTweet.vue';
import TweetAnnotateButtons from '~/components/TweetAnnotateButtons.vue';
import { NewExample } from '~/assets/util';
import { Auth } from 'aws-amplify';

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
      isAdmin: false,
    }
  },
  async asyncData(context) {
    let data = await context.app.$axios.$get("/api/recent_added_tweets");
    return {
      positive: data.PositiveExamples.map(e => NewExample(e)),
      negative: data.NegativeExamples.map(e => NewExample(e)),
      unlabeled: data.UnlabeledExamples.map(e => NewExample(e))
    }
  },
  mounted() {
    Auth.currentAuthenticatedUser()
      .then(user => {
        this.isAdmin = true;
      })
      .catch(err => console.log(err))
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
    },
    getKey: function(example, label) {
      return String(label) + ":" + example.Url + ":" + example.ReferringTweets[0].IdStr;
    },
  },
  head() {
    return {
      title: "最近追加されたTweet一覧",
    };
  },
  components: {
    "example": AdminAnnotateTweet,
    "tweet-annotate-buttons": TweetAnnotateButtons
  }
}
</script>
