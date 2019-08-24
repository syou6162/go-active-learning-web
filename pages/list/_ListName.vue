<template>
  <div>
    <b-form-group label="Please select a label">
      <b-form-radio-group
        v-model="isNew"
        buttons
        button-variant="outline-primary"
        :options="options"
      />
    </b-form-group>
    <b-card-group columns>
      <example 
        v-for="example in examplesFilterByIsNew(isNew)"
        :key="example.Url"
        :example="example"
        :tweets="example.ReferringTweets"
        :is-admin="isAdmin"
      />
    </b-card-group>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Example from '~/models/Example'

import { Auth } from 'aws-amplify';
import { NewExample } from '~/assets/util';

const keywordsByListname: { [key: string]: string[] } = {
  "general": ["機械学習", "Machine Learning", "自然言語処理"],
  "article": ["機械学習", "Machine Learning", "自然言語処理"],
  "github": ["GitHub", "OSS", "Machine Learning"],
  "arxiv": ["arXiv", "Paper", "論文", "Machine Learning"],
  "slide": ["登壇", "発表", "Machine Learning", "機械学習"],
  "video": ["登壇", "発表", "Machine Learning", "機械学習"],
  "event": ["勉強会", "Machine Learning", "機械学習"],
};

const descriptionByListname: { [key: string]: string } = {
  "general": "機械学習に関連する人気のエントリを読むことができます",
  "article": "機械学習に関連する人気のエントリを読むことができます",
  "github": "GitHub上で話題の機械学習に関連するリポジトリを見ることができます",
  "arxiv": "arXiv上で話題の機械学習に関連する論文を読むことができます",
  "slide": "SlideShareやSpeaker Deck上で話題の機械学習に関連する発表資料を読むことができます",
  "video": "YouTube上で話題の機械学習に関連する発表を見ることができます",
  "event": "connpassやTECH PLAY上で話題の機械学習に関連する勉強会を探すことができます",
};

const isNewDayThresholdByListname: { [key: string]: number }= {
  "general": 2.5,
  "article": 3,
  "github": 10,
  "arxiv": 5,
  "slide": 10,
  "video": 10,
  "event": 5,
};

@Component({
  components: {
    Example: () => import('~/components/Example.vue')
  },
  async asyncData(context) {
    const listname = context.route.params.ListName;
    let data = await context.app.$axios.$get(`/api/examples?listName=${listname}`);
    const examples = data.Examples.map((e: Example) => {
      return NewExample(e, {
        "IsNewDayThreshold": isNewDayThresholdByListname[listname],
      })
    }).sort(function(a: Example, b: Example) {
      var aHatebuCount = a.HatenaBookmark.count;
      var bHatebuCount = b.HatenaBookmark.count;

      var aTweetsCount = 0;
      var bTweetsCount = 0;
      if (a.ReferringTweets) {
        aTweetsCount = a.ReferringTweets.length;
      }
      if (b.ReferringTweets) {
        bTweetsCount = b.ReferringTweets.length;
      }

      if (aHatebuCount + aTweetsCount > bHatebuCount + bTweetsCount) {
        return -1;
      } else if (aHatebuCount +aTweetsCount < bHatebuCount + bTweetsCount) {
        return 1;
      } else {
        return 0;
      }
    });
    return {
      title: `ML-News - ${listname}`,
      listname: listname,
      examples: examples,
      loading: false
    };
  },
  head() {
    return {
      title: this.title,
      meta: [
        {
          name: "keywords",
          content: this.keywords(this.listname).join(",")
        },
        {
          name: "description",
          content: this.description(this.listname),
        }
      ],
      link: [
        {
          rel: "canonical",
          href: `https://www.machine-learning.news/list/${this.listname}`
        },
        {
          rel: "alternate",
          type: "application/atom+xml",
          title: `ML-News - ${this.listname}`,
          href: `https://www.machine-learning.news/rss?listName=${this.listname}`
        }
      ]
    };
  }
})

export default class ListNamePage extends Vue {
  layout: string = 'default'

  title: string = "ML-News"
  listname: string = 'general'
  examples: Example[] = []
  isNew: boolean = false
  options: { [key: string]: any }[] = [
    { text: 'All', value: false },
    { text: 'Recent', value: true }
  ]
  isAdmin: boolean = false

  mounted() {
    Auth.currentAuthenticatedUser()
      .then(user => {
        this.isAdmin = true;
      })
      .catch(err => console.log(err))
  }
  examplesFilterByIsNew(isNew: boolean): Example[] {
    return this.examples.filter(function(e: Example) {
      if (!isNew) {
        return true;
      } else {
        return e.IsNew == isNew;
      }
    })
  }
  keywords(listname: string): string[] {
    return keywordsByListname[listname] || []; 
  }
  description(listname: string): string {
    return descriptionByListname[listname] || ""; 
  }
  url(listname: string): string {
    return "https://www.machine-learning.news/list/" + listname;
  }
}
</script>
