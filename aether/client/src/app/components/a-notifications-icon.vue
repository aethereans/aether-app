<template>
  <div
    class="notifications-box"
    v-if="$store.state.notificationsArrived"
    title="Notifications"
    @click="$router.push(notificationsLink)"
  >
    <a-header-icon
      class="notifications-icon"
      v-show="$store.state.unseenNotificationsPresent"
      icon="regular/bell"
      iconunder="circle"
    ></a-header-icon>
    <a-header-icon
      class="notifications-icon"
      v-show="!$store.state.unseenNotificationsPresent"
      icon="regular/bell"
    ></a-header-icon>
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-notifications-icon',
  data() {
    return {}
  },
  computed: {
    notificationsLink(this: any) {
      return (
        '/user/' + this.$store.state.localUser.fingerprint + '/notifications'
      )
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.notifications-icon {
  color: $a-grey-800;
}
</style>
