<template>
  <div class="header-container" :class="{ 'sidebar-closed': !sidebarOpen }">
    <div class="history-container">
      <a-header-icon
        icon="chevron-left"
        :class="{ disabled: !hasPrevious }"
        @click.native="goBackward"
        class="top-left-rounded-hover"
      ></a-header-icon>
      <a-header-icon
        icon="chevron-right"
        :class="{ disabled: !hasForward }"
        @click.native="goForward"
      ></a-header-icon>
    </div>
    <div class="breadcrumbs-container">
      <a-breadcrumbs></a-breadcrumbs>
    </div>
    <a-software-update-icon v-if="newUpdateReady"></a-software-update-icon>
    <a-notifications-icon v-if="localUserExists"></a-notifications-icon>
    <!-- <a-search-icon></a-search-icon> -->
    <div
      class="profile-container"
      @click="toggleUserMenu"
      v-click-outside="onClickOutside"
    >
      <div class="dropdown-container">
        <div class="dropdown is-right" :class="{ 'is-active': userMenuOpen }">
          <div class="dropdown-trigger">
            <template v-if="!$store.state.localUserArrived">
              <div
                class="user-name"
                :class="{ 'read-only': this.localUserReadOnly }"
              >
                Refreshing...
              </div>
            </template>
            <template v-if="$store.state.localUserArrived">
              <div class="info-marker-container" v-if="localUserReadOnly">
                <a-info-marker
                  header="You are in read only mode."
                  text="<p>You haven't created an user. You can read, but you won't be able to post or vote until you create one.</p>"
                ></a-info-marker>
              </div>
              <div
                class="user-name"
                :class="{ 'read-only': this.localUserReadOnly }"
              >
                {{ localUserName }}
              </div>
              <div class="mod-puck-container" v-show="isMod">
                <div class="mod-puck">mod</div>
              </div>
            </template>
            <div class="profile-caret-icon">
              <icon name="chevron-down"></icon>
            </div>
          </div>
          <div class="dropdown-menu" id="dropdown-menu" role="menu">
            <div class="dropdown-content">
              <template v-if="!localUserExists">
                <!-- Local user doesn't exist, not created yet -->
                <!--    <router-link to="/settings/newuser" class="dropdown-item">
                  Create user
                </router-link> -->
                <div class="button-container">
                  <router-link
                    class="button is-link is-outlined is-small join-button"
                    to="/newuser"
                  >
                    JOIN AETHER
                  </router-link>
                </div>
                <hr class="dropdown-divider" />
              </template>
              <template v-if="!localUserReadOnly">
                <a-avatar-block
                  nofingerprint="true"
                  :user="$store.state.localUser"
                  :clickable="true"
                  :imageheight="96"
                ></a-avatar-block>
                <hr class="dropdown-divider" />
                <router-link
                  :to="'/user/' + $store.state.localUser.fingerprint"
                  class="dropdown-item"
                >
                  Profile
                </router-link>
              </template>
              <router-link to="/settings" class="dropdown-item">
                Preferences
              </router-link>
              <router-link to="/membership" class="dropdown-item">
                Supporter benefits
              </router-link>
              <hr class="dropdown-divider" />
              <router-link to="/intro" class="dropdown-item">
                Beginner's guide
              </router-link>
              <router-link to="/about" class="dropdown-item">
                About
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="windows-toggles-container" v-if="isWindows">
      <template v-if="isWindows10">
        <div class="windows-window-toggles">
          <div class="win-btn minimise" @click="minimiseWindow">
            <span>&#xE921;</span>
          </div>
          <div
            class="win-btn maximise"
            v-show="!maximised"
            @click="maximiseWindow"
          >
            <span>&#xE922;</span>
          </div>
          <div
            class="win-btn restore"
            v-show="maximised"
            @click="restoreWindow"
          >
            <span>&#xE923;</span>
          </div>
          <div class="win-btn close" @click="closeWindow">
            <span>&#xE8BB;</span>
          </div>
        </div>
      </template>
      <template v-else>
        <!-- Older versions of Windows that don't have Segoe MDL2  -->
        <!-- We have to provide our own icons here. -->
        <div class="windows-window-toggles">
          <div class="win-btn minimise" @click="minimiseWindow">
            <icon name="regular/window-minimize"></icon>
          </div>
          <div
            class="win-btn maximise"
            v-show="!maximised"
            @click="maximiseWindow"
          >
            <icon name="regular/square"></icon>
          </div>
          <div
            class="win-btn restore"
            v-show="maximised"
            @click="restoreWindow"
          >
            <icon name="window-maximize"></icon>
          </div>
          <div class="win-btn close" @click="closeWindow">
            <!-- <icon name="regular/times"></icon> -->
            <img class="close-btn" src="../ext_dep/svg/close-btn.svg" alt="" />
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
<script lang="ts">
const remote = require('electron').remote
let os = require('os')
const isDev = require('electron-is-dev')
if (isDev) {
  var Vue = require('../../../node_modules/vue/dist/vue.js') // Production
} else {
  var Vue = require('../../../node_modules/vue/dist/vue.min.js') // Production
}
var Vuex = require('../../../node_modules/vuex').default
var mixins = require('../mixins/mixins')
var globalMethods = require('../services/globals/methods')
export default Vue.extend({
  name: 'a-header',
  mixins: [mixins.localUserMixin],
  data() {
    return {
      userMenuOpen: false,
      maximised: false,
    }
  },
  computed: {
    ...Vuex.mapState(['sidebarOpen', 'newUpdateReady']),
    localUserName(this: any) {
      if (this.localUserReadOnly) {
        return 'Menu'
      }
      return globalMethods.GetUserName(this.$store.state.localUser)
    },
    isMod(this: any) {
      if (
        this.$store.state.modModeEnabledArrived &&
        this.$store.state.modModeEnabled
      ) {
        return true
      }
      return false
    },
    hasPrevious(this: any) {
      return this.$store.state.historyMaxCaret > 1
    },
    hasForward(this: any) {
      return (
        this.$store.state.historyCurrentCaret <
        this.$store.state.historyMaxCaret
      )
    },
    isWindows(this: any) {
      return process.platform === 'win32'
    },
    osMajorVersion(this: any) {
      return parseInt(os.release().split('.')[0])
    },
    isWindows10(this: any) {
      if (!this.isWindows) {
        return false
      }
      if (this.osMajorVersion < 10) {
        // W10 and future should have our native typeface that provides the native interface.
        return false
      }
      return true
    },
  },
  methods: {
    ...Vuex.mapActions(['setSidebarState']),
    toggleUserMenu(): void {
      ;(<any>this)['userMenuOpen']
        ? ((<any>this)['userMenuOpen'] = false)
        : ((<any>this)['userMenuOpen'] = true)
    },
    onClickOutside() {
      ;(<any>this)['userMenuOpen'] = false
    },
    toggleSidebarState() {
      this.$store.state.sidebarOpen === true
        ? this.setSidebarState(false)
        : this.setSidebarState(true)
    },
    signIn(this: any) {},
    signOut(this: any) {},
    goForward(this: any) {
      // if (!this.hasForward) {
      //   return
      // }
      this.$store.dispatch('registerNextActionIsHistoryMoveForward')
      this.$router.go(+1)
      // this.$router.push({ path: this.$routerHistory.next().path })
    },
    goBackward(this: any) {
      // if (!this.hasPrevious) {
      //   return
      // }
      this.$store.dispatch('registerNextActionIsHistoryMoveBack')
      this.$router.go(-1)
      // this.$router.push({ path: this.$routerHistory.previous().path })
    },
    minimiseWindow(this: any) {
      remote.getCurrentWindow().minimize()
    },
    maximiseWindow(this: any) {
      remote.getCurrentWindow().maximize()
      this.maximised = true
    },
    restoreWindow(this: any) {
      remote.getCurrentWindow().unmaximize()
      this.maximised = false
    },
    closeWindow(this: any) {
      remote.getCurrentWindow().close()
    },
  },
})
</script>
<style lang="scss" scoped>
@import '../scss/globals';
.header-container {
  -webkit-app-region: drag;
  /* ^ You must mark all interactive objects within this as NON-draggable to make this work, otherwise they will just be effectively unclickable. */
  height: $top-bar-height;
  background-color: #25272a; // $mid-base * 0.9;
  /*border-bottom: 1px solid #111;*/
  box-shadow: $line-separator-shadow, $line-separator-shadow-castleft-light;
  position: relative; // z-index: 3;
  display: flex;
  flex: 1;
  border-radius: 10px 0 0 0;
  min-width: 0;

  .top-left-rounded-hover {
    &:hover {
      border-radius: 10px 0 0 0;
    }
  }

  &.sidebar-closed {
    border-radius: 0;
  }
  .history-container {
    display: flex;
    .disabled {
      opacity: 0.15;
      cursor: default;
      &:hover {
        background-color: $a-transparent;
      }
    }
  }

  .breadcrumbs-container {
    flex: 1;
    min-width: 0;
  }

  .profile-container {
    -webkit-app-region: no-drag;
    display: flex; // padding: 0 15px;
    padding-right: 15px;
    padding-left: 7px;

    &:hover {
      background-color: rgba(255, 255, 255, 0.05);
    }

    .dropdown-container {
      display: flex;

      .profile-header {
        // width: 225px;
        width: 200px;
      }

      .dropdown-trigger {
        display: flex;
        cursor: pointer;

        .user-name {
          &.read-only {
            // font-family: "SSP Semibold Italic"; // color: $a-grey-400;
            font-family: 'SSP Semibold';
          }
        }
        .info-marker-container {
          padding: 0 3px;
          padding-right: 7px;
          margin-top: 1px;
          fill: $a-grey-400;
        }

        .profile-caret-icon {
          display: flex;
          padding-left: 8px;
          padding-top: 2px;
          svg {
            margin: auto;
          }
        }
      }
      .dropdown {
        margin: auto;
        .dropdown-divider {
          background-color: $dropdown-divider-color;
        }
      }
    }
  }
}

