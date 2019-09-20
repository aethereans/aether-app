<template>
  <div id="app">
    <transition name="component-fade" mode="out-in">
      <div class="app-loader-container" v-show="preloadSpinnerVisible">
        <transition name="fade" appear>
          <a-main-app-loader v-show="preloadSpinnerVisible"></a-main-app-loader>
        </transition>
      </div>
    </transition>
    <transition name="fade" appear>
      <div class="onboard-route-container" v-if="onboardVisible">
        <router-view name="onboarding"></router-view>
      </div>
    </transition>
    <template v-if="mainInterfaceVisible">
      <div class="main-top">
        <a-side-header></a-side-header>
        <a-header></a-header>
      </div>
      <div class="app-content">
        <div class="side-block" v-show="$store.state.sidebarOpen">
          <a-sidebar></a-sidebar>
        </div>
        <div
          class="main-block"
          :class="{ 'sidebar-closed': !$store.state.sidebarOpen }"
        >
          <router-view></router-view>
        </div>
      </div>
    </template>
  </div>
</template>
<script lang="ts">
// Check out renderermain.ts to see what's normally ought to be here.
export default {
  data() {
    return {}
  },
  computed: {
    preloadSpinnerVisible(this: any) {
      return (
        !this.$store.state.onboardCompleteStatusArrived ||
        !this.$store.state.firstRunAfterUpdateStatusArrived
      )
    },
    onboardVisible(this: any) {
      return (
        this.$store.state.onboardCompleteStatusArrived &&
        !this.$store.state.onboardCompleteStatus
      )
    },
    mainInterfaceVisible(this: any) {
      return !this.onboardVisible
    },
  },
}
</script>
<style lang="scss">
@import '../scss/bulmastyles';
@import '../scss/typefaces';
@import '../scss/globals';
@import '../scss/tooltipstyles';
@import '../scss/markdowned';

::selection {
  background: #b4d5fe;
  /* Fix for Chromium issue #641509
        https://bugs.chromium.org/p/chromium/issues/detail?id=641509 */
}

:focus {
  outline: 7px auto $a-orange;
  /*
      'auto' is the trick here.
      If you make this 'solid' to get the actual colour and opaqueness you want, what happens is that the focus event is triggered at every click, and that leaves a focus ring on the menu button after the click. Auto is smarter than that, but it also controls the opaqueness of the colour for some reason, as well as its thickness.

      Given the option between
      a) Disabling focus entirely (accesibility, usability via keyboard sucks, nonstarted)
      b) Having a solid styled outline, which makes all clicks leave a focus like a footprint
      c) Just styling the colour of the focus and be OK with the uglier transparency-based, fuzzy, but natively handled smart focus chooser..

      I'll go with C.
     */
}

* {
  box-sizing: border-box;
  user-select: none; // ^ Don't forget to make sure all content is still selectable as default! Above is meant for things like button text and weird flashes of blue selection only.
  // &:focus {
  //   outline: 3px solid $a-orange; // This is for accessibility - let's make this as high contrast as possible.
  // }
  a {
    // color: inherit;
    text-decoration: none;
    color: $a-cerulean;
    cursor: pointer;
    &:hover {
      position: relative; // color: inherit;
      background-color: rgba(255, 255, 255, 0.05);
      &::before {
        position: absolute;
        width: $link-hover-ghost-extension-length;
        height: 100%;
        background-color: rgba(255, 255, 255, 0.05);
        content: '';
        left: -$link-hover-ghost-extension-length;
        top: 0;
        border-radius: 2px 0 0 2px;
      }
      &::after {
        position: absolute;
        width: $link-hover-ghost-extension-length;
        height: 100%;
        background-color: rgba(255, 255, 255, 0.05);
        content: '';
        right: -$link-hover-ghost-extension-length;
        top: 0;
        border-radius: 0 2px 2px 0;
      }
    }
  }
  ::after,
  ::before {
    box-sizing: inherit;
  }
}

video::-webkit-media-controls-panel {
  background-color: rgba(25, 25, 25, 0.9);
}

html {
  height: 100%; // background-color: $dark-base;
  background-color: $mid-base;
  overflow-y: auto;
  text-rendering: optimizeLegibility;
  text-size-adjust: 100%;
  -webkit-font-smoothing: antialiased;
  overflow: hidden;
}

body {
  font-family: 'SSP Bold';
  height: 100%;
  display: flex;
  box-sizing: border-box;
  color: $a-grey-800;
  margin: 0;
  font-weight: 400;
  line-height: 1.5;
  font-size: 1rem;
}

.main-block * {
  user-select: text;
}

@include generateScrollbar($a-grey-100);
</style>

<style lang="scss" scoped>
@import '../scss/globals';
#app {
  width: 100%;
  min-height: 100%;
  display: flex;
}

.main-top {
  background-color: #16222a; // $dark-base * 0.8;
  z-index: 3;
}

.side-block {
  width: $sidebar-width; // max-width: $max-sidebar-width;
  // min-width: $min-sidebar-width;
  background-color: $dark-base;
  display: flex;
  flex-direction: column; // box-shadow: 0 1px 2px rgba(0, 0, 0, 0.35);
}

.main-block {
  flex: 1;
  overflow-y: scroll; // box-shadow: $line-separator-shadow-castleft;
  box-shadow: -2px 0 2px -1px rgba(0, 0, 0, 0.2);
  border-radius: 0 0 0 10px;
  background-color: $mid-base;
  height: 100%;
  z-index: 2;

  &.sidebar-closed {
    border-radius: 0;
  }
}

.location {
  padding-bottom: 50px;
}

#app {
  display: flex;
  flex-direction: column;
  flex: 1; // border-radius: 3px 0 0 3px;
  // overflow: hidden;
  .main-top {
    display: flex;
  }
}

.app-content {
  flex: 1;
  display: flex;
  @include generateScrollbar($mid-base);
  background-color: #16222a; // $dark-base * 0.8;
} // .side-top {
//   -webkit-app-region: drag;
//   /* ^ You must mark all interactive objects within this as NON-draggable to make this work, otherwise they will just be effectively unclickable. */
//   height: $top-bar-height;
//   background-color: $dark-base*0.8;
//   /*border-bottom: 1px solid #111;*/
//   box-shadow: 0 1px 2px rgba(0, 0, 0, 0.35);
//   display: flex;
//   padding-right: 5px;
//   position: relative;
//   z-index: 2;
// }
.location {
  flex: 1;
  width: 100%; // overflow-y: scroll;
  display: flex;
  flex-direction: column;
}

.app-loader-container {
  // width: 100%;
  // height: 100%;
  display: flex;
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 999;
  background-color: $mid-base;
}

.onboard-route-container {
  display: flex;
  min-height: 100%;
}

.fade-enter-active,
.fade-leave-active {
  transition-property: opacity;
  transition-duration: 0.25s;
}

.fade-enter-active {
  transition-delay: 0.25s;
}

.fade-enter,
.fade-leave-active {
  opacity: 0;
}
</style>
