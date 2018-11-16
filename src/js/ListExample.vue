<template>
  <div>
    <vue-headful 
      v-bind:title="title" 
      v-bind:keywords="keywords(listname).join(',')" 
      v-bind:description="description(listname)" 
      v-bind:url="url(listname)"
      />
    <div v-if="loading">
      Now loading...
    </div>
    <div v-else-if="error">
      Fail to retrieve from API server. Error: {{ error }}
    </div>
    <div v-else>
      <b-form-group label="Please select a label">
        <b-form-radio-group
          buttons
          v-model="isNew"
          button-variant="outline-primary"
          :options="options" />
      </b-form-group>
      <b-card-group columns>
        <example 
          v-for="example in examplesFilterByIsNew(isNew)"
          v-bind:key="example.Url"
          v-bind:example="example"
          ></example>
      </b-card-group>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Example from './Example.vue';
import NewExample from './util';

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
      error: null,
      loading: true,
    }
  },
  mounted() {
    this.fetchList(this.$route.params.listname)
  },
  watch: {
    '$route' (to, from) {
      this.fetchList(to.params.listname)
    }
  },
  methods: {
    fetchList(listname) {
      let self = this;
      this.loading = true;
      this.error = null;

      axios.get("/api/examples?listName=" + listname)
        .then(response => {
          this.examples = response.data.Examples.map(e => 
            NewExample(e, {
              "IsNewDayThreshold": isNewDayThresholdByListname[listname],
            })
          ).sort(function(a, b) {
            if (a.HatenaBookmarks.count > b.HatenaBookmarks.count) return -1;
            if (a.HatenaBookmarks.count < b.HatenaBookmarks.count) return 1;
            return 0;
          });
          this.listname = this.$route.params.listname;
          this.title = `ML News - ${this.listname}`;
          this.loading = false;
        }).catch(function (error) {
          if (error.response) {
            self.loading = false;
            self.error = error.response.statusText;
          }
        });
    },
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
