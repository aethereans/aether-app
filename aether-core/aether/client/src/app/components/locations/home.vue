<template>
  <div class="location">
    <div class="home">
      <template v-if="!$store.state.homeViewArrived">
        <div class="spinner-container">
          <a-spinner :delay="300" :hidetext="true"></a-spinner>
        </div>
      </template>
      <template v-else>
        <template v-if="$store.state.onboardCompleteStatusArrived">
          <a-home-header></a-home-header>
          <a-bootstrapper></a-bootstrapper>
          <a-thread-listitem
            v-for="thr in bracketedHomeViewThreads"
            :thread="thr"
            :key="thr.Fingerprint"
          ></a-thread-listitem>
          <div
            class="load-more-carrier"
            v-show="
              !(loadMoreCaret + loadMoreBatchSize >= homeViewThreads.length)
            "
          >
            <a
              class="button is-warning is-outlined load-more-button"
              @click="loadMore()"
            >
              LOAD MORE
            </a>
          </div>
          <a-no-content
            no-content-text="Hey there! You're not subbed to any communities yet. <br><br>To see something here, you can sub to some by clicking <a href='#/globalscope'>here</a> after the first bootstrap is done."
            :quoteDisabled="true"
            v-if="homeViewThreads.length === 0"
          ></a-no-content>
          <a-fin-puck
            v-show="loadMoreCaret + loadMoreBatchSize >= homeViewThreads.length"
          ></a-fin-puck>
        </template>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
var Vuex = require('../../../../node_modules/vuex').default
export default {
  name: 'home',
  data() {
    return {
      breadcrumbState: [],
      loadMoreCaret: 0,
      loadMoreBatchSize: 25,
    }
  },
  methods: {
    ...Vuex.mapActions(['setBreadcrumbs']),
    loadMore(this: any) {
      if (
        this.loadMoreCaret + this.loadMoreBatchSize >=
        this.homeViewThreads.length
      ) {
        return
      }
      this.loadMoreCaret = this.loadMoreCaret + this.loadMoreBatchSize
    },
  },
  computed: {
    ...Vuex.mapState(['homeViewThreads']),
    bracketedHomeViewThreads(this: any) {
      return this.homeViewThreads.slice(
        0,
        this.loadMoreCaret + this.loadMoreBatchSize
      )
    },
  },
  beforeMount(this: any) {},
  mounted(this: any) {
    console.log('home mounted routine gets called')
    this.setBreadcrumbs()
  },
  updated(this: any) {
    console.log('home updated routine gets called')
    // this['setBreadcrumbs']([])
    this.setBreadcrumbs()
  },
}
</script>

<style lang="scss" scoped>
.location .home {
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
}
</style>