.button-container {
  width: 100%;
  display: flex;
  padding: 5px 15px;
  .join-button {
    flex: 1;
    font-family: 'SCP Bold';
    font-size: 14px;
  }
}

.mod-puck-container {
  display: flex;
  .mod-puck {
    font-family: 'SCP Bold';
    letter-spacing: 2px;
    font-size: 90%;
    border-radius: 5px;
    padding: 0 3px 0 6px;
    margin: auto 0px auto 8px;
    border: 1px solid $a-purple;
    color: $a-purple;
    line-height: 135%;
  }
}

.windows-toggles-container {
  -webkit-app-region: no-drag;
  display: flex;
  .windows-window-toggles {
    display: flex;
    margin: auto;
    font-family: 'Segoe MDL2 Assets'; // Latter for Win 7/8
    height: 100%;
    .win-btn {
      width: 46px;
      font-size: 10px; // height: 32px;
      display: flex;
      cursor: default; // Windows doesn't do cursor:pointer on those.
      span {
        margin: auto;
      }
      &:hover {
        background-color: rgba(255, 255, 255, 0.05);
        &.close {
          background-color: $a-red-80;
        }
      }
      svg {
        // This is for W10<. W10 uses font-family so it's not svg.
        margin: auto;
        width: 14px;
        height: 14px;
      }
      img.close-btn {
        width: 18px;
        height: 18px;
        margin: auto;
      }
    }
  }
}
</style>
