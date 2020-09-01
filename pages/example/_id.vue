<template>
  <div class="mx-auto example">
    <b-card
      v-if="example"
      tag="article"
    >
      <h1 class="h4">
        {{ example | getTitle(1000, '...') }}
      </h1>
      <img
        v-if="example.OgImage"
        v-lazy="example.OgImage"
        :alt="example.Title"
        :title="example.Title" 
        class="img-thumbnail img-responsive ogimage"
        onerror="this.style.display='none'"
      >
      <p class="card-text">
        {{ example | getDescription(1000, '...') }}
      </p>
      <div v-if="example.HatenaBookmark.count > 0">
        <a
          :href="example.HatenaBookmark.entry_url"
          class="hatena-bookmark-link"
          target="_blank"
          rel="noopener"
        >{{ example.HatenaBookmark.count }} users</a>:
        <hatena-bookmark-icon 
          v-for="b in example.HatenaBookmark.bookmarks.slice(0, 9)"
          :key="b.user"
          :bookmark="b"
        />
      </div>
      <div v-if="example.ReferringTweets.Tweets.length > 0">
        <span class="tweets-count">{{ example.ReferringTweets.Count }} mentions</span>:
        <twitter-icon
          v-for="tweet in example.ReferringTweets.Tweets.slice(0, 8)"
          :key="tweet.IdStr"
          :tweet="tweet"
        />
      </div>
      <div v-if="keywords.length > 0">
        Keywords: 
        <b-button
          v-for="k in keywords" 
          :key="k"
          :href="'/search?query=' + k"
          variant="outline-primary"
          size="sm"
          class="keyword"
        >
          {{ k }}
        </b-button>
      </div>
      Date: {{ example.CreatedAt }}
      <annotate-buttons
        v-if="isAdmin"
        :example="example"
      />
      <b-card-footer>
        <img
          v-if="example.Favicon"
          v-lazy="example.Favicon"
          :alt="example.Title"
          class="example-favicon-img"
          onerror="this.style.display='none'"
        >
        <a
          :href="example.FinalUrl"
          target="_blank"
          rel="noopener"
        >{{ example | getDomain }} {{ example | getUserName }}</a>
        <b-button
          :href="tweetShareLink()"
          target="_blank"
          rel="noopener"
          class="float-right"
          variant="primary"
          size="sm"
        >
          Tweet
        </b-button>
      </b-card-footer>
    </b-card>
    <div
      v-if="tweetsWithPositiveLabelOrPositiveScore.length > 0"
      id="referring-tweet"
    >
      <h2 class="h4">
        Referring Tweets
      </h2>
      <b-list-group>
        <b-list-group-item
          v-for="t in tweetsWithPositiveLabelOrPositiveScore.slice(0, 9)"
          :key="t.IdStr"
          class="tweet-item"
        >
          <img
            v-if="t.ProfileImageUrl"
            :src="t.ProfileImageUrl"
            class="tweet-full-text-icon-img"
            onerror="this.src='/img/twitter_icon.png'"
          >
          <div class="tweet-screen-name-and-full-text-container">
            <a
              :href="'https://twitter.com/' + t.ScreenName + '/status/' + t.IdStr"
              target="_blank"
              rel="noopener"
            >@{{ t.ScreenName }}</a>
            <span v-html="fullTextWithLinks(t.FullText)" />
            <div class="tweet-footer">
              <span class="tweet-retweet-count">{{ t.RetweetCount }} RT</span>, <span class="tweet-favorite-count"> {{ t.FavoriteCount }} Fav</span>
              <span class="tweet-created-at">{{ t.CreatedAt }}</span>
            </div>
          </div>
          <tweet-annotate-buttons
            v-if="isAdmin"
            :tweet="t"
          />
        </b-list-group-item>
      </b-list-group>
    </div>
    <div
      v-if="hasBookmarksWithComment"
      id="bookmark-comments"
    >
      <h2 class="h4">
        Bookmark Comments
      </h2>
      <b-list-group>
        <b-list-group-item
          v-for="b in bookmarksWithComment.slice(0, 9)"
          :key="b.user"
          class="hatena-bookmark-item"
        >
          <img
            v-lazy="'https://cdn.profile-image.st-hatena.com/users/' + b.user+ '/profile.png'"
            class="hatena-bookmark-comment-user-icon-img"
          >
          <div class="hatena-bookmark-user-link-and-comment-container">
            <a 
              :href="'http://b.hatena.ne.jp/' + b.user"
              target="_blank"
              rel="noopener"
              class="hatena-bookmark-user-link"
            >
              id:{{ b.user }}
            </a>
            <span class="hatena-bookmark-comment">
              {{ b.comment }}
            </span>
            <span class="hatena-bookmark-timestamp">
              {{ b.timestamp }}
            </span>
          </div>
        </b-list-group-item>
      </b-list-group>
    </div>
    <h2
      v-if="similarExamples.length > 0"
      class="h4"
    >
      Related Entries
    </h2>
    <b-list-group>
      <b-list-group-item
        v-for="example in similarExamples"
        :key="example.FinalUrl"
        class="similar-example"
      >
        <img
          v-if="example.Favicon"
          :src="example.Favicon"
          class="similar-example-favicon-img"
          onerror="this.src='/img/website_icon.png'"
        >
        <img
          v-else
          src="/img/website_icon.png"
          class="similar-example-favicon-img"
        >
        <div class="similar-example-title-container">
          <b-button
            :href="example | getExampleUrl"
            class="float-right"
            size="sm"
          >
            Read more
          </b-button>
          {{ example | getTitle(100, '...') }}
          <div>
            <a
              :href="example.FinalUrl"
              target="_blank"
              rel="noopener"
            >{{ example | getDomain }} {{ example | getUserName }}</a>
          </div>
          <div>
            <span class="hatena-bookmark-link">{{ example.HatenaBookmark.count }} users</span>,
            <span class="tweets-count">{{ example.ReferringTweets.Count }} mentions</span>
            <span class="example-created-at">{{ example.CreatedAt }}</span>
          </div>
        </div>
      </b-list-group-item>
    </b-list-group>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'nuxt-property-decorator'
