<template>
  <div>
    <div v-if="loading">
      Now loading...
    </div>
    <div v-else-if="error">
      Fail to retrieve from API server. Error: {{ error }}
    </div>
    <div v-else>
      <vue-headful 
        v-bind:title="title" 
        />
      <b-form-group label="Please select a label">
        <b-form-radio-group
          buttons
          v-model="isNew"
          button-variant="outline-primary"
          :options="options" />
      </b-form-group>
      <b-card-group columns>
        <example 
          v-for="example in examplesFilterByIsNew(isNew)"
          v-bind:key="example.Url"
          v-bind:example="example"
          ></example>
      </b-card-group>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Example from './Example.vue';
import NewExample from './util';

export default {
  data () {
    return {
      title: "ML News",
      listname: 'general',
      examples: [],
      isNew: 0,
      options: [
        { text: 'All', value: 0 },
        { text: 'Recent', value: 1 },
      ],
      error: null,
      loading: true,
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
      let self = this;
      this.loading = true;
      this.error = null;

      axios.get("/api/examples?listName=" + listname)
        .then(response => {
          this.examples = response.data.map(e => NewExample(e));
          this.listname = this.$route.params.listname;
          this.title = `ML News - ${this.listname}`;
          this.loading = false;
        }).catch(function (error) {
          if (error.response) {
            self.error = error.response.statusText;
          }
        });
    },
    examplesFilterByIsNew: function(isNew) {
      return this.examples.filter(function(e) {
        if (isNew == 0) {
          return true;
        } else {
          return e.IsNew == isNew;
        }
      })
    }
  },
  components: {
    "example": Example,
  }
}
</script>
