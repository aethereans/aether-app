<template>
  <div class="board-tabs" v-if="!isThread()">
    <router-link
      v-for="tab in tabslist"
      class="board-tab-item"
      :class="{
        'inexact-selected': inexactHighlight,
        'router-link-exact-active': isExactExceptQuery(tab),
      }"
      :key="tab.link"
      :to="tab.link"
      :title="tab.name"
    >
      <div class="tab-name-text">{{ tab.name }}</div>
    </router-link>
  </div>
</template>

<script lang="ts">
// var Vuex = require('../../../node_modules/vuex').default
export default {
  name: 'a-tabs',
  props: ['tabslist', 'inexactHighlight', 'exactExceptQueryHighlight'],
  /*
    Inexact highlight: do we want to show highlight on inexact (parent) match? That is, if you're in cars/sedans, do you want your 'cars' tab to be active? By default, it's not - we look for exact matches for highlight. If you enable this, it does inexact matching. Be mindful that this might mean multiple tabs highlighting at once, if they're not at the same level of hierarchy.
  */
  data() {
    return {}
  },
  computed: {
    // ...Vuex.mapState(['currentBoardFp']),
  },
  methods: {
    isThread: function (this: any) {
      if (this.$route.params.threadfp) {
        return true
      } else {
        return false
      }
    },
    isExactExceptQuery(this: any, tab: any) {
      if (!this.exactExceptQueryHighlight) {
        return false
      }
      let url = this.$route.fullPath
      url = url.substring(0, url.indexOf('?'))
      if (tab.link !== url) {
        return false
      }
      return true
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.board-tabs {
  * {
    user-select: none;
  }
  display: flex;
  border-bottom: 1px solid rgba(0, 0, 0, 0.25);

  .board-tab-item {
    display: flex; // flex: 1;
    height: 38px;
    min-width: 70px;
    color: $a-grey-500;
    &:hover {
      background-color: rgba(255, 255, 255, 0.05);
      color: $a-grey-800; // Disable link hover prettifier
      @extend %link-hover-ghost-extenders-disable;
    }
    .tab-name-text {
      margin: auto;
      padding: 0 20px;
      font-family: 'SSP Semibold';
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    &.selected {
      box-shadow: 0 3px 0 0 #f99157 inset; // why not border: because border pushes the text inside down, so when a tab is selected, text moves. Not good.
      color: $a-grey-800;
      background-color: rgba(255, 255, 255, 0.1); // font-family: "SSP/ Bold";
      // background: linear-gradient(to bottom, rgba(255, 255, 255, 0.1), $a-transparent 90%);
      &:hover {
        background-color: rgba(255, 255, 255, 0.15);
      }
    }
    &.router-link-exact-active {
      @extend .selected;
    }
    &.router-link-active.inexact-selected {
      @extend .selected;
    }
  }
}
</style>
