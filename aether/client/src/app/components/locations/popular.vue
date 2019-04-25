<template>
  <div class="location">
    <div class="popular">
      <template v-if="!$store.state.popularViewArrived">
        <div class="spinner-container">
          <a-spinner :delay="300" :hidetext="true"></a-spinner>
        </div>
      </template>
      <template v-else>
        <a-popular-header></a-popular-header>
        <a-bootstrapper></a-bootstrapper>
        <a-thread-listitem
          v-for="thr in bracketedPopularViewThreads"
          :thread="thr"
          :key="thr.Fingerprint"
        ></a-thread-listitem>
        <div
          class="load-more-carrier"
          v-show="
            !(loadMoreCaret + loadMoreBatchSize >= popularViewThreads.length)
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
          no-content-text="Your popular list is empty. This usually happens when the app is bootstrapping itself. When this process is complete, the latest content will appear here. <br><br>In the meanwhile, make yourself comfortable, and feel free to take a look&nbsp;around."
          :quoteDisabled="true"
          v-if="popularViewThreads.length === 0"
        ></a-no-content>
        <a-fin-puck
          v-show="
            loadMoreCaret + loadMoreBatchSize >= popularViewThreads.length
          "
        ></a-fin-puck>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
var Vuex = require('../../../../node_modules/vuex').default
export default {
  name: 'popular',
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
        this.popularViewThreads.length
      ) {
        return
      }
      this.loadMoreCaret = this.loadMoreCaret + this.loadMoreBatchSize
    },
  },
  computed: {
    ...Vuex.mapState(['popularViewThreads']),
    bracketedPopularViewThreads(this: any) {
      return this.popularViewThreads.slice(
        0,
        this.loadMoreCaret + this.loadMoreBatchSize
      )
    },
  },
  beforeMount(this: any) {},
  mounted(this: any) {
    console.log('popular mounted routine gets called')
    // this.setBreadcrumbs()
  },
  updated(this: any) {
    console.log('popular updated routine gets called')
    // this.setBreadcrumbs()
  },
}
</script>

<style lang="scss" scoped>
.location .popular {
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
