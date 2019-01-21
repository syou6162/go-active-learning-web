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
        v-bind:key="example.Url"
        v-bind:example="example"
        v-bind:isAdmin="isAdmin"
        ></example>
    </b-card-group>
  </div>
</template>

<script>
import axios from 'axios';
import Example from '~/components/Example.vue';
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
    let data = await context.app.$axios.$get("/api/recent_added_examples");
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
    }
  },
  head() {
    return {
      title: "最近追加された事例一覧",
      meta: [
        {
          name: "keywords",
          context: "最近追加された事例一覧"
        }
      ],
      link: [
        {
          rel: "canonical",
          href: "https://www.machine-learning.news/recent-added-examples"
        }
      ]
    };
  },
  components: {
    "example": Example
  }
}
</script>
