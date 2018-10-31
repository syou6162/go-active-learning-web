<template>
  <b-card-group columns>
    <example 
      v-for="example in examples"
      v-bind:key="example.Url"
      v-bind:example="example"
      ></example>
  </b-card-group>
</template>

<script>
import axios from 'axios';
import Example from './Example.vue';
import NewExample from './util';

export default {
  data () {
    return {
      listname: 'general',
      examples: []
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
      this.examples = [];
      axios.get("/api/examples?listName=" + listname)
      .then(response => {
        this.examples = response.data.map(e => NewExample(e));
      });
    }
  },
  components: {
    "example": Example,
  }
}
</script>
