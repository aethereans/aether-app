<template>
  <div class="onboard-sublocation">
    <div class="onboard-carrier-root">
      <img
        class="logo-outline"
        src="../../../ext_dep/images/Logo-v3-outline.svg"
        alt=""
      />
      <div class="enter-box">
        <router-link class="button is-outlined start-button" to="/onboard/1">
          BEGIN
        </router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'onboardroot',
  data() {
    return {}
  },
}
</script>

<style lang="scss" scoped>
@import '../../../scss/globals';
.onboard-sublocation {
  display: flex;
  flex: 1;
  min-height: 100%;
  .onboard-carrier-root {
    margin: auto;
    height: 180px;
    display: flex; // width: 240px;
    animation-duration: 1.6s;
    animation-name: SIDE_LOAD;
    overflow: hidden;
    width: 245px;

    .logo-outline {
      width: 155px;
      height: 155px;
    }

    .enter-box {
      display: flex;
      margin-left: 10px;
      height: 155px;
      .start-button {
        margin: auto;
        @extend %link-hover-ghost-extenders-disable;
        background-color: $a-transparent;
        color: $a-grey-800;
        &:hover {
          background-color: $a-grey-800;
          color: $mid-base;
        }
      }
    }
  }
}

@keyframes SIDE_LOAD {
  0% {
    width: 155px;
  }
  60% {
    width: 155px;
  }
  100% {
    width: 245px;
  }
}
</style>
