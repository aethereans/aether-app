<template>
  <div class="location global" id="scrolltarget-container-target">
    <div class="globalscope">
      <a-global-header></a-global-header>
      <a-tabs v-show="tabsVisible()" :tabslist="tabslist"></a-tabs>
      <router-view></router-view>
    </div>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'globalscope',
    data() {
      return {}
    },
    computed: {
      tabslist(this: any) {
        return [{
          'name': 'ALL',
          'link': '/globalscope'
        }, {
          'name': 'SUBBED',
          'link': '/globalscope/subbed'
        }]
      }
    },
    methods: {
      tabsVisible(this: any): boolean {
        if (this.$store.state.route.name === 'Global>NewBoard') {
          return false
        }
        return true
      }
    }
  }
</script>

<style lang="scss" scoped>

</style>