<template>
  <div class="sidebar-container">
    <div class="sidebar-group specials-and-subs">
      <div class="sidebar-items-list">
        <router-link class="special-sidebar-item" to="/">
          <div class="sidebar-item-icon"></div>
          <div class="sidebar-item-text">Home</div>
        </router-link>
        <router-link class="special-sidebar-item" to="/popular">
          <div class="sidebar-item-icon"></div>
          <div class="sidebar-item-text">Popular</div>
        </router-link>
        <router-link class="special-sidebar-item" to="/new">
          <div class="sidebar-item-icon"></div>
          <div class="sidebar-item-text">New</div>
        </router-link>
        <router-link class="special-sidebar-item" to="/searchscope">
          <div class="sidebar-item-icon"></div>
          <div class="sidebar-item-text">Search</div>
        </router-link>
        <router-link
          tag="div"
          class="sidebar-group-header"
          to="/globalscope/subbed"
        >
          <div class="header-icon"></div>
          <div class="header-text">SUBS</div>
        </router-link>
        <template v-if="ambientBoardsArrived">
          <router-link
            class="sidebar-item iterable"
            v-for="board in ambientBoards"
            :key="board.fingerprint"
            :to="'/board/' + board.fingerprint"
            :class="{
              updated:
                board.lastnewthreadarrived > board.lastseen && board.notify,
            }"
          >
            <div class="sidebar-item-icon">
              <div class="hashimage-carrier">
                <a-hashimage
                  :hash="board.fingerprint"
                  height="20px"
                ></a-hashimage>
              </div>
            </div>
            <div class="sidebar-item-text">{{ board.name }}</div>
            <div class="sidebar-item-notifier">
              <div
                class="notifier-dot"
                v-show="
                  board.lastnewthreadarrived > board.lastseen && board.notify
                "
              ></div>
            </div>
            <div class="silenced-icon" v-show="!board.notify">
              <icon name="regular/bell-slash"></icon>
            </div>
          </router-link>
          <router-link
            class="sidebar-item iterable browse-boards"
            :to="'/globalscope'"
          >
            <div class="sidebar-item-icon">
              <icon name="plus"></icon>
            </div>
            <div class="sidebar-item-text">Browse communities</div>
          </router-link>
        </template>
        <template v-else>
          <div class="spinner-container">
            <div class="spinner-box">
              <a-spinner :hidetext="true" :delay="500"></a-spinner>
            </div>
          </div>
        </template>
      </div>
    </div>
    <router-link tag="div" class="sidebar-group status" to="/status">
      <div class="sidebar-group-header">
        <div class="header-icon"></div>
        <div class="header-text">STATUS</div>
        <!-- <div class="flex-spacer"></div> -->
        <div class="infomark-container">
          <a-info-marker
            header="If both of the lights are green, everything is fine, and you can ignore this pane."
            text="<p>If you're seeing red lights, you should click <a href='#/status'><b>here</b></a> and see what's wrong. If it's not obvious, consider posting a screenshot of your status view to the forum, so we can help you better.</p><p>An untreated warning or error can slow down updates, or completely block you from receiving them.</p>"
          ></a-info-marker>
        </div>
      </div>
      <div class="status-body">
        <div class="status-row">
          <div class="status-row-text">FRONTEND</div>
          <div class="flex-spacer"></div>
          <div class="status-row-value">
            <div class="status-dot" :class="frontendDotColour"></div>
          </div>
        </div>
        <div class="status-row">
          <div class="status-row-text">BACKEND</div>
          <div class="flex-spacer"></div>
          <div class="status-row-value">
            <div class="status-dot" :class="backendDotColour"></div>
          </div>
        </div>
        <!-- <div class="status-row">
          <div class="status-row-text">MINTER</div>
          <div class="flex-spacer"></div>
          <div class="status-row-value" :class="{'active': minterIsWorking}">{{minterState}}</div>
        </div> -->
        <div class="status-row">
          <div class="status-row-text">LAST UPDATE</div>
          <div class="flex-spacer"></div>
          <div class="status-row-value" :class="{ active: refreshIsWorking }">
            {{ refresherState }}
          </div>
        </div>
      </div>
    </router-link>
  </div>
</template>

