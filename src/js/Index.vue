<template>
  <b-card no-body>
    <b-tabs card>
      <div v-for="listname in listnames">
        <b-tab v-bind:title="listname">
          <b-card-group columns>
            <list-example 
              v-for="example in examples_by_listname[listname]"
              v-bind:key="example.Url"
              v-bind:example="example"
              ></list-example>
          </b-card-group>
        </b-tab>
      </div>
      <b-tab v-bind:id="recent-added-examples" title="recent-added-examples">
        <recent-added-examples></recent-added-examples>
      </b-tab>
    </b-tabs>
  </b-card>
</template>

<script>
import axios from 'axios';
import ListExample from './ListExample.vue';
import RecentAddedExamples from './RecentAddedExamples.vue';

export default {
  data () {
    return {
      listnames: ["general", "twitter", "github", "arxiv", "slideshare", "speakerdeck"],
      examples_by_listname: {
        general: [],
        twitter: [],
        github: [],
        arxiv: [],
        slideshare: [],
        speakerdeck: []
      }
    }
  },
  mounted() {
    for (let listname of this.listnames) {
      axios.get("/api/examples?listName=" + listname)
        .then(response => {
          this.examples_by_listname[listname] = response.data;
        });
    }
  },
  components: {
    "list-example": ListExample,
    "recent-added-examples": RecentAddedExamples
  }
}
</script>
