<template>
  <div>
    <div v-if="loading">
      Now loading...
    </div>
    <div v-else-if="error">
      Fail to retrieve from API server. Error: {{ error }}
    </div>
    <div v-else class="mx-auto" style="max-width: 40rem;">
      <vue-headful 
        v-bind:title="title" 
        v-bind:description="example | getDescriptionForSearchEngine"
        v-bind:url="example.FinalUrl"
        />
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
        Date: 
        {{ example.CreatedAt.tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm") }}
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
          <b-list-group-item v-for="t in example.ReferringTweets.slice(0, 3)" :key="t.ScreenName">
            <img v-if="t.ProfileImageUrl" style="width: 24px; height: 24px;" v-lazy="t.ProfileImageUrl" onerror="this.style.display='none'" />
            <a v-bind:href="'https://twitter.com/' + t.ScreenName + '/status/' + t.IdStr">@{{ t.ScreenName }}</a>
            {{ t.FullText }}
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
          <b-button v-bind:href="'/example/' + encodeURIComponent(example.FinalUrl)" class="float-right" size="sm">Read more</b-button>
          <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
          {{ example | getTitle(100, '...') }}
        </b-list-group-item>
      </b-list-group>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Example from './Example.vue';
import HatenaBookmarkIcon from './HatenaBookmarkIcon.vue';
import AnnotateButtons from './AnnotateButtons.vue';
import TwitterIcon from './TwitterIcon.vue';
import { NewExample, filterBookmarksWithComment, IsAdmin } from './util';

export default {
  data () {
    return {
      title: "ML News",
      url: this.$route.params.url,
      example: null,
      similarExamples: [],
      keywords: [],
      error: null,
      loading: true,
      isAdmin: false,
    }
  },
  mounted() {
    this.fetchExample(this.$route.params.url);
  },
  methods: {
    fetchExample(url) {
      let self = this;
      this.loading = true;
      this.error = null;

      this.examples = [];
      axios.get("/api/example?url=" + encodeURIComponent(url))
        .then(response => {
          this.example = NewExample(response.data.Example);
          this.similarExamples = response.data.SimilarExamples.filter(function(e) {
            return e.Label === 1 || e.Score > 0.0;
          });
          this.keywords = response.data.Keywords;
          this.title = `ML News - ${this.example.Title}`;
          this.loading = false;
        }).catch(function (error) {
          if (error.response) {
            self.loading = false;
            self.error = error.response.statusText;
          }
        });
      this.isAdmin = IsAdmin();
    }
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
    "annotate-buttons": AnnotateButtons 
  }
}
</script>
