<template>
  <a v-bind:href="example.Url">
    <b-card v-bind:title="example | getTitle(75, '...')">
      <p class="card-text">{{ example | getDescription(250, '...') }}</p>
      <b-card-footer>
        {{ example | getDomain }} {{ example | getUserName }}
      </b-card-footer>
    </b-card>
  </a>
</template>

<script>
function getDomain(example) {
  var url = example.Url;
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
  props: ['example'],
  filters: {
    getTitle: function(example, length, omission) {
      var title = example.Title ? example.Title : example.Url;
      var domain = getDomain(example);
      if ('arxiv.org' === domain) {
        return title;
      }
      return truncate(title, length, omission);
    },
    getDescription: function(example, length, omission) {
      var title = example.Title ? example.Title : example.Url;
      var desc = example.OgDescription ? example.OgDescription : (example.Description ? example.Description : example.CleanedText);
      var domain = getDomain(example);
      if ('arxiv.org' === domain) {
        desc = '';
      }
      return truncate(desc, length, omission);
    },
    getDomain: getDomain,
    getUserName: function(example) {
      var domain = getDomain(example);
      var url = example.Url;
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
