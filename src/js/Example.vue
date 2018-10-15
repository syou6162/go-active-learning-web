<template>
  <a v-bind:href="example.Url">
    <b-card no-body>
      <b-card-body v-bind:title="example | getTitle(75, '...')" class="m-1 p-2">
        <div v-if="example.OgImage === ''">
          <p class="card-text">{{ example | getDescription(250, '...') }}</p>
        </div>
        <div v-else class="d-flex justify-content-between">
          <p class="card-text">{{ example | getDescription(75, '...') }}</p>
          <img class="img-thumbnail img-responsive" style="width: 128px; height: 96px; margin: 3px;" v-bind:src="example.OgImage" onerror="this.style.display='none'" />
        </div>
        <b-card-footer>
          <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-bind:src="example.Favicon" onerror="this.style.display='none'" />
          {{ example | getDomain }} {{ example | getUserName }}
        </b-card-footer>
      </b-card-body>
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
