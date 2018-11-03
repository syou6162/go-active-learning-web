<template>
  <div>
    <b-card no-body>
      <b-card-body v-bind:title="example | getTitle(75, '...')" title-tag="h5" class="m-1 p-2">
        <b-card-footer>
          <b-button v-bind:href="'/example/' + encodeURIComponent(example.FinalUrl)" class="float-right" size="sm">Read more</b-button>
          <img v-if="example.Favicon" style="width: 16px; height: 16px;" v-lazy="example.Favicon" onerror="this.style.display='none'" />
          <a v-bind:href="example.FinalUrl">{{ example | getDomain }} {{ example | getUserName }}</a>
        </b-card-footer>
      </b-card-body>
    </b-card>
  </div>
</template>

<script>
export default {
  data () {
    return {
      modalShow: false
    }
  },
  props: ['example'],
}
</script>
