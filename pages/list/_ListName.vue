<template>
  <div>
    <div>
      <a :href="rssLink(listname)">
        <img 
          src="/img/rss_icon.png"
          class="rss-icon"
        >
      </a>
      <b-form-group label="Please select a label">
        <b-form-radio-group
          v-model="filterType"
          buttons
          button-variant="outline-primary"
          :options="getFilterOptions(isAdmin)"
        />
      </b-form-group>
    </div>
    <b-card-group columns>
      <example 
        v-for="example in examplesFilterByFilterType(filterType)"
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

import Auth from '@aws-amplify/auth';
import { NewExample } from '~/plugins/util';

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
  "slide": "SlideShareやSpeaker DeckやGoogle Slides上で話題の機械学習に関連する発表資料を読むことができます",
  "video": "YouTube上で話題の機械学習に関連する発表を見ることができます",
  "event": "connpassやTECH PLAY上で話題の機械学習に関連する勉強会を探すことができます",
};

const isNewDayThresholdByListname: { [key: string]: number }= {
  "general": 1,
  "article": 1.5,
  "github": 5,
  "arxiv": 2.5,
  "slide": 5,
  "video": 5,
  "event": 5,
};

enum FilterType {
  All,
  Recent,
  Unlabeled
}

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

      var aTweetsCount = a.ReferringTweets.Count;
      var bTweetsCount = b.ReferringTweets.Count;
      
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
          href: this.rssLink(this.listname),
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
  filterType: FilterType = FilterType.All
  isAdmin: boolean = false

  mounted() {
    Auth.currentAuthenticatedUser()
      .then(user => {
        this.isAdmin = true;
      })
      .catch(err => console.log(err))
  }
  examplesFilterByFilterType(filterType: FilterType): Example[] {
    return this.examples.filter(function(e: Example) {
      switch(filterType) {
        case FilterType.All:
          return true
        case FilterType.Recent:
          return e.IsNew
        case FilterType.Unlabeled:
          return e.Label == 0;
      }
    })
  }
  getFilterOptions(isAdmin: boolean): { [key: string]: any }[] {
    let filterOptions: { [key: string]: any }[] = [
      { text: 'All', value: FilterType.All },
      { text: 'Recent', value: FilterType.Recent }
    ];
    if (isAdmin) {
      filterOptions.push({ text: 'Unlabeled', value: FilterType.Unlabeled })
    }
    return filterOptions;
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
  rssLink(listname: string): string {
    return "https://www.machine-learning.news/rss?listName=" + listname;
  }
}
</script>

<style>
.rss-icon {
  float: right;
  width: 32px;
  height: 32px;
  margin: 10px;
}
</style>
