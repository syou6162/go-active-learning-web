<template>
  <div class="mx-auto example">
    <b-card v-if="example" tag="article">
      <h1 class="h4">{{ example | getTitle(1000, '...') }}</h1>
      <img
        v-if="example.OgImage"
        v-lazy="example.OgImage"
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
        >{{ example.HatenaBookmark.count }} users</a>:
        <hatena-bookmark-icon 
          v-for="b in example.HatenaBookmark.bookmarks.slice(0, 9)"
          :key="b.user"
          :bookmark="b"
        />
      </div>
      <div v-if="example.ReferringTweets && example.ReferringTweets.length > 0">
        <span class="tweets-count">{{ example.ReferringTweets.length }} mentions</span>:
        <twitter-icon
          v-for="tweet in example.ReferringTweets.slice(0, 8)"
          :key="tweet.ScreenName"
          :tweet="tweet"
        />
      </div>
      <div v-if="keywords.length > 0">
        Keywords: 
        <b-link 
          v-for="k in keywords" 
          :key="k"
          :href="'/search?query=' + k"
        >
          {{ k }}
        </b-link>
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
          class="example-favicon-img"
          onerror="this.style.display='none'"
        >
        <a :href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
      </b-card-footer>
    </b-card>
    <div v-if="example.ReferringTweets && example.ReferringTweets.length > 0">
      <h4>Referring Tweets</h4>
      <b-list-group>
        <b-list-group-item
          v-for="t in tweetsWithPositiveLabelOrPositiveScore.slice(0, 9)"
          :key="t.ScreenName"
        >
          <img
            v-if="t.ProfileImageUrl"
            v-lazy="t.ProfileImageUrl"
            class="tweet-icon-img"
            onerror="this.style.display='none'"
          >
          <a :href="'https://twitter.com/' + t.ScreenName + '/status/' + t.IdStr">@{{ t.ScreenName }}</a>
          {{ t.FullText }}
          <div class="tweet-footer">
            <span class="tweet-retweet-count">{{t.RetweetCount}} RT</span>, <span class="tweet-favorite-count"> {{t.FavoriteCount}} Fav</span>
            <span class="tweet-created-at">{{t.CreatedAt}}</span>
          </div>
          <tweet-annotate-buttons
            v-if="isAdmin"
            :tweet="t"
          />
        </b-list-group-item>
      </b-list-group>
    </div>
    <div v-if="hasBookmarksWithComment">
      <h4>Bookmark Comments</h4>
      <b-list-group>
        <b-list-group-item
          v-for="b in bookmarksWithComment.slice(0, 9)"
          :key="b.user"
        >
          <img
            v-lazy="'https://cdn.profile-image.st-hatena.com/users/' + b.user+ '/profile.png'"
            class="hatena-bookmark-user-icon-img"
          >
          <a :href="'http://b.hatena.ne.jp/' + b.user">id:{{ b.user }}</a>
          {{ b.comment }}
        </b-list-group-item>
      </b-list-group>
    </div>
    <h4 v-if="similarExamples.length > 0">
      Related Entries
    </h4>
    <b-list-group>
      <b-list-group-item
        v-for="example in similarExamples"
        :key="example.FinalUrl"
      >
        <b-button
          :href="example | getExampleUrl"
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
        {{ example | getTitle(100, '...') }}
      </b-list-group-item>
    </b-list-group>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Auth } from 'aws-amplify';
import Example from '~/components/Example.vue';
import Tweet from '~/components/Tweet.vue';
import { NewExample, filterBookmarksWithComment } from '~/assets/util';

@Component({
  components: {
    HatenaBookmarkIcon: () => import('~/components/HatenaBookmarkIcon.vue'),
    TwitterIcon: () => import('~/components/TwitterIcon.vue'),
    AnnotateButtons: () => import('~/components/AnnotateButtons.vue'),
    TweetAnnotateButtons: () => import('~/components/TweetAnnotateButtons.vue')
  },
  asyncData({ app, params, error }) {
    return app.$axios.$get(`/api/example?id=${params.id}`)
      .then((data) => {
        return {
          title: `[ML-News] ${data.Example.Title}`,
          example: NewExample(data.Example),
          similarExamples: data.SimilarExamples.filter(function(e) {
            return e.Label === 1 || e.Score > 0.0;
          }),
          keywords: data.Keywords,
        }
      })
      .catch((err) => {
        const statusCode = err.response.status;
        error({ 
          statusCode: statusCode,
          message: err.response.data.error
        });
      })
  },
  head() {
    const tweets = this.tweetsWithPositiveLabelOrPositiveScore.map(t => "@" + t.ScreenName + "「" + t.FullText + "」").slice(0, 3);
    const bookmarks = filterBookmarksWithComment(this.example).map(b => "id:" + b.user + "「"+ b.comment + "」").slice(0, 3);
    const description = tweets.join("\n") + bookmarks.join("\n");
    const robotsContent = this.example.Label == -1 ? "noindex, nofollow" : "index, follow";

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
          name: "og:image",
          content: this.example.OgImage 
        },
        {
          name: "robots",
          content: robotsContent
        }
      ],
      link: [
        {
          rel: "canonical",
          href: `https://www.machine-learning.news/example/${this.example.Id}`
        }
      ]
    };
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
  get hasBookmarksWithComment(): boolean {
    return filterBookmarksWithComment(this.example).length > 0;
  }
  get bookmarksWithComment() {
    return filterBookmarksWithComment(this.example);
  }
  get tweetsWithPositiveLabelOrPositiveScore(): Tweet[] {
    return this.example.ReferringTweets.filter(function(t: Tweet) {
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
.hatena-bookmark-link {
  color: #ff4166;
}
.hatena-bookmark-user-icon-img {
  width: 24px;
  height: 24px;
  margin: 2px;
}
.tweets-count, .tweet-retweet-count, .tweet-favorite-count {
  color: #ff4166;
}
.tweet-icon-img {
  width: 24px;
  height: 24px;
  margin: 2px;
}
.tweet-created-at {
  color: #999;
  float: right;
}
</style>