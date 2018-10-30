<template>
  <div>
    <b-card no-body>
      <b-card-body v-bind:title="example | getTitle(75, '...')" title-tag="h5" class="m-1 p-2">
        <b-card-footer>
          <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
          <a v-bind:href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
          <b-button @click="modalShow = !modalShow" class="float-right" size="sm">Read more</b-button>
        </b-card-footer>
      </b-card-body>
    </b-card>
    <b-modal v-model="modalShow" hide-footer=true>
      <img v-if="example.OgImage" class="img-thumbnail img-responsive" style="width: 128px; height: 96px; margin: 3px; float: right;" v-lazy="example.OgImage" onerror="this.style.display='none'" />
      <p class="card-text">{{ example | getDescription(500, '...') }}</p>
      <a v-bind:href="example.FinalUrl">Permalink</a>
    </b-modal>
  </div>
</template>

<script>
function getDomain(example) {
  var url = example.FinalUrl;
  return url.replace('http://','').replace('https://','').split(/[/?#]/)[0];
}

function truncate(str, length, omission) {
  str = str ? str : '';
  var length = length ? parseInt(length, 10) : 20;
  var omission = omission ? omission.toString() : '...';

  if (str.length <= length) {
    return str;
  }
  else {
    return str.substring(0, length) + omission;
  }
}

export default {
  data () {
    return {
      modalShow: false
    }
  },
  props: ['example'],
  filters: {
    getTitle: function(example, length, omission) {
      var title = example.Title ? example.Title : example.Url;
      var domain = getDomain(example);
      return truncate(title, length, omission);
    },
    getDescription: function(example, length, omission) {
      var title = example.Title ? example.Title : example.Url;
      var desc = example.OgDescription ? example.OgDescription : example.Description;
      return truncate(desc, length, omission);
    },
    getDomain: getDomain,
    getUserName: function(example) {
      var domain = getDomain(example);
      var url = example.FinalUrl;
      var paths = url.replace('http://','').replace('https://','').split(/[/?#]/);
      if (paths.length === 0) {
        return;
      } else if ('twitter.com' === domain) {
        return '(@' + paths[1] + ')';
      } else if ('github.com' === domain) {
        return '(@' + paths[1] + ')';
      } else if ('qiita.com' === domain) {
        return '(@' + paths[1] + ')';
      } else if ('www.slideshare.net' === domain) {
        return '(id:' + paths[1] + ')';
      } else if ('speakerdeck.com' === domain) {
        return '(id:' + paths[1] + ')';
      } else {
        return;
      }
    }
  }
}
</script>
