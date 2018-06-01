<template>
  <b-col md=4>
    <b-card v-bind:href="example.Url" v-if="example.Title">
      {{ example.Title }}
    </b-card>
    <b-card v-bind:href="example.Url" v-else>
      {{ example.Url }}
    </b-card>
  </b-col>
</template>

<script>
export default {
  props: ['example'],
}
</script>
