<template>
  <div class="global-sublocation">
    <div class="subbed-root">
      <template v-if="!loadingComplete">
        <div class="spinner-container">
          <a-spinner></a-spinner>
        </div>
      </template>
      <template v-else>
        <a-board-listitem
          v-for="board in subbedBoards"
          :key="board.fingerprint"
          :board="board"
        ></a-board-listitem>
        <a-fin-puck></a-fin-puck>
      </template>
    </div>
  </div>
</template>
<script lang="ts">
var Vuex = require('../../../../../node_modules/vuex').default
export default {
  name: 'subbedroot',
  data() {
    return {}
  },
  computed: {
    ...Vuex.mapState(['allBoards', 'allBoardsLoadComplete']),
    loadingComplete(this: any) {
      return this.allBoardsLoadComplete
    },
    subbedBoards(this: any) {
      let subbed = []
      let vm = this
      for (var i = 0; i < this['allBoards'].length; i++) {
        ;(function(i) {
          if (vm.allBoards[i].subscribed) {
            subbed.push(vm.allBoards[i])
          }
        })(i)
      }
      return subbed
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
