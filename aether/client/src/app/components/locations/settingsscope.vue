<template>
  <div class="location settings-location" id="scrolltarget-container-target">
    <div class="settingsscope">
      <div class="settings-navigation">
        <router-link class="nav-item" to="/settings">Preferences</router-link>
        <router-link class="nav-item" to="/settings/defaults"
          >Defaults</router-link
        >
        <router-link class="nav-item" to="/modship">Mod mode</router-link>
        <router-link class="nav-item" to="/membership"
          >Supporter benefits</router-link
        >
        <router-link class="nav-item" to="/status">App status</router-link>
        <div class="spacer"></div>
        <router-link class="nav-item" to="/settings/shortcuts"
          >Shortcuts</router-link
        >
        <router-link
          v-if="!localUserExists && localUserArrived"
          class="nav-item"
          to="/newuser"
          >Create new user</router-link
        >
        <router-link class="nav-item" to="/intro">Beginner's guide</router-link>
        <router-link class="nav-item" to="/adminsquickstart"
          >Admin's quickstart</router-link
        >
        <router-link class="nav-item" to="/sfwlist"
          >Safe-for-work list</router-link
        >
        <template v-if="isDev">
          <router-link class="nav-item" to="/namemint">Name minter</router-link>
        </template>
        <div class="spacer"></div>
        <a class="nav-item" href="https://www.getaether.net/faq"
          >FAQ
          <span class="icon-container">
            <icon name="external-link-alt"></icon>
          </span>
        </a>
        <a class="nav-item" href="https://meta.getaether.net"
          >Forums
          <span class="icon-container">
            <icon name="external-link-alt"></icon>
          </span>
        </a>
        <div class="spacer"></div>
        <router-link class="nav-item" to="/about">About</router-link>
        <router-link class="nav-item" to="/changelog">Changelog</router-link>
      </div>
      <router-view></router-view>
    </div>
  </div>
</template>

<script lang="ts">
var ElcUtil = require('../../../../node_modules/electron-util')
var mixins = require('../../mixins/mixins')
export default {
  name: 'settingsscope',
  mixins: [mixins.localUserMixin],
  data() {
    return {}
  },
  computed: {
    isDev(this: any) {
      return ElcUtil.is.development
    },
  },
}
</script>

<style lang="scss">
.settings-sublocation {
  // flex: 1;
  // max-width: 770px; // 80 characters monospaced
  max-width: 545px;
  font-size: 110%;
  padding-right: 5px;
}
</style>

<style lang="scss" scoped>
@import '../../scss/globals';
.location {
  display: flex;
}

.settings-location {
  padding-bottom: 20px;
  padding-right: 14px;
  padding-left: 20px;
}

.settingsscope {
  font-family: 'SSP Regular';
  display: flex; // flex-direction: column;
  margin: auto;
  margin-top: 100px;
  padding: 30px; // width: 800px;
  margin: 20px 20px 0 20px; // background-color: rgba(255, 255, 255, 0.05); // border: 1px solid rgba(0, 0, 0, 0.25);
  margin: 20px auto 0 auto; // background-color: rgba(255, 255, 255, 0.05); // border: 1px solid rgba(0, 0, 0, 0.25);
  background-color: rgba(0, 0, 0, 0.25);
  border-radius: 3px; // padding-right: 20%;
  width: fit-content;
}

.settings-navigation {
  font-family: 'SSP Regular';
  color: $a-grey-500;
  // width: 150px;
  min-width: 150px;
  display: flex;
  flex-direction: column; // padding: 0 10px;
  margin-right: 20px; // color: $a-grey-500;
  padding-top: 77px;
  a {
    color: inherit;
  }
}

a.router-link-exact-active {
  font-family: 'SSP Black';
  color: $a-grey-800;
  letter-spacing: 0px;
}

// a.router-link-active {
//   @extend a.router-link-exact-active;
// }

.spacer {
  height: 24px;
}

.nav-item {
  // font-family: "SCP Regular";
  // font-size: 100%;
  letter-spacing: 0.65px;
  .icon-container {
    svg {
      margin: auto;
      width: 12px;
      height: 12px;
    }
  }
}
</style>
