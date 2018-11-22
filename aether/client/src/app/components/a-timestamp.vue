<template>
  <div class="timestamp" v-if="creation > 0">
    {{computeTimestamp()}}
  </div>
</template>

<script lang="ts">
  var globalMethods = require('../services/globals/methods')
  export default {
    name: 'a-timestamp',
    props: {
      creation: {
        type: Number,
        default: 0
      },
      lastupdate: {
        type: Number,
        default: 0
      }
    },
    data() {
      return {}
    },
    computed: {
      nowTs: function(): number {
        return Math.floor(Date.now() / 1000)
      }
    },
    methods: {
      computeTimestamp(this: any): string {
        let ts = this.creation > this.lastupdate ? this.creation : this.lastupdate
        let edited = this.creation > this.lastupdate ? "" : " (edited)"
        return globalMethods.TimeSince(ts) + edited
      },
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .timestamp {
    font-family: "SSP Regular Italic";
    color: $a-grey-600;
  }
</style>