<script lang="ts">
// var feapiconsumer = require('../services/feapiconsumer/feapiconsumer')
var Vuex = require('../../../node_modules/vuex').default
var globalMethods = require('../services/globals/methods')
export default {
  name: 'a-sidebar',
  data() {
    return {
      allBoards: {},
    }
  },
  computed: {
    ...Vuex.mapState([
      'ambientBoards',
      'ambientStatus',
      'ambientBoardsArrived',
    ]),
    refresherState(this: any): string {
      if (this.refreshIsWorking) {
        return 'Updating'
      }
      let rawrefreshts =
        this.ambientStatus.frontendambientstatus.lastrefreshtimestamp
      if (rawrefreshts === 0) {
        return 'Unknown'
      }
      if (globalMethods.TimeSince(rawrefreshts) === 'just now') {
        return 'Just now'
      }
      return globalMethods.TimeSince(rawrefreshts)
    },
    refreshIsWorking(this: any): boolean {
      return this.ambientStatus.frontendambientstatus.refresherstatus !== 'Idle'
    },
    minterState(this: any): string {
      // make sure that if state is 100% it does not count.
      let state = 'Idle'
      let someItemsPresent = false
      if (
        this.ambientStatus.inflights.boardsList.length > 0 ||
        this.ambientStatus.inflights.threadsList.length > 0 ||
        this.ambientStatus.inflights.postsList.length > 0 ||
        this.ambientStatus.inflights.votesList.length > 0 ||
        this.ambientStatus.inflights.keysList.length > 0 ||
        this.ambientStatus.inflights.truststatesList.length > 0
      ) {
        someItemsPresent = true
      }
      if (!someItemsPresent) {
        return state
      }
      for (let val of this.ambientStatus.inflights.boardsList) {
        if (val.status.completionpercent < 100) {
          state = 'Minting'
          break
        }
      }
      for (let val of this.ambientStatus.inflights.threadsList) {
        if (val.status.completionpercent < 100) {
          state = 'Minting'
          break
        }
      }
      for (let val of this.ambientStatus.inflights.postsList) {
        if (val.status.completionpercent < 100) {
          state = 'Minting'
          break
        }
      }
      for (let val of this.ambientStatus.inflights.votesList) {
        if (val.status.completionpercent < 100) {
          state = 'Minting'
          break
        }
      }
      for (let val of this.ambientStatus.inflights.keysList) {
        if (val.status.completionpercent < 100) {
          state = 'Minting'
          break
        }
      }
      for (let val of this.ambientStatus.inflights.truststatesList) {
        if (val.status.completionpercent < 100) {
          state = 'Minting'
          break
        }
      }
      return state
    },
    minterIsWorking(this: any): boolean {
      if (this.minterState !== 'Idle') {
        return true
      }
      return false
    },
    frontendDotColour(this: any) {
      if (
        this.$store.state.dotStates.frontendDotState === 'status_section_ok'
      ) {
        return 'green'
      }
      if (
        this.$store.state.dotStates.frontendDotState === 'status_section_warn'
      ) {
        return 'yellow'
      }
      if (
        this.$store.state.dotStates.frontendDotState === 'status_section_fail'
      ) {
        return 'red'
      }
      return 'blue'
    },
    backendDotColour(this: any) {
      if (this.$store.state.dotStates.backendDotState === 'status_section_ok') {
        return 'green'
      }
      if (
        this.$store.state.dotStates.backendDotState === 'status_section_warn'
      ) {
        return 'yellow'
      }
      if (
        this.$store.state.dotStates.backendDotState === 'status_section_fail'
      ) {
        return 'red'
      }
      return 'blue'
    },
    renderTimestamp(this: any, ts: any) {
      if (ts !== 'Unknown') {
        return globalMethods.TimeSince(ts)
      }
      return ts
    },
  },
  mounted(this: any) {
    // let vm = this
    // console.log(vm)
    // console.log("side bar is created")
    // feapiconsumer.GetAllBoards(function(result: any) {
    //   console.log("callback gets called")
    //   vm.allBoards = result
    //   console.log(result)
    //   console.log(vm)
    // })
    // setTimeout(function() {
    //   console.log("all boards length")
    //   console.log(vm.allBoards.length)
    // }, 10000)
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';

.sidebar-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  @include generateScrollbar($a-grey-100);
  .sidebar-group {
    // box-shadow: $line-separator-shadow-v2;
    // box-shadow: 0 3px 1px -2px rgba(0, 0, 0, 0.25);
    padding: 5px 0 10px 0;
    &.specials-and-subs {
      flex: 1;
      height: 0; // https://stackoverflow.com/a/14964944 or: I love CSS
      overflow-y: scroll;
    }
    &.global-locations {
    }

    &.status {
      height: 103px; // 3 rows
      // height: 120px; 4 rows
      background-color: #16222a; // $dark-base * 0.8;
      cursor: pointer;
      .status-body {
        cursor: pointer;
      }
      &:hover {
        .sidebar-group-header {
          color: $a-grey-800;
        }
      }
    }

    .sidebar-group-header {
      font-size: 80%;
      padding: 5px 10px;
      cursor: pointer;
      letter-spacing: 1.5px;
      color: $a-grey-500;
      display: flex;
      .flex-spacer {
        flex: 1;
      }
      .infomark-container {
        display: flex;
        margin-left: 4px;
        .info-marker {
          margin: auto; // opacity: 0.8;
          // margin-right: -1px;
          fill: $a-grey-100 !important;
        }
      }
    }

    .special-sidebar-item {
      @extend .sidebar-item;
      &.router-link-exact-active {
        @extend .selected;
      }
    }

    .sidebar-item {
      width: 95%;
      margin: 1px 2% 1px 3%;
      padding: 5px 10px;
      border-radius: 3px;
      cursor: pointer;
      font-family: 'SSP Semibold';
      display: flex;
      color: $a-grey-400;
      word-break: break-word;

      .sidebar-item-icon {
        // width: 18px;
        // height: 18px;
        display: flex;
        svg {
          margin: auto;
          width: 16px;
          height: 14px;
          margin-right: 2px;
        }
        .hashimage-carrier {
          margin: auto;
          margin-right: 10px;
        }
      }
      &:hover {
        background-color: rgba(255, 255, 255, 0.05);
        color: $a-grey-800;
        @extend %link-hover-ghost-extenders-disable;
      }

      &.selected {
        color: $a-grey-800;
        background-color: rgba(255, 255, 255, 0.1); // font-family: "SSP Bold";
        &:hover {
          background-color: rgba(255, 255, 255, 0.15);
        }
      } // &.router-link-exact-active {
      //   @extend .selected;
      // }
      .sidebar-item-text {
        flex: 1;
        line-height: 120%;
      }
      .sidebar-item-notifier {
        display: flex;
        padding-left: 10px;
        width: 18px;
        .notifier-dot {
          margin: auto;
          width: 8px;
          height: 8px;
          border-radius: 4px;
          background-color: $a-cerulean;
        }
      }

      .silenced-icon {
        // padding-right: 3px;
        display: flex;
        margin-right: -3px;
        svg {
          margin: auto;
          width: 15px;
          height: 15px;
        }
      }

      &.updated {
        color: $a-grey-500; // font-family: "SSP Bold";
      }
    }
  }
}

.iterable {
  &.router-link-active {
    @extend .selected;
  }
}

.special-sidebar-item {
  @extend .sidebar-item;

  &.router-link-exact-active {
    @extend .selected;
  }
}

.browse-boards {
  &.router-link-active {
    @extend .selected;
  }
}

.status-body {
  cursor: default;
  .status-row {
    padding: 0 10px;
    display: flex; // font-family: "SSP Semibold";
    color: $a-grey-400;
    font-size: 85%;
    letter-spacing: 1.5px;
    .status-row-value {
      display: flex;
      &.active {
        color: $a-grey-800;
      }
    }
    .status-dot {
      width: 9px;
      height: 9px;
      border-radius: 9px;
      margin: auto;
      &.green {
        background-color: $a-green;
      }
      &.yellow {
        background-color: $a-yellow;
      }
      &.red {
        background-color: $a-red;
      }
      &.blue {
        background-color: $a-cerulean;
      }
    }
  }
}

.flex-spacer {
  flex: 1;
}

.spinner-container {
  display: flex;
  .spinner-box {
    margin: auto;
    padding-top: 50px;
  }
}
</style>

<style lang="scss">
@import '../scss/globals';
.sidebar-group-header {
  .infomark-container {
    .info-marker .info-anchor svg {
      fill: $a-grey-400;
    }
  }
}
</style>
