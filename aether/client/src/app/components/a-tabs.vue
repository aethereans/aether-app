<template>
  <div class="board-tabs" v-if="!isThread()">

    <router-link v-for="tab in tabslist" class="board-tab-item" :key="tab.link" :to="tab.link" :title="tab.name">
      <div class="tab-name-text">{{tab.name}}</div>
    </router-link>
    <!-- <router-link class="board-tab-item" :to="'/board/'+this.currentBoardFp+'/info'">
      <div class="tab-name-text">INFO</div>
    </router-link>
    <router-link class="board-tab-item" :to="'/board/'+this.currentBoardFp+'/modactivity'">
      <div class="tab-name-text">MODERATION</div>
    </router-link>
    <router-link class="board-tab-item" :to="'/board/'+this.currentBoardFp+'/elections'">
      <div class="tab-name-text">ELECTIONS</div>
    </router-link> -->
  </div>
</template>

<script lang="ts">
  // var Vuex = require('../../../node_modules/vuex').default
  export default {
    name: 'a-tabs',
    props: ['tabslist'],
    data() {
      return {}
    },
    computed: {
      // ...Vuex.mapState(['currentBoardFp']),
    },
    methods: {
      isThread: function(this: any) {
        if (this.$route.params.threadfp) {
          return true
        } else {
          return false
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
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
        font-family: "SSP Semibold";
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
      &.selected {
        box-shadow: 0 3px 0 0 #F99157 inset; // why not border: because border pushes the text inside down, so when a tab is selected, text moves. Not good.
        color: $a-grey-800;
        background-color: rgba(255, 255, 255, 0.1); // font-family: "SSP/ Bold";
        // background: linear-gradient(to bottom, rgba(255, 255, 255, 0.1), $a-transparent 90%);
        &:hover {
          background-color: rgba(255, 255, 255, 0.15)
        }
      }
      &.router-link-exact-active {
        @extend .selected;
      }
    }
  }
</style>