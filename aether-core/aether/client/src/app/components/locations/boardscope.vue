<template>
  <div class="location board" id="scrolltarget-container-target">
    <!-- Container target is important to get the fix-to-bottom in create panes to work. -->
    <template v-if="!$store.state.currentBoardLoadComplete">
      <div class="spinner-container">
        <a-spinner :delay="300"></a-spinner>
      </div>
    </template>
    <template v-else>
      <template v-if="entityNotFound">
        <a-notfound></a-notfound>
      </template>
      <template v-else>
        <div class="boardscope">
          <a-boardheader></a-boardheader>
          <a-tabs v-show="tabsVisible()" :tabslist="tabslist"></a-tabs>
          <router-view></router-view>
        </div>
      </template>
    </template>
  </div>
</template>

<script lang="ts">
export default {
  name: 'boardscope',
  computed: {
    tabslist(this: any) {
      let modTabsList = [
        {
          name: 'POPULAR',
          link: '/board/' + this.$store.state.currentBoardFp,
        },
        {
          name: 'NEW',
          link: '/board/' + this.$store.state.currentBoardFp + '/new',
        },
        {
          name: 'INFO',
          link: '/board/' + this.$store.state.currentBoardFp + '/info',
        },
        {
          name: 'REPORTS',
          link: '/board/' + this.$store.state.currentBoardFp + '/reports',
        },
        {
          name: 'MOD ACTIVITY',
          link: '/board/' + this.$store.state.currentBoardFp + '/modactivity',
        },
        {
          name: 'ELECTIONS',
          link: '/board/' + this.$store.state.currentBoardFp + '/elections',
        },
      ]
      let nonModTabsList = [
        {
          name: 'POPULAR',
          link: '/board/' + this.$store.state.currentBoardFp,
        },
        {
          name: 'NEW',
          link: '/board/' + this.$store.state.currentBoardFp + '/new',
        },
        {
          name: 'INFO',
          link: '/board/' + this.$store.state.currentBoardFp + '/info',
        },
        {
          name: 'MOD ACTIVITY',
          link: '/board/' + this.$store.state.currentBoardFp + '/modactivity',
        },
        {
          name: 'ELECTIONS',
          link: '/board/' + this.$store.state.currentBoardFp + '/elections',
        },
      ]
      if (this.isMod) {
        return modTabsList
      }
      return nonModTabsList
    },
    entityNotFound(this: any) {
      return this.$store.state.currentBoard.fingerprint.length === 0
    },
    isMod(this: any) {
      if (
        this.$store.state.modModeEnabledArrived &&
        this.$store.state.modModeEnabled
      ) {
        return true
      }
      return false
    },
  },
  data(this: any) {
    return {}
  },
  methods: {
    tabsVisible: function (this: any) {
      if (this.$route.name === 'Board>NewThread') {
        return false
      } else {
        return true
      }
    },
  },
  mounted(this: any) {},
  updated(this: any) {},
}
</script>

<style lang="scss" scoped>
.spinner-container {
  display: flex;
  .spinner {
    margin: auto;
    padding-top: 50px;
  }
}
</style>
