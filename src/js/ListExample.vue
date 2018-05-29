<template>
  <li>
    <a v-bind:href="example.Url" v-if="example.Title">
      {{ example.Title }}
    </a>
    <a v-bind:href="example.Url" v-else>
      {{ example.Url }}
    </a>
  </li>
</template>

<script>
export default {
  props: ['example'],
}
</script>
