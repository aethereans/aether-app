<template>
  <div class="sfwlist-box" title="Enable/Disable SFW Lists">
    <a-header-icon
      id="toggleSFWListIcon"
      class="sfwlist-icon-safe"
      v-show="!$store.state.ambientStatus.frontendambientstatus.sfwlistdisabled"
      @click.native="toggleSFWListDisabled"
      icon="ban"
    ></a-header-icon>
    <a-header-icon
      id="toggleSFWListIcon"
      class="sfwlist-icon-unsafe"
      v-show="$store.state.ambientStatus.frontendambientstatus.sfwlistdisabled"
      @click.native="toggleSFWListDisabled"
      icon="ban"
    ></a-header-icon>
  </div>
</template>

<script lang="ts">
var fe = require('../services/feapiconsumer/feapiconsumer')
export default {
  name: 'a-sfwlist-icon',
  data(this: any) {
    return {}
  },
  methods: {
    async toggleSFWListDisabled(this: any) {
      fe.SendSFWListDisabledStatus(
        !this.$store.state.ambientStatus.frontendambientstatus.sfwlistdisabled,
        function () {}
      )
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.sfwlist-icon-safe {
  color: $a-grey-800;
}

.sfwlist-icon-unsafe {
  color: $a-red;
}
</style>
