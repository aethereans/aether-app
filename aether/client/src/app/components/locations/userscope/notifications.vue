<template>
  <div class="user-sublocation" v-if="$store.state.notificationsArrived">
    <div class="user-notifications">
      <a-notification-entity v-for="n in bracketedNotifications" :notification="n"></a-notification-entity>
      <div class="load-more-carrier" v-show="!(loadMoreCaret + loadMoreBatchSize >= notifications.length)">
        <a class="button is-warning is-outlined load-more-button" @click="loadMore()">
      LOAD MORE
            </a>
      </div>
      <a-no-content no-content-text="You have no notifications." v-if="notifications.length === 0"></a-no-content>
      <a-fin-puck v-show="(loadMoreCaret + loadMoreBatchSize >= notifications.length)"></a-fin-puck>
    </div>
  </div>
</template>

<script lang="ts">
  var Vuex = require('../../../../../node_modules/vuex').default
  var fe = require('../../../services/feapiconsumer/feapiconsumer')
  export default {
    name: 'notifications',
    data() {
      return {
        loadMoreCaret: 0,
        loadMoreBatchSize: 50,
      }
    },
    computed: {
      ...Vuex.mapState(['notifications']),
      bracketedNotifications(this: any) {
        return this.notifications.slice(0, this.loadMoreCaret + this.loadMoreBatchSize)
      }
    },
    methods: {
      loadMore(this: any) {
        if (this.loadMoreCaret + this.loadMoreBatchSize >= this.notifications.length) {
          return
        }
        this.loadMoreCaret = this.loadMoreCaret + this.loadMoreBatchSize
      }
    },
    mounted(this: any) {
      fe.markSeen()
    },
    updated() {}
  }
</script>

<style lang="scss" scoped>
  @import "../../../scss/bulmastyles";
  .load-more-carrier {
    display: flex;
    padding: 20px 0 0 0;
    .load-more-button {
      margin: auto;
    }
  }
</style>