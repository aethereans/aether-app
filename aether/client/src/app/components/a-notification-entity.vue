<template>
  <div @click="clickNotification" class="notification-entity" :class="{'unread': !notification.read}">
    <template v-if="notificationType === 'thread'">
      <div class="icon-container thread">
        <icon>
          <icon class="under-icon" name="circle" scale="1.5"></icon>
          <icon class="base-icon" name="asterisk"></icon>
        </icon>
      </div>
    </template>
    <template v-else>
      <div class="icon-container post">
        <icon>
          <icon class="under-icon" name="circle" scale="1.5"></icon>
          <icon class="base-icon" name="ellipsis-h"></icon>
        </icon>
      </div>
    </template>
    <div class="text-container">
      <template v-if="notification.responsepostsusersList.length > 0">
        <a-username styling="inline" :notificationowner="notification.responsepostsusersList[0]"></a-username>
      </template>
      {{notification.text}}
    </div>
    <div class="flex-spacer"></div>
    <div class="time-container">
      <a-timestamp :creation="notification.creationtimestamp"></a-timestamp>
    </div>
  </div>
</template>

<script lang="ts">
  // var globalMethods = require('../services/globals/methods')
  var fe = require('../services/feapiconsumer/feapiconsumer')
  // var vuexStore = require('../store/index').default
  export default {
    name: 'a-notification-entity',
    props: {
      notification: {
        type: Object,
        default: {}
      }
    },
    data() {
      return {}
    },
    computed: {
      notificationType(this: any) {
        if (this.notification.type === 1) {
          return "thread"
        }
        return "post"
      },
      parentFingerprint(this: any) {
        if (this.notificationType === "thread") {
          return this.notification.parentthread.fingerprint
        }
        return this.notification.parentpost.fingerprint
      },
      link(this: any) {
        if (this.notificationType === "thread") {
          let b = this.notification.parentthread.board
          let t = this.notification.parentthread.fingerprint
          let highlightSelectors: any = []
          for (let val of this.notification.responsepostsList) {
            highlightSelectors.push('post-' + val)
          }
          // let link = "/board/" + b + "/thread/" + t + '?highlightSelectors=' + JSON.stringify(this.notification.responsepostsList)
          let link = "/board/" + b + "/thread/" + t + '?focusSelector=' + this.focusSelector + '&parentSelector=' + this.focusSelector + '&highlightSelectors=' + JSON.stringify(this.notification.responsepostsList)
          return link
        }
        // post
        let b = this.notification.parentpost.board
        let t = this.notification.parentpost.thread
        // let p = "post-" + this.notification.parentpost.fingerprint
        let highlightSelectors: any = []
        for (let val of this.notification.responsepostsList) {
          highlightSelectors.push('post-' + val)
        }
        let link = "/board/" + b + "/thread/" + t + '?focusSelector=' + this.focusSelector + '&parentSelector=' + this.focusSelector + '&highlightSelectors=' + JSON.stringify(this.notification.responsepostsList)
        return link
      },
      focusSelector(this: any) {
        if (this.notificationType === "thread") {
          // return ""
          /* Show first response */
          return this.notification.responsepostsList[0]

        }
        // post
        // return "post-" + this.notification.parentpost.fingerprint
        return this.notification.parentpost.fingerprint
      }
    },
    methods: {
      clickNotification(this: any) {
        fe.markRead(this.parentFingerprint)
        console.log('notification link')
        console.log(this.link)
        this.$router.push(this.link)
      }
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .notification-entity {
    display: flex;
    flex-direction: row;
    padding: 10px 15px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    font-family: "SSP Regular";
    color: $a-grey-500;
    cursor: pointer;
    &:hover {
      background-color: rgba(255, 255, 255, 0.05);
    }
    .icon-container {
      display: flex;
      margin-right: 10px;
      opacity: 0.6;
      svg {
        margin: auto;
        color: $a-grey-800;
      }
      &.thread .under-icon {
        fill: $a-grey-300;
      }
      &.post .under-icon {
        fill: $a-grey-300;
      }
      &.post .base-icon {
        -webkit-filter: drop-shadow( -5px -5px 5px #000);
      }
    }
    .time-container {
      .timestamp {
        color: $a-grey-500;
        text-align: right;
      }
    }
    &.unread {
      font-family: "SSP Bold";
      color: $a-grey-800;
      .icon-container {
        opacity: 1;
        &.thread .under-icon {
          fill: $a-earth; // fill: $a-orange;
        }
        &.post .under-icon {
          // fill: $a-purple-80;
          fill: $a-earth;
        }
      }
      .time-container {
        .timestamp {
          color: $a-grey-800;
        }
      }
    }
  }

  .flex-spacer {
    flex: 1;
  }
</style>