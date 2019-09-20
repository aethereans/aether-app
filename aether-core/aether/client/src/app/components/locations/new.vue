<template>
  <div class="location">
    <div class="new">
      <template v-if="!$store.state.newViewArrived">
        <div class="spinner-container">
          <a-spinner :delay="300" :hidetext="true"></a-spinner>
        </div>
      </template>
      <template v-else>
        <a-new-header></a-new-header>
        <a-bootstrapper></a-bootstrapper>
        <template v-for="entity in bracketedNewList">
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
        </template>
        <div
          class="load-more-carrier"
          v-show="
            !(loadMoreCaret + loadMoreBatchSize >= mergedNewList.length)
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
          no-content-text="Your new list is empty. This usually happens when the app is bootstrapping itself. When this process is complete, the latest content will appear here. <br><br>In the meanwhile, make yourself comfortable, and feel free to take a look&nbsp;around."
          :quoteDisabled="true"
          v-if="mergedNewList.length === 0"
        ></a-no-content>
        <a-fin-puck
          v-show="
            loadMoreCaret + loadMoreBatchSize >= mergedNewList.length
          "
        ></a-fin-puck>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
var Vuex = require('../../../../node_modules/vuex').default
export default {
  name: 'new',
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
        this.mergedNewList.length
      ) {
        return
      }
      this.loadMoreCaret = this.loadMoreCaret + this.loadMoreBatchSize
    },
  },
  computed: {
    ...Vuex.mapState(['newView']),
    mergedNewList(this: any) {
      return this.newView.postsList
        .concat(this.newView.threadsList)
        .sort((a: any, b: any) => b.creation - a.creation)
    },
    bracketedNewList(this: any) {
      return this.mergedNewList.slice(
        0,
        this.loadMoreCaret + this.loadMoreBatchSize
      )
    },
  },
  beforeMount(this: any) {},
  mounted(this: any) {
    console.log('new mounted routine gets called')
    // this.setBreadcrumbs()
  },
  updated(this: any) {
    console.log('new updated routine gets called')
    // this.setBreadcrumbs()
  },
}
</script>

<style lang="scss" scoped>
.location .new {
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
