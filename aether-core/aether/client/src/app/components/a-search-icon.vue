<template>
  <div
    class="notifications-box"
    title="Search"
    @click="$router.push(searchLink)"
  >
    <a-header-icon class="header-search-icon" icon="search"></a-header-icon>
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-search-icon',
  data() {
    return {}
  },
  computed: {
    searchLink(this: any) {
      return '/searchscope'
    },
  },
}
</script>
<style lang="scss">
.header-search-icon {
  svg {
    height: 14px;
  }
}
</style>

<style lang="scss" scoped>
@import '../scss/globals';
.header-search-icon {
  color: $a-grey-800;
}
</style>
