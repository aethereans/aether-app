<template>
  <div class="spinner" v-show="visible">
    <!-- <template v-if="hidetext">
      <vue-simple-spinner class="spinner-style narrow-container" :spacing="7" message="" :speed="0.33"></vue-simple-spinner>
    </template>
    <template v-else>
      <vue-simple-spinner class="spinner-style" :spacing="7" message="Loading" :speed="0.33"></vue-simple-spinner>
    </template> -->
    <template>
      <vue-simple-spinner
        class="spinner-style narrow-container"
        :spacing="7"
        message=""
        :speed="0.33"
      ></vue-simple-spinner>
    </template>
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-spinner',
  props: {
    delay: {
      type: Number,
      default: 100,
    },
    hidetext: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      visible: false,
    }
  },
  beforeMount(this: any) {
    if (this.delay === 0) {
      this.visible = true
      return
    }
    let vm = this
    setTimeout(function() {
      vm.visible = true
    }, this.delay)
  },
}
</script>

<style lang="scss" scoped>
// @import "../scss/globals";
.spinner {
  display: flex;
  font-family: 'SCP Bold';
  font-size: 14.4px;
  & > div {
    width: 75px;
  }
  .spinner-style.narrow-container {
    width: 32px;
  }
}
</style>
<style lang="scss">
@import '../scss/globals';
.spinner-style {
  // These styles had to be marked !important because the underlying component is actually setting them in its inline styles. Unfortunate.
  .vue-simple-spinner {
    border-color: $a-orange $a-cerulean $a-grey-400 !important;
  }
  .vue-simple-spinner-text {
    color: $a-grey-600 !important;
  }
}
</style>

<!-- Spinner config: https://jsfiddle.net/dzwillia/kc5ka8zu -->
