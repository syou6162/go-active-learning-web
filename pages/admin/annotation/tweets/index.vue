<template>
  <div>
    <b-form-group label="Please select a label">
      <b-form-radio-group
        v-model="label"
        buttons
        button-variant="outline-primary"
        :options="options"
      />
    </b-form-group>
    <b-card-group columns>
      <admin-annotate-tweets
        v-for="example in searchExamplesByLabel(label)"
        :key="getKey(example, label)"
        :example="example"
        :tweets="example.ReferringTweets.Tweets"
        :is-admin="isAdmin"
      />
    </b-card-group>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Example from '~/models/Example'
import { NewExample } from '~/plugins/util';
import { Auth } from 'aws-amplify';

@Component({
  components: {
    AdminAnnotateTweets: () => import('~/components/AdminAnnotateTweets.vue')
  },
  async asyncData(context) {
    let data = await context.app.$axios.$get("/api/recent_added_tweets");
    return {
      positive: data.PositiveExamples.map(e => NewExample(e)),
      negative: data.NegativeExamples.map(e => NewExample(e)),
      unlabeled: data.UnlabeledExamples.map(e => NewExample(e))
    }
  },
  head() {
    return {
      title: "最近追加されたTweet一覧",
    };
  }
})

export default class AdminAnnotationTweet extends Vue {
  label: number = 0
  options: { [key: string]: any}[] = [
    { text: 'Unlabeled', value: 0 },
    { text: 'Positive', value: 1 },
    { text: 'Negative', value: -1 },
  ]
  positive: Example[] = []
  negative: Example[] = []
  unlabeled: Example[] = []
  isAdmin: boolean = false
  mounted() {
    Auth.currentAuthenticatedUser()
      .then(user => {
        this.isAdmin = true;
      })
      .catch(err => console.log(err))
  }
  searchExamplesByLabel(label: number) {
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
  getKey(example: Example, label: number) {
    return String(label) + ":" + example.Url + ":" + example.ReferringTweets.Tweets[0].IdStr;
  }
}
</script>
