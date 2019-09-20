<template>
  <div class="user-sublocation" v-if="$store.state.notificationsArrived">
    <div class="user-notifications">
      <a-notification-entity
        v-for="n in bracketedNotifications"
        :key="n.link"
        :notification="n"
      ></a-notification-entity>
      <div class="buttons-carrier">
        <div class="buttons-container">
          <a
            class="button is-warning is-outlined load-more-button"
            @click="loadMore()"
            v-show="
              !(loadMoreCaret + loadMoreBatchSize >= notifications.length)
            "
          >
            LOAD MORE
          </a>
          <a
            class="button is-warning is-outlined markallasread-button"
            :class="{ 'is-loading': markAllAsReadInProgress }"
            @click="markAllAsRead()"
            v-show="$store.state.unseenNotificationsPresent"
          >
            MARK ALL AS READ
          </a>
        </div>
      </div>
      <a-no-content
        no-content-text="You have no notifications."
        v-if="notifications.length === 0"
      ></a-no-content>
      <a-fin-puck
        v-show="loadMoreCaret + loadMoreBatchSize >= notifications.length"
      ></a-fin-puck>
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
      markAllAsReadInProgress: false,
    }
  },
  computed: {
    ...Vuex.mapState(['notifications']),
    bracketedNotifications(this: any) {
      return this.notifications.slice(
        0,
        this.loadMoreCaret + this.loadMoreBatchSize
      )
    },
  },
  methods: {
    loadMore(this: any) {
      if (
        this.loadMoreCaret + this.loadMoreBatchSize >=
        this.notifications.length
      ) {
        return
      }
      this.loadMoreCaret = this.loadMoreCaret + this.loadMoreBatchSize
    },
    markAllAsRead(this: any) {
      let vm = this
      vm.markAllAsReadInProgress = true
      fe.markAllAsRead(function() {
        vm.markAllAsReadInProgress = false
      })
    },
  },
  mounted(this: any) {
    fe.markSeen()
  },
  updated() {},
}
</script>

<style lang="scss" scoped>
.buttons-carrier {
  display: flex;
  padding: 20px 0 0 0;
  .buttons-container {
    margin: auto;
  }
  .load-more-button,
  .markallasread-button {
    // margin: auto;
    margin: 0 5px;
  }
}
</style>
