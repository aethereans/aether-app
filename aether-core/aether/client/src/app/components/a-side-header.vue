<template>
  <div
    class="side-header"
    :class="{ 'sidebar-collapsed': !$store.state.sidebarOpen }"
  >
    <div class="app-name" @click="goHome">
      <span :class="{ 'non-fullscreen': !isFullscreen }"> aether </span>
    </div>
    <!-- <div class="spacer"></div> -->
    <div class="sidebar-expander-box">
      <a-header-icon
        icon="chevron-up"
        @click.native="toggleSidebarState"
        v-show="$store.state.sidebarOpen"
      ></a-header-icon>
      <a-header-icon
        icon="chevron-down"
        @click.native="toggleSidebarState"
        v-show="!$store.state.sidebarOpen"
      ></a-header-icon>
    </div>
  </div>
</template>

<script lang="ts">
var Vuex = require('../../../node_modules/vuex').default
export default {
  name: 'a-side-header',
  data() {
    return {}
  },
  methods: {
    ...Vuex.mapActions(['setSidebarState']),
    toggleSidebarState() {
      this.$store.state.sidebarOpen === true
        ? this.setSidebarState(false)
        : this.setSidebarState(true)
    },
    goHome(this: any) {
      this.$router.push('/')
    },
  },
  computed: {
    isFullscreen(this: any) {
      if (process.platform !== 'darwin') {
        return true //
      }
      return this.$store.state.appIsFullscreen
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';

.side-header {
  display: flex;
  width: $sidebar-width; // flex: 1;
  // max-width: $max-sidebar-width; // flex: 1;
  // min-width: $min-sidebar-width;
  -webkit-app-region: drag;
  /* ^ You must mark all interactive objects within this as NON-draggable to make this work, otherwise they will just be effectively unclickable. */
  // height: $top-bar-height;
  background-color: #16222a; // $dark-base * 0.8;
  /*border-bottom: 1px solid #111;*/
  // box-shadow: $line-separator-shadow;
  box-shadow: -1px 1px 1px rgba(0, 0, 0, 0.35);
  display: flex;
  padding-right: 5px;
  position: relative; // z-index: 3;
  .spacer {
    flex: 1;
  }
  &.sidebar-collapsed {
    // I know this is nitpicky, but I want to have the shadow look exactly consistent in visual weight when the underlying colour changes.
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.42);
  }
}

.app-name {
  font-family: 'SCP Semibold';
  display: flex; // margin-left: 15px;
  flex: 1;
  cursor: pointer;
  span {
    margin: auto; // margin-left: 15px;
    color: $a-grey-300; // margin-left: 15px;
    // transition: letter-spacing 0.15s cubic-bezier(0.68, -0.55, 0.265, 1.55), margin-right 0.2s cubic-bezier(0.68, -0.55, 0.265, 1.55);
    transition: margin-right 0.225s cubic-bezier(0.68, -0.55, 0.265, 1.55);
    position: relative; // letter-spacing: 15px;
    letter-spacing: 1px;
    bottom: 1px; // margin-right: 110px;
    // margin-right: 100px;
    margin-right: 119px;
    &.non-fullscreen {
      // In this case, we need to centre the app name.
      // margin: auto;
      margin-right: 44px;
    }
  }
}
</style>