import { Auth } from 'aws-amplify';
import * as Sentry from '@sentry/browser';
import { Autolinker, AutolinkerConfig } from 'autolinker';
import { MetaInfo } from 'vue-meta'
import { Example, getTitle, getDescription, getDomain, getUserName, getExampleUrl, filterBookmarksWithComment } from '~/models/Example'
import Tweet from '~/models/Tweet';
import { NewExample } from '~/plugins/util';

import HatenaBookmarkIcon from '~/components/HatenaBookmarkIcon.vue'
import TwitterIcon from '~/components/TwitterIcon.vue'
import AnnotateButtons from '~/components/AnnotateButtons.vue'
import TweetAnnotateButtons from '~/components/TweetAnnotateButtons.vue'

@Component({
  components: {
    HatenaBookmarkIcon,
    TwitterIcon,
    AnnotateButtons,
    TweetAnnotateButtons,
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
  },
  asyncData({ app, params, error }) {
    return app.$axios.$get(`/api/example?id=${params.id}`, {timeout: 5000})
      .then((data) => {
        return {
          title: `[ML-News] ${data.Example.Title}`,
          example: NewExample(data.Example),
          similarExamples: data.SimilarExamples.filter(function(e) {
            return e.Label === 1 || e.Score > 0.0;
          }).map(function(e) {
            return NewExample(e)
          }),
          keywords: data.Keywords,
        }
      })
      .catch((err) => {
        let errObj = {};
        if (!err.response) {
          errObj["message"] = "Network error";
        } else {
          errObj["statusCode"] = err.response.status;
          errObj["message"] = err.response.data;
        }
        Sentry.captureException(err);
        return error(errObj);
      });
  }
})

