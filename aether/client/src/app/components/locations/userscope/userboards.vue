<template>
  <div class="user-sublocation">
    <div class="user-boards">
      <template v-if="!loadingComplete">
        <div class="spinner-container">
          <a-spinner></a-spinner>
        </div>
      </template>
      <template v-else>
        <template v-if="isSelf">
          <template v-for="board in inflightCreates.slice().reverse()">
            <a-board-entity :inflightStatus="board.status" :board="board.entity" uncompiled="true" :refresher="fetchData"></a-board-entity>
            <div class="divider"></div>
          </template>
        </template>
        <div v-for="board in boardsList" :key="board.Fingerprint">
          <a-board-entity :board="board" uncompiled="true"></a-board-entity>
          <div class="divider"></div>
        </div>
        <div class="load-more-carrier" v-show="loadMoreVisible">
          <a class="button is-warning is-outlined load-more-button" @click="loadMore">
        LOAD MORE
              </a>
        </div>
        <a-no-content no-content-text="No communities created in retained history." quoteDisabled="true" v-if="boardsList.length === 0 && inflightCreates.length === 0">
        </a-no-content>
        <a-fin-puck v-show="!loadMoreVisible"></a-fin-puck>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
  var globalMethods = require('../../../services/globals/methods')
  var fe = require('../../../services/feapiconsumer/feapiconsumer')
  var Vuex = require('../../../../../node_modules/vuex').default
  export default {
    name: 'userboards',
    data() {
      return {
        boardsList: [],
        currentUserFp: '',
        loadingComplete: false,
        limit: 25,
        offset: 0,
        lastLoadSize: 0,
      }
    },
    computed: {
      ...Vuex.mapState(['currentUserEntity']),
      /*----------  Inflight computeds  ----------*/
      inflightCreates(this: any) {
        let inflightCreates = []
        for (let val of this.$store.state.ambientStatus.inflights.boardsList) {
          if (val.status.eventtype !== 'CREATE') {
            continue
          }
          inflightCreates.push(val)
        }
        return inflightCreates
      },
      isSelf(this: any) {
        if (globalMethods.IsUndefined(this.$store.state.currentUserEntity) || globalMethods.IsUndefined(this.$store.state.localUser)) {
          return false
        }
        if (this.$store.state.currentUserEntity.fingerprint !== this.$store.state.localUser.fingerprint) {
          return false
        }
        return true
      },
      loadMoreVisible(this: any) {
        if (this.lastLoadSize < this.limit) {
          return false
        }
        return true
      }
    },
    methods: {
      fetchData(this: any, targetuserfp: string) {
        let vm = this
        fe.GetUncompiledEntityByKey('Board', targetuserfp, 0, 0, function(resp: any) {
          console.log(resp)
          vm.boardsList = resp.boardsList
          vm.currentUserFp = targetuserfp
          vm.loadingComplete = true
          vm.lastLoadSize = resp.postsList.length
        })
      },
      loadMore(this: any) {
        this.offset = this.offset + this.limit
        this.fetchData(this.currentUserFp)
      }
    },
    beforeMount(this: any) {
      if (typeof this.currentUserEntity === 'undefined') {
        return
      }
      this.fetchData(this.currentUserEntity.fingerprint)
    },
    updated(this: any) {
      if (typeof this.currentUserEntity === 'undefined') {
        return
      }
      if (this.currentUserEntity.fingerprint === this.currentUserFp) {
        return
      }
      this.fetchData(this.currentUserEntity.fingerprint)
    }
  }
</script>

<style lang="scss" scoped>
  .user-sublocation {
    .user-boards {}
  }

  .divider {
    width: 100%;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .spinner-container {
    display: flex;
    .spinner {
      margin: auto;
      padding-top: 50px;
    }
  }

  .load-more-carrier {
    display: flex;
    padding: 20px 0 0 0;
    .load-more-button {
      margin: auto;
    }
  }

  .board-entity {
    border-bottom: none;
  }
</style>