<template>
  <div class="board-sublocation">
    <div class="board-root">
      <template v-if="$store.state.route.name === 'Board>ThreadsNewList'">
        <a-thread-listitem
          v-for="thr in inflightNewThreads.slice().reverse()"
          :thread="thr.entity"
          :inflightStatus="thr.status"
          :key="thr.Fingerprint"
        ></a-thread-listitem>
      </template>
      <a-thread-listitem
        v-for="thr in bracketedCurrentBoardsThreads"
        :thread="thr"
        :key="thr.Fingerprint"
      ></a-thread-listitem>
      <a-no-content
        no-content-text="There are no threads yet. You should write something."
        v-if="hasNoContent"
      ></a-no-content>
      <div
        class="load-more-carrier"
        v-show="
          !(loadMoreCaret + loadMoreBatchSize >= currentBoardsThreads.length)
        "
      >
        <a
          class="button is-warning is-outlined load-more-button"
          @click="loadMore()"
        >
          LOAD MORE
        </a>
      </div>
      <a-fin-puck
        v-show="
          loadMoreCaret + loadMoreBatchSize >= currentBoardsThreads.length
        "
      ></a-fin-puck>
    </div>
  </div>
</template>

<script lang="ts">
var Vuex = require('../../../../../node_modules/vuex').default
var Mousetrap = require('../../../../../node_modules/mousetrap')
export default {
  name: 'boardroot',
  data() {
    return {
      currentBoardsThreadsNew: [],
      loadMoreCaret: 0,
      loadMoreBatchSize: 25,
    }
  },
  computed: {
    ...Vuex.mapState(['currentBoardsThreads']),
    inflightNewThreads(this: any) {
      let inflightNewThreads = []
      for (let val of this.$store.state.ambientStatus.inflights.threadsList) {
        if (val.status.eventtype !== 'CREATE') {
          continue
        }
        if (this.$store.state.currentBoard.fingerprint !== val.entity.board) {
          continue
        }
        inflightNewThreads.push(val)
      }
      console.log('inflight new threads')
      console.log(inflightNewThreads)
      return inflightNewThreads
    },
    hasNoContent(this: any) {
      if (
        this.$store.state.route.name === 'Board' ||
        this.$store.state.route.name === undefined
      ) {
        if (this.currentBoardsThreads.length > 0) {
          return false
        }
        return true
      }
      // If Board>ThreadsNewList
      if (this.inflightNewThreads.length > 0) {
        return false
      }
      if (this.currentBoardsThreads.length > 0) {
        return false
      }
      return true
    },
    bracketedCurrentBoardsThreads(this: any) {
      return this.currentBoardsThreads.slice(
        0,
        this.loadMoreCaret + this.loadMoreBatchSize
      )
    },
  },
  methods: {
    loadMore(this: any) {
      if (
        this.loadMoreCaret + this.loadMoreBatchSize >=
        this.currentBoardsThreads.length
      ) {
        return
      }
      this.loadMoreCaret = this.loadMoreCaret + this.loadMoreBatchSize
    },
  },
  mounted(this: any) {
    // console.log('this is currentboardsthreads')
    // console.log(this.currentBoardsThreads)
    let vm = this
    Mousetrap.bind('mod+n', function() {
      vm.$router.push(
        '/board/' + vm.$store.state.currentBoard.fingerprint + '/newthread'
      )
    })
  },
  beforeDestroy(this: any) {
    Mousetrap.unbind('mod+n')
  },
  updated() {},
}
</script>

<style lang="scss" scoped>
.load-more-carrier {
  display: flex;
  padding: 20px 0 0 0;
  .load-more-button {
    margin: auto;
  }
}
</style>
