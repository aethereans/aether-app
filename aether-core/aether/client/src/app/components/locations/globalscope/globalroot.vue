<template>
  <div class="global-sublocation">
    <div class="global-root">
      <template v-if="!loadingComplete">
        <div class="spinner-container">
          <a-spinner></a-spinner>
        </div>
      </template>
      <template v-else>
        <a-board-listitem
          v-for="board in bracketedSfwlistedBoards"
          :key="board.fingerprint"
          :board="board"
        ></a-board-listitem>
        <template
          v-if="
            bracketedSfwlistedBoards.length >=
            sfwlistedCaret + sfwlistedBatchSize
          "
        >
          <div class="load-more-carrier">
            <a
              class="button is-warning is-outlined load-more-button"
              @click="loadMoreSFWListed()"
            >
              LOAD MORE
            </a>
          </div>
        </template>
        <template v-else>
          <a-fin-puck></a-fin-puck>
        </template>
        <div
          class="non-sfwlist-info"
          v-show="
            !$store.state.ambientStatus.frontendambientstatus.sfwlistdisabled
          "
        >
          <div class="non-sfwlist-info-carrier">
            <p>
              <span class="nonsfw-header"
                >There are some communities marked as not safe-for-work.
              </span>
              <br />
              <br />You can view the list of non-SFW communities below. (New
              communities are considered non-SFW at first.)
              <router-link to="/sfwlist">More info</router-link>
            </p>
            <div class="show-non-sfw" v-if="!nonSFWListedVisible">
              <a @click="toggleNonSFWListed()">Show non-SFW communities</a>
              <br />
            </div>
          </div>
        </div>
        <template v-if="nonSFWListedVisible">
          <a-board-listitem
            v-for="board in bracketedNonSFWListedBoards"
            :key="board.fingerprint"
            :board="board"
          ></a-board-listitem>
          <template
            v-if="
              bracketedNonSFWListedBoards.length >=
              nonsfwlistedCaret + nonsfwlistedBatchSize
            "
          >
            <div class="load-more-carrier">
              <a
                class="button is-warning is-outlined load-more-button"
                @click="loadMoreNonSFWListed()"
              >
                LOAD MORE
              </a>
            </div>
          </template>
          <template v-else>
            <a-fin-puck></a-fin-puck>
          </template>
        </template>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
var Vuex = require('../../../../../node_modules/vuex').default
export default {
  name: 'globalroot',
  data() {
    return {
      nonSFWListedVisible: false,
      // nonSFWListedFirstLoadedItem: 0,
      // nonSFWListedLastLoadedItem: 25,
      // sfwlistedFirstLoadedItem: 0,
      // sfwlistedLastLoadedItem: 25,
      // batchSize: 25,
      moreInfoOpen: true,
      sfwlistedCaret: 0,
      sfwlistedBatchSize: 25,
      nonsfwlistedCaret: 0,
      nonsfwlistedBatchSize: 25,
    }
  },
  computed: {
    ...Vuex.mapState(['allBoards', 'allBoardsLoadComplete']),
    loadingComplete(this: any) {
      return this.allBoardsLoadComplete
    },
    bracketedSfwlistedBoards(this: any) {
      let sfwlisted = []
      let vm = this
      if (
        this.$store.state.ambientStatus.frontendambientstatus.sfwlistdisabled
      ) {
        // SFW list is disabled. All communities are sfw listed communities.
        sfwlisted = this.allBoards
      } else {
        for (var i = 0; i < this.allBoards.length; i++) {
          ;(function (i) {
            if (vm.allBoards[i].sfwlisted) {
              sfwlisted.push(vm.allBoards[i])
            }
          })(i)
        }
      }
      return sfwlisted.slice(0, this.sfwlistedCaret + this.sfwlistedBatchSize)
    },
    bracketedNonSFWListedBoards(this: any) {
      let nonSFWListed = []
      let vm = this
      for (var i = 0; i < this['allBoards'].length; i++) {
        ;(function (i) {
          if (!vm.allBoards[i].sfwlisted) {
            nonSFWListed.push(vm.allBoards[i])
          }
        })(i)
      }
      return nonSFWListed.slice(
        0,
        this.nonsfwlistedCaret + this.nonsfwlistedBatchSize
      )
    },
  },
  methods: {
    toggleNonSFWListed(this: any) {
      this.nonSFWListedVisible = !this.nonSFWListedVisible
    },
    loadMoreNonSFWListed(this: any) {
      // this.nonSFWListedLastLoadedItem = this.nonSFWListedLastLoadedItem + this.batchSize
      this.nonsfwlistedCaret =
        this.nonsfwlistedCaret + this.nonsfwlistedBatchSize
    },
    loadMoreSFWListed(this: any) {
      // this.sfwlistedLastLoadedItem = this.sfwlistedLastLoadedItem + this.batchSize
      this.sfwlistedCaret = this.sfwlistedCaret + this.sfwlistedBatchSize
    },
    toggleMoreInfo(this: any) {
      this.moreInfoOpen
        ? (this.moreInfoOpen = false)
        : (this.moreInfoOpen = true)
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../../../scss/globals';
.non-sfwlist-info {
  font-family: 'SCP Regular';
  margin: 20px;
  padding: 20px;
  background-color: rgba(0, 0, 0, 0.25);
  border-radius: 3px;
  a {
    font-family: 'SCP Bold';
    &:hover {
      color: $a-grey-800;
    }
  }
  .bold {
    font-family: 'SCP Bold';
  }
  .header {
    font-family: 'SCP Bold';
    font-size: 150%;
    margin-bottom: 15px;
    margin-top: 50px;
  }
  .non-sfwlist-info-carrier {
    font-family: 'SCP Semibold';
    max-width: 750px;
    .nonsfw-header {
      // font-family: 'SCP Bold';
    }
  }
}

.button {
  font-family: 'SSP Bold';
}

.load-more-carrier {
  display: flex;
  padding: 20px 0 0 0;
  .load-more-button {
    margin: auto;
  }
}

.spinner-container {
  display: flex;
  .spinner {
    margin: auto;
    padding-top: 50px;
  }
}
</style>
