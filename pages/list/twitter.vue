<template>
  <div>
    <b-card-group columns>
      <div
        v-for="example in examples"
        :key="example.Id"
      >
        <div
          v-for="tweet in example.ReferringTweets.Tweets"
          :key="tweet.IdStr"
          :tweet="tweet"
          :isAdmin="isAdmin"
        >
          <b-card no-body>
            <b-card-body
              :title="example | getTitle(75, '...')"
              title-tag="h5"
              class="m-1 p-2"
            >
              <img
                v-if="tweet.ProfileImageUrl"
                v-lazy="tweet.ProfileImageUrl"
                class="tweet-full-text-icon-img"
                onerror="this.src='/img/twitter_icon.png'"
              >
              <div class="tweet-screen-name-and-full-text-container">
                <a 
                  :href="'https://twitter.com/' + tweet.ScreenName + '/status/' + tweet.IdStr"
                  target="_blank"
                  rel="noopener"
                >
                  @{{ tweet.ScreenName }}
                </a>
                <span class="tweet-full-text" v-html="fullTextWithLinks(tweet.FullText)" />
                <div class="tweet-footer">
                  <span class="tweet-retweet-count">{{ tweet.RetweetCount }} RT</span>, <span class="tweet-favorite-count"> {{ tweet.FavoriteCount }} Fav</span>
                  <span class="tweet-created-at">{{ tweet.CreatedAt }}</span>
                </div>
              </div>
              <tweet-annotate-buttons
                v-if="isAdmin"
                :tweet="tweet"
              />
              <b-card-footer>
                <b-button
                  :to="example | getExampleUrl"
                  class="float-right"
                  size="sm"
                >
                  Read more
                </b-button>
                <img
                  v-if="example.Favicon"
                  v-lazy="example.Favicon"
                  class="example-favicon-img"
                  onerror="this.style.display='none'"
                >
                <a 
                  :href="example.FinalUrl"
                  target="_blank"
                  rel="noopener"
                >
                  {{ example | getDomain }} {{ example | getUserName }}
                </a>
              </b-card-footer>
            </b-card-body>
          </b-card>
        </div>
      </div>
    </b-card-group>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'nuxt-property-decorator';
import { Autolinker, AutolinkerConfig } from 'autolinker';
import { MetaInfo } from 'vue-meta'
import { Example, getTitle, getDescription, getDomain, getUserName, getExampleUrl } from '~/models/Example'

import { Auth } from 'aws-amplify';
import { NewExample } from '~/plugins/util';
import TweetAnnotateButtons from '~/components/TweetAnnotateButtons.vue'

@Component({
  components: {
    TweetAnnotateButtons,
  },
  async asyncData(context) {
    let data = await context.app.$axios.$get(`/api/tweets`);
    const examples = data.Examples.map(function(e) {
      return NewExample(e)
    });
    return {
      title: `ML-News - Twitter`,
      examples: examples,
      loading: false
    };
  },
  filters: {
    getTitle(example: Example, length: number, omission: string): string {
      return getTitle(example, length, omission)
    },
    getDescription(example: Example, length: number, omission: string): string {
      return getDescription(example, length, omission)
    },
    getDomain(example: Example): string {
      return getDomain(example)
    },
    getUserName(example: Example): string {
      return getUserName(example)
    },
    getExampleUrl(example: Example): string {
      return getExampleUrl(example)
    },
  }  
})

export default class ListNamePage extends Vue {
  title: string = "ML-News"
  examples: Example[] = []
  isAdmin: boolean = false

  mounted() {
    Auth.currentAuthenticatedUser()
      .then(user => {
        this.isAdmin = true;
      })
      .catch(err => console.log(err))
  }
  fullTextWithLinks(fullText: string): string {
    const opts: AutolinkerConfig = { 
      mention: 'twitter',
      hashtag: 'twitter'
      };
    return Autolinker.link(fullText, opts);
  }
  head(): MetaInfo {
    return {
      title: this.title,
      meta: [
        {
          name: "description",
          content: `機械学習に関するTweetを見れます`,
        }
      ],
      link: [
        {
          rel: "canonical",
          href: `https://www.machine-learning.news/list/twitter`
        },
      ]
    };
  }
}
</script>

<style scoped>
.example-favicon-img {
  width: 16px;
  height: 16px;
}
.tweet-screen-name-and-full-text-container {
  overflow: hidden;
}
.tweets-count, .tweet-retweet-count, .tweet-favorite-count {
  color: #ff4166;
}
.tweet-full-text-icon-img {
  float: left;
  width: 32px;
  height: 32px;
  margin: 0 10px 0 0;
}
.tweet-created-at {
  color: #999;
  float: right;
}
.tweet-full-text {
  font-size: 16px;
  margin: 0 0 4px;
  color: #55606a;
}
</style>
