<template>
  <div class="patreon-button" @click="patreonLinkClicked">
    <img src="../ext_dep/images/patreon.png" alt="" class="patreon-button-image">
  </div>
</template>

<script lang="ts">

export default {
  name: 'a-patreon-button',
  data () {
    return {
    }
  },
  methods: {
    patreonLinkClicked(this: any) {
      const metrics = require('../services/metrics/metrics')()
      metrics.SendRaw("Patreon link clicked")
      window.location.href = "https://patreon.com/nehbit"
    }
  }
}
</script>

<style lang="scss" scoped>
.patreon-button {
  cursor:pointer;
  .patreon-button-image {
    width:175px;
    border-radius: 3px;
    overflow:hidden;
    &:hover {
      filter: hue-rotate(-8deg);
    }
  }
}
</style>