export default class ExamplePage extends Vue {
  title: string = "ML-News"
  example: Example | null = null
  similarExamples: Example[] = []
  keywords: string[] = []
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
  tweetShareLink(): string {
    const txt = "👀";
    const hashtag = "ml_news";
    const url = this.example ? `https://www.machine-learning.news/example/${this.example.Id}` : "https://www.machine-learning.news";
    return `https://twitter.com/intent/tweet?text=${txt}&hashtags=${hashtag}&url=${url}`;
  }
  get hasBookmarksWithComment(): boolean {
    if (this.example === null) return false;
    return filterBookmarksWithComment(this.example).length > 0;
  }
  get bookmarksWithComment() {
    if (this.example === null) return [];
    return filterBookmarksWithComment(this.example);
  }
  tweetsWithPositiveLabelOrPositiveScore(): Tweet[] {
    if (this.example === null) return [];
    return this.example.ReferringTweets.Tweets.filter(function(t: Tweet) {
      return t.Label == 1 || t.Score > 0.0;
    }).sort(function(a: Tweet, b: Tweet) {
      if (a.Score > b.Score) {
        return -1;
      } else if (a.Score < b.Score) {
        return 1;
      } else {
        return 0;
      }
    });
  }
  head(): MetaInfo {
    const tweets = this.tweetsWithPositiveLabelOrPositiveScore().map(t => "@" + t.ScreenName + "「" + t.FullText + "」").slice(0, 3);
    const bookmarks = this.example ? filterBookmarksWithComment(this.example).map(b => "id:" + b.user + "「"+ b.comment + "」").slice(0, 3) : [];
    const description = tweets.join("\n") + bookmarks.join("\n");
    const robotsContent = this.example === null || this.example.Label == -1 ? "noindex, nofollow" : "index, follow";

    return {
      title: this.title,
      meta: [
        {
          name: "keywords",
          content: this.keywords.join(",")
        },
        {
          name: "description",
          content: description
        },
        {
          name: "og:title",
          content: this.title
        },
        {
          name: "og:type",
          content: "article"
        },
        {
          name: "og:description",
          content: description 
        },
        {
          name: "og:url",
          content: this.example ? `https://www.machine-learning.news/example/${this.example.Id}` : "https://www.machine-learning.news"
        },
        {
          name: "og:image",
          content: this.example ? this.example.OgImage : ""
        },
        {
          name: "robots",
          content: robotsContent
        }
      ],
      link: [
        {
          rel: "canonical",
          href: this.example ? `https://www.machine-learning.news/example/${this.example.Id}` : "https://www.machine-learning.news"
        }
      ]
    }
  }
}
</script>

<style scoped>
.example {
  max-width: 40rem;
}
.ogimage {
  width: 128px; 
  height: 96px; 
  margin: 3px; 
  float: right;
}
.example-favicon-img {
  width: 16px;
  height: 16px;
}
.tweet-item, .hatena-bookmark-item, .similar-example {
  padding: 10px 10px;
}
.hatena-bookmark-link {
  color: #ff4166;
}
.hatena-bookmark-comment-user-icon-img {
  float: left;
  width: 32px;
  height: 32px;
  margin: 0 10px 0 0;
}
.tweet-screen-name-and-full-text-container, .hatena-bookmark-user-link-and-comment-container, .similar-example-title-container {
  overflow: hidden;
}
.hatena-bookmark-user-link {
  margin: 0 0 0 0;
}
.hatena-bookmark-comment {
  margin: 0 4px 0 0;
}
.hatena-bookmark-timestamp {
  display: block;
  color: #999;
  float: right;
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
.example-created-at {
  color: #999;
  float: right;
  margin: 0 0 4px;
  line-height: 16px;
}
.keyword {
  padding: 3px 6px;
  margin: 1px 4px 1px 0;
}
.similar-example-favicon-img {
  float: left;
  width: 32px;
  height: 32px;
  margin: 0 10px 0 0;
}
</style>
