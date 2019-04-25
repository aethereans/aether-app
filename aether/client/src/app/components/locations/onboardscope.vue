<template>
  <transition name="fade" appear>
    <div class="location onboarding">
      <div class="onboardscope">
        <router-view></router-view>
      </div>
    </div>
  </transition>
</template>

<script lang="ts">
export default {
  name: 'onboardscope',
  data() {
    return {}
  },
}
</script>

<style lang="scss" scoped>
.location.onboarding {
  padding-bottom: 0;
  .onboardscope {
    flex: 1;
    display: flex;
  }
}
</style>
