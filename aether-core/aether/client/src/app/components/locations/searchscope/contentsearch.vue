<template>
  <div class="search-sublocation">
    <div class="search-user">
      <template v-if="!loadingComplete">
        <div class="spinner-container">
          <a-spinner></a-spinner>
        </div>
      </template>
      <template v-else>
        <!-- Load complete here -->
        <div class="search-form-container">
          <a-composer id="searchform" :spec="searchSpec"></a-composer>
        </div>
        <template v-if="stmState === 'idle'">
          <a-no-content
            no-content-text="Enter your search term above to start."
          ></a-no-content>
        </template>
        <template v-if="stmState === 'searchInProgress'">
          <div class="spinner-container">
            <a-spinner></a-spinner>
          </div>
        </template>
        <template v-if="stmState === 'searchComplete'">
          <template v-if="primaryList.length > 0">
            <div class="results-container">
              <div class="list-item-container" v-for="entity in primaryList">
                <template v-if="typeof entity.parent === 'undefined'">
                  <a-thread-listitem
                    :thread="entity"
                    :key="entity.Fingerprint"
                  ></a-thread-listitem>
                </template>
                <template v-else>
                  <a-post-listitem :post="entity"></a-post-listitem>
                  <!-- Post -->
                </template>
              </div>
              <template v-if="primaryListLoadMoreVisible">
                <div class="load-more-carrier">
                  <a
                    class="button is-warning is-outlined load-more-button"
                    @click="primaryListLoadMore"
                  >
                    LOAD MORE
                  </a>
                </div>
              </template>
            </div>
          </template>
          <template v-else>
            <a-no-content
              :no-content-text="
                'No results found for \'' +
                currentSearchTerm +
                '\'. A shorter query might help.'
              "
              :quote-disabled="true"
            ></a-no-content>
          </template>
          <!-- Secondary list: NSFW communities -->
          <template v-if="secondaryListNotEmpty">
            <div
              class="non-sfwlist-info"
              v-show="
                !$store.state.ambientStatus.frontendambientstatus
                  .sfwlistdisabled
              "
            >
              <div class="non-sfwlist-info-carrier">
                <p>
                  <span class="nonsfw-header"
                    >There are some results marked as not safe-for-work.</span
                  >
                  <br />
                  <br />You can view the non-SFW search results below. (New
                  communities are considered non-SFW at first.)
                  <router-link to="/sfwlist">More info</router-link>
                </p>
                <div class="show-non-sfw" v-if="!secondaryListVisible">
                  <a @click="toggleSecondaryListVisible"
                    >Show non-SFW results</a
                  >
                  <br />
                </div>
              </div>
            </div>
            <template v-if="secondaryListVisible">
              <div class="list-item-container" v-for="entity in secondaryList">
                <template v-if="typeof entity.parent === 'undefined'">
                  <a-thread-listitem
                    :thread="entity"
                    :key="entity.Fingerprint"
                  ></a-thread-listitem>
                </template>
                <template v-else>
                  <a-post-listitem :post="entity"></a-post-listitem>
                  <!-- Post -->
                </template>
              </div>
              <template v-if="secondaryListLoadMoreVisible">
                <div class="load-more-carrier">
                  <a
                    class="button is-warning is-outlined load-more-button"
                    @click="secondaryListLoadMore"
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
          <template
            v-if="
              !primaryListLoadMoreVisible &&
              !secondaryListNotEmpty &&
              !secondaryListVisible
            "
          >
            <!-- Primary list fin puck - only visible if there's nothing in the secondary list. -->
            <a-fin-puck></a-fin-puck>
          </template>
        </template>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
/*
  This initialises a finite state machine.
  https://brilliant.org/wiki/finite-state-machines/

  This implementation below seems to be the best way to use this with Vue. Vue entities do not need access to the machine itself. I tried putting it in data, but since it's a function call, it takes some time to init, and other parts of the component consider it immediately available, which creates a race condition. Instead of that, I've moved it outside the component and I'm initialising this in sequential order, and then fire the whole thing on component creation.

  Why do we even need that?

  Because we have a lot of implicit state in our components. Loading states, talking to API conditions, UI objects and so on. We have been dealing with this with 'toggles' in data until now, however, the issue is that those toggles are not only non-standard, they are also very easy to mess up because you missed a state or a subtle state transition happens in a place that you don't actually expect. Formalising these transitions into state machines (thus making it obvious which states can transition to others, and emitting big red flashing errors when a state tries to move to one not allowed from that one) helps us clarify the code in a manner that is not only graphable visually.

  See visualiser here:
  https://statecharts.github.io/xstate-viz/

*/

var xstate = require('xstate')
var feapiconsumer = require('../../../services/feapiconsumer/feapiconsumer')

