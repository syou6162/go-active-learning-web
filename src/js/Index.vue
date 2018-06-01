<template>
  <b-card no-body>
    <b-tabs card>
      <div v-for="listname in listnames">
        <b-tab v-bind:title="listname">
          <b-container class="bv-example-row">
            <b-row>
              <list-example 
                v-for="example in examples_by_listname[listname]"
                v-bind:key="example.Url"
                v-bind:example="example"
                ></list-example>
            </b-row>
          </b-container>
        </b-tab>
      </div>
    </b-tabs>
  </b-card>
</template>

<script>
import axios from 'axios';
import ListExample from './ListExample.vue';

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
    "list-example": ListExample
  }
}
</script>
