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
      <example 
        v-for="example in searchExamplesByLabel(label)"
        :key="example.Url"
        :example="example"
        :tweets="example.ReferringTweets"
        :is-admin="isAdmin"
      />
    </b-card-group>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'nuxt-property-decorator';
import { MetaInfo } from 'vue-meta'
import Example from '~/components/Example.vue';
import { NewExample } from '~/plugins/util';
import { Auth } from 'aws-amplify';

@Component({
  components: {
    Example,
  },
  async asyncData(context) {
    let data = await context.app.$axios.$get("/api/recent_added_examples");
    return {
      positive: data.PositiveExamples.map(e => NewExample(e)),
      negative: data.NegativeExamples.map(e => NewExample(e)),
      unlabeled: data.UnlabeledExamples.map(e => NewExample(e))
    }
  },
})

export default class AdminAnnotationExample extends Vue {
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
  searchExamplesByLabel(label: number): Example[] {
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
  head(): MetaInfo {
    return {
      title: "最近追加された事例一覧",
      meta: [
        {
          name: "keywords",
          content: "最近追加された事例一覧"
        }
      ],
      link: [
        {
          rel: "canonical",
          href: "https://www.machine-learning.news/recent-added-examples"
        }
      ]
    };
  }
}
</script>