function initialiseSTM(actions: Object, context: Object) {
  let machine: any,
    service: any,
    initialState: any = {}
  machine = xstate.Machine(
    {
      id: 'contentsearch',
      context: context,
      initial: 'idle',
      states: {
        idle: {
          on: {
            SEARCH_STARTED: {
              target: 'searchInProgress',
              actions: ['search'],
            },
          },
        },
        searchInProgress: {
          on: {
            RESULT_ARRIVED: {
              target: 'searchComplete',
              actions: ['resetCarets'],
            },
            // If result doesn't return in 30 seconds, timeout TODO
          },
        },
        searchComplete: {
          on: {
            SEARCH_STARTED: {
              target: 'searchInProgress',
              actions: ['search'],
            },
            INCR_PRIMARY_LIST_CARET: {
              actions: ['incrPrimaryCaret'],
            },
            INCR_SECONDARY_LIST_CARET: {
              actions: ['incrSecondaryCaret'],
            },
          },
        },
      },
    },
    {
      actions: actions,
    }
  )
  service = xstate.interpret(machine)
  initialState = machine.initial
  return {
    machine: machine,
    service: service,
    initialState: initialState,
  }
}
export default {
  name: 'contentsearch',
  data(this: any) {
    return {
      loadingComplete: true,
      currentSearchTerm: '',
      secondaryListVisible: false,
      searchSpec: {
        fields: [
          {
            id: 'searchTerm',
            visibleName: '',
            description: ``,
            placeholder: 'Search content',
            maxCharCount: 100,
            minCharCount: 1,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: false,
            spaceDisabled: false,
            preInfoEnabled: false,
            preInfoText: '',
            searchField: true,
          },
        ],
        commitActionName: 'SEARCH',
        commitAction: this.submitSearch,
        cancelActionName: '',
        cancelAction: function () {},
        fixToBottom: true,
        autofocus: true,
        preventClearAfterSuccessfulCommit: true,
      },
      stmService: {},
      stmCurrState: {},
      stmActions: {
        search: function (_: any, event: any) {
          // _ = context: any, comes from machine context defined above.
          feapiconsumer.SendContentSearchRequest(
            event.vars.searchTerm,
            function () {
              console.log('successfully sent')
            }
          )
        },
        incrPrimaryCaret(context: any, _: any) {
          context.primaryListCaret += context.primaryListBatchSize
        },
        incrSecondaryCaret(context: any, _: any) {
          context.secondaryListCaret += context.secondaryListBatchSize
        },
        resetCarets(context: any, _: any) {
          console.log('reset carets run')
          context.primaryListCaret = 0
          context.secondaryListCaret = 0
        },
      },
      stmContext: {
        primaryListCaret: 0,
        primaryListBatchSize: 25,
        secondaryListCaret: 0,
        secondaryListBatchSize: 25,
      },
      unsubFromMutationFunc: function () {},
    }
  },
  created(this: any) {
    let vm = this
    let stm = initialiseSTM(this.stmActions, this.stmContext)
    this.stmService = stm.service
    this.stmCurrState = stm.initialState
    this.stmService
      .onTransition(function (state: any) {
        vm.stmCurrState = state
      })
      .start()
  },
  beforeMount(this: any) {
    /*
      This subscribes to the save search result event, and whenever we get something like it, it'll fire a state change. If we're in a state that expects it, we'll know we've received our results.

      We also want to unsubscribe whenever this component gets destroyed so as to not create a memory leak.
    */
    this.unsubFromMutationFunc = this.$store.subscribe(
      (mutation: any, _: any) => {
        // _ : state
        if (mutation.type !== 'SAVE_SEARCH_RESULT') {
          return
        }
        // We have a response.
        this.sendEvent({ type: 'RESULT_ARRIVED', vars: {} })
      }
    )
  },
  mounted(this: any) {
    // We look at the search query - if there's a query present, we need to make the search API call for it, so we can show refreshed results for that query. This effectively is only useful when the user is navigating with back/forward, allows us to capture the searched term in history from the query string and do a re-search.
    if (typeof this.$store.state.route.query.searchHashtag !== 'undefined') {
      this.$store.state.route.query.searchTerm =
        '#' + this.$store.state.route.query.searchHashtag
    }
    if (
      typeof this.$store.state.route.query !== 'undefined' &&
      typeof this.$store.state.route.query.searchTerm !== 'undefined'
    ) {
      this.sendEvent({
        type: 'SEARCH_STARTED',
        vars: { searchTerm: this.$store.state.route.query.searchTerm },
      })
      this.searchSpec.fields[0].content =
        this.$store.state.route.query.searchTerm
      // this.searchSpec.fields[0]._touched = true
      // ^ Not necessary, since if it was not valid the last time around to submit, it would not have gotten into here.
    }
  },
  beforeDestroy(this: any) {
    this.stmService.stop()
    this.unsubFromMutationFunc()
  },
  computed: {
    stmState(this: any) {
      return this.stmCurrState.value
    },
    mergedRawContentList(this: any) {
      return this.$store.state.postsSearchResult
        .concat(this.$store.state.threadsSearchResult)
        .sort((a: any, b: any) => a.viewmetaSearchscore - b.viewmetaSearchscore)
    },
    primaryList(this: any) {
      // return this.$store.state.threadsSearchResult.slice(
      //   0,
      //   this.stmContext.primaryListCaret + this.stmContext.primaryListBatchSize
      // )

      let primaryList = []
      let vm = this
      if (
        this.$store.state.ambientStatus.frontendambientstatus.sfwlistdisabled
      ) {
        // SFW list is disabled. All communities are sfw listed communities.
        primaryList = this.mergedRawContentList
      } else {
        for (var i = 0; i < this.mergedRawContentList.length; i++) {
          ;(function (i) {
            if (vm.mergedRawContentList[i].compiledcontentsignals.modblocked) {
              return
              // If modblocked by a valid mod that the user considers valid, don't show. (Heads up, we don't have modblocked for boards themselves yet.)
            }
            if (vm.mergedRawContentList[i].viewmetaSfwlisted) {
              primaryList.push(vm.mergedRawContentList[i])
            }
          })(i)
        }
      }
      return primaryList.slice(
        0,
        this.stmContext.primaryListCaret + this.stmContext.primaryListBatchSize
      )
    },
    primaryListLoadMoreVisible(this: any) {
      return (
        this.primaryList.length >=
        this.stmContext.primaryListCaret + this.stmContext.primaryListBatchSize
      )
    },
    secondaryList(this: any) {
      // return this.$store.state.threadsSearchResult.slice(
      //   0,
      //   this.stmContext.primaryListCaret + this.stmContext.primaryListBatchSize
      // )

      let secondaryList = []
      let vm = this
      if (
        this.$store.state.ambientStatus.frontendambientstatus.sfwlistdisabled
      ) {
        // SFW list is disabled. All communities are sfw listed communities.
        secondaryList = this.mergedRawContentList
      } else {
        for (var i = 0; i < this.mergedRawContentList.length; i++) {
          ;(function (i) {
            if (vm.mergedRawContentList[i].compiledcontentsignals.modblocked) {
              return
              // If modblocked by a valid mod that the user considers valid, don't show. (Heads up, we don't have modblocked for boards themselves yet.)
            }
            if (!vm.mergedRawContentList[i].viewmetaSfwlisted) {
              secondaryList.push(vm.mergedRawContentList[i])
            }
          })(i)
        }
      }
      return secondaryList.slice(
        0,
        this.stmContext.secondaryListCaret +
          this.stmContext.secondaryListBatchSize
      )
    },
    secondaryListLoadMoreVisible(this: any) {
      return (
        this.secondaryList.length >=
        this.stmContext.secondaryListCaret +
          this.stmContext.secondaryListBatchSize
      )
    },
    secondaryListNotEmpty(this: any) {
      // This is useful because we want to show 'this is the end of the sfw results' only when there are actual nsfw results.
      return this.secondaryList.length > 0
    },
  },
  methods: {
    sendEvent(this: any, event: any) {
      this.stmService.send(event)
    },
    submitSearch(this: any, fields: any) {
      if (fields[0].id !== 'searchTerm') {
        return
      }
      let st = fields[0].content
      this.sendEvent({ type: 'SEARCH_STARTED', vars: { searchTerm: st } })
      this.currentSearchTerm = fields[0].content
      this.$router.replace({ query: { searchTerm: st } })
    },
    primaryListLoadMore(this: any) {
      this.sendEvent({ type: 'INCR_PRIMARY_LIST_CARET' })
    },
    secondaryListLoadMore(this: any) {
      this.sendEvent({ type: 'INCR_SECONDARY_LIST_CARET' })
    },
    toggleSecondaryListVisible(this: any) {
      this.secondaryListVisible = this.secondaryListVisible ? false : true
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../../../scss/globals';
.spinner-container {
  display: flex;
  .spinner {
    margin: auto;
    padding-top: 50px;
  }
}
.search-form-container {
  // padding: 5px 5px 0px 10px;
  // border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding: 15px 20px 0px 20px;
}

.load-more-carrier {
  display: flex;
  padding: 20px 0 0 0;
  .load-more-button {
    margin: auto;
  }
}

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
</style>
