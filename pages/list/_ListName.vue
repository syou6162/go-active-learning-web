<template>
  <div>
    <b-form-group label="Please select a label">
      <b-form-radio-group
        buttons
        v-model="isNew"
        button-variant="outline-primary"
        :options="options" />
    </b-form-group>
    <b-card-group columns>
      <no-ssr>
      <example 
        v-for="example in examplesFilterByIsNew(isNew)"
        v-bind:key="example.Url"
        v-bind:example="example"
        v-bind:tweets="example.ReferringTweets"
        v-bind:isAdmin="isAdmin"
        ></example>
      </no-ssr>
    </b-card-group>
  </div>
</template>

<script>
import axios from 'axios';
import { Auth } from 'aws-amplify';
import Example from '~/components/Example.vue';
import { NewExample } from '~/assets/util';
import { IsAdmin } from '~/plugins/amplify.js';

const keywordsByListname = {
  "general": ["機械学習", "Machine Learning", "自然言語処理"],
  "article": ["機械学習", "Machine Learning", "自然言語処理"],
  "github": ["GitHub", "OSS", "Machine Learning"],
  "arxiv": ["arXiv", "Paper", "論文", "Machine Learning"],
  "slide": ["登壇", "発表", "Machine Learning", "機械学習"],
};

const descriptionByListname = {
  "general": "機械学習に関連する人気のエントリを読むことができます",
  "article": "機械学習に関連する人気のエントリを読むことができます",
  "github": "GitHub上で話題の機械学習に関連するリポジトリを見ることができます",
  "arxiv": "arXiv上で話題の機械学習に関連する論文を読むことができます",
  "slide": "SlideShareやSpeaker Deck上で話題の機械学習に関連する発表資料を読むことができます",
};

const isNewDayThresholdByListname = {
  "general": 2.5,
  "article": 3,
  "github": 10,
  "arxiv": 5,
  "slide": 10,
};

export default {
  layout: 'default',
  data () {
    return {
      title: "ML News",
      listname: 'general',
      examples: [],
      isNew: 0,
      options: [
        { text: 'All', value: 0 },
        { text: 'Recent', value: 1 },
      ],
      isAdmin: false,
    }
  },
  mounted() {
    Auth.currentAuthenticatedUser()
      .then(user => {
        this.isAdmin = true;
      })
      .catch(err => console.log(err))
  },
  async asyncData(context) {
    const listname = context.route.params.ListName;
    let data = await context.app.$axios.$get(`/api/examples?listName=${listname}`);
    const examples = data.Examples.map(e => {
      return NewExample(e, {
        "IsNewDayThreshold": isNewDayThresholdByListname[listname],
      })
    }).sort(function(a, b) {
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
      title: `ML News - ${listname}`,
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
        }
      ]
    };
  },
  methods: {
    examplesFilterByIsNew: function(isNew) {
      return this.examples.filter(function(e) {
        if (isNew == 0) {
          return true;
        } else {
          return e.IsNew == isNew;
        }
      })
    },
    keywords: function(listname) {
      return keywordsByListname[listname] || []; 
    },
    description: function(listname) {
      return descriptionByListname[listname] || ""; 
    },
    url: function(listname) {
      return "https://www.machine-learning.news/list/" + listname;
    },
  },
  components: {
    "example": Example,
  }
}
</script>
