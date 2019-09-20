<template>
  <div class="header-icon-container">
    <template v-if="!stacked">
      <icon :name="icon"></icon>
    </template>
    <template v-else>
      <icon>
        <icon
          :name="iconunder"
          scale="1.5"
          class="icon-background-color"
        ></icon>
        <icon :name="icon"></icon>
      </icon>
    </template>
  </div>
</template>
<script lang="ts">
const isDev = require('electron-is-dev')
if (isDev) {
  var Vue = require('../../../node_modules/vue/dist/vue.js') // Production
} else {
  var Vue = require('../../../node_modules/vue/dist/vue.min.js') // Production
}
export default Vue.extend({
  // props: ['icon']
  props: {
    icon: {
      type: String,
      default: '',
    },
    iconunder: {
      type: String,
      default: '',
    },
  },
  computed: {
    stacked(this: any) {
      if (this.iconunder.length > 0) {
        return true
      }
      return false
    },
  },
})
</script>
<style lang="scss">
@import '../scss/globals';
.header-icon-container {
  -webkit-app-region: no-drag;
  height: $top-bar-height;
  width: $top-bar-height;
  display: flex;
  cursor: pointer;
  svg {
    margin: auto;
  }
  &:hover {
    background-color: rgba(255, 255, 255, 0.05);
  }
  .icon-background-color {
    fill: $a-cerulean-90;
  }
}
</style>
