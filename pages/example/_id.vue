<template>
  <div class="mx-auto" style="max-width: 40rem;">
    <b-card v-if="example" v-bind:title="example | getTitle(1000, '...')" tag="article">
      <img v-if="example.OgImage" class="img-thumbnail img-responsive" style="width: 128px; height: 96px; margin: 3px; float: right;" v-lazy="example.OgImage" onerror="this.style.display='none'" />
      <p class="card-text">
        {{ example | getDescription(1000, '...') }}
      </p>
      <div v-if="example.HatenaBookmark.count > 0">
        <a v-bind:href="example.HatenaBookmark.entry_url" style="color: #ff4166;">{{ example.HatenaBookmark.count }} users</a>:
        <hatena-bookmark-icon 
          v-for="b in example.HatenaBookmark.bookmarks.slice(0, 9)"
          v-bind:bookmark="b"
          v-bind:key="b.user"
          ></hatena-bookmark-icon>
      </div>
      <div v-if="example.ReferringTweets && example.ReferringTweets.length > 0">
        <span style="color: #ff4166;">{{ example.ReferringTweets.length }} mentions</span>:
        <twitter-icon
           v-for="tweet in example.ReferringTweets.slice(0, 8)"
           v-bind:tweet="tweet"
           v-bind:key="tweet.ScreenName"
           ></twitter-icon>
      </div>
      <div v-if="keywords.length > 0">
        Keywords: 
        <b-link 
          v-bind:href="'/search?query=' + k" 
          v-for="k in keywords"
          v-bind:key="k"
          >
          {{ k }}
        </b-link>
      </div>
      Date: {{ example.CreatedAt }}
      <annotate-buttons
        v-if="isAdmin"
        v-bind:example="example"
        ></annotate-buttons>
      <b-card-footer>
        <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
        <a v-bind:href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
      </b-card-footer>
    </b-card>
    <div v-if="example.ReferringTweets && example.ReferringTweets.length > 0">
      <h4>Referring Tweets</h4>
      <b-list-group>
        <b-list-group-item v-for="t in example.ReferringTweets.slice(0, 9)" :key="t.ScreenName">
          <img v-if="t.ProfileImageUrl" style="width: 24px; height: 24px;" v-lazy="t.ProfileImageUrl" onerror="this.style.display='none'" />
          <a v-bind:href="'https://twitter.com/' + t.ScreenName + '/status/' + t.IdStr">@{{ t.ScreenName }}</a>
          {{ t.FullText }}
          <tweet-annotate-buttons
            v-if="isAdmin"
            v-bind:tweet="t"
            ></tweet-annotate-buttons>
        </b-list-group-item>
      </b-list-group>
    </div>
    <div v-if="hasBookmarksWithComment">
      <h4>Bookmark Comments</h4>
      <b-list-group>
        <b-list-group-item v-for="b in bookmarksWithComment.slice(0, 3)" :key="b.user">
          <img style="width: 24px; height: 24px; margin: 2px" v-lazy="'https://cdn.profile-image.st-hatena.com/users/' + b.user+ '/profile.png'" />
          <a v-bind:href="'http://b.hatena.ne.jp/' + b.user">id:{{ b.user }}</a>
          {{ b.comment }}
        </b-list-group-item>
      </b-list-group>
    </div>
    <h4 v-if="similarExamples.length > 0">Related Entries</h4>
    <b-list-group>
      <b-list-group-item v-for="example in similarExamples" :key="example.FinalUrl">
        <b-button v-bind:href="example | getExampleUrl" class="float-right" size="sm">Read more</b-button>
        <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
        {{ example | getTitle(100, '...') }}
      </b-list-group-item>
    </b-list-group>
  </div>
</template>

<script>
import axios from 'axios';
import { Auth } from 'aws-amplify';
import Example from '~/components/Example.vue';
import HatenaBookmarkIcon from '~/components/HatenaBookmarkIcon.vue';
import TwitterIcon from '~/components/TwitterIcon.vue';
import AnnotateButtons from '~/components/AnnotateButtons.vue';
import TweetAnnotateButtons from '~/components/TweetAnnotateButtons.vue';
import { NewExample, filterBookmarksWithComment } from '~/assets/util';

export default {
  data () {
    return {
      title: "ML News",
      url: this.$route.params.url,
      example: null,
      similarExamples: [],
      keywords: [],
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
  asyncData({ app, params, error }) {
    return app.$axios.$get(`/api/example?id=${params.id}`)
      .then((data) => {
        return {
          title: `ML News - ${data.Example.Title}`,
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
    const tweets = this.example.ReferringTweets.map(t => "@" + t.ScreenName + "「" + t.FullText.substr(0, 100) + "...」").slice(0, 3);
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
          name: "og:type",
          content: "article"
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
  },
  computed: {
    hasBookmarksWithComment() {
      return filterBookmarksWithComment(this.example).length > 0;
    },
    bookmarksWithComment() {
      return filterBookmarksWithComment(this.example);
    },
  },
  components: {
    "hatena-bookmark-icon": HatenaBookmarkIcon,
    "twitter-icon": TwitterIcon,
    "annotate-buttons": AnnotateButtons,
    "tweet-annotate-buttons": TweetAnnotateButtons
  }
}
</script>
