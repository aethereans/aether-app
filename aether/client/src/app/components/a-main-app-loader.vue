<template>
  <!-- <transition appear> -->

  <div class="main-app-loader">
    <div class="spinner-container" v-show="visible">
      <vue-simple-spinner class="spinner-style" :size="125" :line-size="4" :speed="0.33"></vue-simple-spinner>
    </div>
  </div>
  <!-- </transition> -->
</template>

<script lang="ts">
  export default {
    name: 'a-main-app-loader',
    data() {
      return {
        delay: 200,
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
    }
  }
</script>

<style lang="scss" scoped>
  .main-app-loader {
    height: 150px;
    margin: auto;
    .spinner-container {
      position: relative;
      .aether-core {
        position: absolute; // opacity: 0.4;
        width: 125px;
        left: 0;
        top: 0;
      }
    }
  }
</style>