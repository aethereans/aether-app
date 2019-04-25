<template>
  <div class="location search" id="scrolltarget-container-target">
    <div class="searchscope">
      <a-search-header></a-search-header>
      <a-tabs :tabslist="tabslist" :exactExceptQueryHighlight="true"></a-tabs>
      <router-view></router-view>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'searchscope',
  data() {
    return {}
  },
  computed: {
    tabslist(this: any) {
      return [
        {
          name: 'COMMUNITIES',
          link: '/searchscope',
        },
        {
          name: 'CONTENT',
          link: '/searchscope/content',
        },
        {
          name: 'PEOPLE',
          link: '/searchscope/user',
        },
      ]
    },
  },
  methods: {},
}
</script>

<style lang="scss" scoped></style>
