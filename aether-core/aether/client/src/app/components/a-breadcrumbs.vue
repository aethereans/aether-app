<template>
  <div class="breadcrumbs">
    <router-link class="breadcrumb root" to="/" title="Home">
      <img
        class="root-img"
        src="../ext_dep/images/logo-sm-h.png"
        alt=""
        :class="{ 'at-home': isAtHome }"
      />
      <div class="breadcrumb-text" v-show="isAtHome">Home</div>
    </router-link>
    <router-link
      class="breadcrumb"
      v-for="bc in breadcrumbs"
      :key="bc.Link"
      :title="bc.VisibleName"
      :to="'/' + bc.Link"
    >
      <img
        class="soft-chevron"
        src="../ext_dep/svg/softer-chevron.svg"
        alt=""
      />
      <div class="breadcrumb-text">{{ bc.VisibleName }}</div>
    </router-link>
  </div>
</template>

<script lang="ts">
var Vuex = require('../../../node_modules/vuex').default
export default {
  name: 'a-breadcrumbs',
  data() {
    return {}
  },
  computed: {
    ...Vuex.mapState(['breadcrumbs']),
    isAtHome(this: any) {
      return this.$store.state.route.name === 'Home'
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.breadcrumbs {
  flex: 1;
  display: flex;
  margin-left: 20px;
}

.breadcrumb {
  font-family: 'SSP Semibold'; // letter-spacing: 0.3px;
  -webkit-app-region: no-drag;
  padding-left: 9px;
  padding-right: 11px;
  height: 38px;
  position: relative;
  left: -3px;
  cursor: pointer;
  display: flex; // flex: 1;
  margin-left: -6px;
  min-width: 0;
  &:hover {
    background-color: rgba(255, 255, 255, 0.05);
  }
  @extend %link-hover-ghost-extenders-disable;
  &.root {
    flex: none;
  }
  .breadcrumb-text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
  }
}

.breadcrumb.root .breadcrumb-text {
  padding-left: 5px;
}

.soft-chevron {
  /*width: 36px;*/
  height: 37px;
}

.breadcrumb-text {
  padding-top: 7px;
  color: $a-grey-800;
  padding-left: 12px;
}

.root-img {
  width: 24px;
  height: 24px;
  margin-top: 7px;
  margin-right: -1px;
  &.at-home {
    margin-right: 5px;
  }
}
</style>
