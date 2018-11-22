<template>
  <div class="guidelight-bucket">
    <!-- Unknown -->
    <div class="guide-light" hasGuidelightTooltip title="Status data from this component has not loaded." v-show="notifyState==='status_section_unknown' || notifyState==='status_subsection_unknown'">
      <div class="guidelight-dot blue"></div>
    </div>
    <!-- All OK -->
    <div class="guide-light" hasGuidelightTooltip title="All components are working fine." v-show="notifyState==='status_section_ok' || notifyState==='status_subsection_ok'">
      <div class="guidelight-dot green"></div>
    </div>
    <!-- Section warn -->
    <div class="guide-light" hasGuidelightTooltip title="Some components are operating in degraded mode, or are under stress." v-show="notifyState==='status_section_warn'">
      <div class="guidelight-dot yellow"></div>
    </div>
    <!-- Subsection warn -->
    <div class="guide-light" hasGuidelightTooltip title="This component works, but is operating in degraded mode, or is under stress." v-show="notifyState==='status_subsection_warn'">
      <div class="guidelight-dot yellow"></div>
    </div>
    <!-- Section fail -->
    <div class="guide-light" hasGuidelightTooltip title="Some components are having significant issues and they are not functioning correctly." v-show="notifyState==='status_section_fail'">
      <div class="guidelight-dot red"></div>
    </div>
    <!-- Subsection fail -->
    <div class="guide-light" hasGuidelightTooltip title="This component is having significant issues and it is not functioning correctly." v-show="notifyState==='status_subsection_fail'">
      <div class="guidelight-dot red"></div>
    </div>
  </div>
</template>

<script lang="ts">
  var Tooltips = require('../services/tooltips/tooltips')
  export default {
    name: 'a-guidelight',
    props: {
      notifyState: {
        type: String,
        default: "ok" // warn, fail
      }
    },
    data() {
      return {}
    },
    mounted(this: any) {
      Tooltips.MountGuidelightTooltip()
    },
    updated(this: any) {},
    computed: {},
    methods: {}
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .guidelight-bucket {
    display: flex;
  }

  .guide-light {
    display: flex;
    cursor: pointer;
    .guidelight-dot {
      margin: auto;
      width: 36px;
      height: 12px;
      border-radius: 100px;
      background-color: $a-cerulean;
      &.green {
        background-color: $a-green;
      }
      &.yellow {
        background-color: $a-yellow;
      }
      &.red {
        background-color: $a-red;
      }
      &.blue {
        background-color: $a-cerulean;
      }
    }
  }
</style>