<template>
  <div class="board-header" v-if="!headerInvisible()">
    <div class="flex-carrier">
      <div class="board-name-container">
        <div class="board-name-image">
          <a-hashimage
            :hash="$store.state.currentBoard.fingerprint"
            height="48px"
          ></a-hashimage>
        </div>
        <div class="board-name-text">
          <template v-if="$store.state.route.name === 'Board>NewThread'">
            New Thread
          </template>
          <template v-else>
            {{ $store.state.currentBoard.name }}
          </template>
        </div>
      </div>
    </div>
    <div class="flex-carrier">
      <div class="board-actions-container" v-show="!newThreadButtonHidden">
        <a
          class="button is-danger is-outlined"
          hasTooltip
          title="Subscribe"
          @click="
            subToBoard({
              fp: $store.state.currentBoard.fingerprint,
              notify: true,
            })
          "
          v-show="!$store.state.currentBoard.subscribed"
        >
          SUB
          <span class="in-button-icon">
            <icon name="plus"></icon>
          </span>
        </a>
        <a
          class="button is-outlined unsub-button"
          @click="unsubFromBoard({ fp: $store.state.currentBoard.fingerprint })"
          hasTooltip
          title="Unsubscribe"
          v-show="$store.state.currentBoard.subscribed"
        >
          SUBBED
        </a>

        <a
          class="button is-warning notifications-button is-outlined"
          hasTooltip
          title="Notifications dots on, click to turn off"
          @click="silenceBoard({ fp: $store.state.currentBoard.fingerprint })"
          v-show="
            $store.state.currentBoard.subscribed &&
              $store.state.currentBoard.notify
          "
        >
          <icon name="bell"></icon>
        </a>
        <a
          class="button is-warning is-outlined notifications-button"
          hasTooltip
          title="Notifications off, click to turn on"
          @click="unsilenceBoard({ fp: $store.state.currentBoard.fingerprint })"
          v-show="
            $store.state.currentBoard.subscribed &&
              !$store.state.currentBoard.notify
          "
        >
          <icon name="regular/bell-slash"></icon>
        </a>

        <template v-if="!localUserReadOnly">
          <router-link
            :to="
              '/board/' + $store.state.currentBoard.fingerprint + '/newthread'
            "
            class="button is-primary is-outlined"
            hasTooltip
            title="Create a new thread"
            >NEW THREAD</router-link
          >
        </template>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var Vuex = require('../../../node_modules/vuex').default
var Tooltips = require('../services/tooltips/tooltips')
var globalMethods = require('../services/globals/methods')
var mixins = require('../mixins/mixins')
export default {
  name: 'a-boardheader',
  mixins: [mixins.localUserMixin],
  data() {
    return {}
  },
  methods: {
    ...Vuex.mapActions([
      'subToBoard',
      'unsubFromBoard',
      'silenceBoard',
      'unsilenceBoard',
    ]),
    headerInvisible: function(this: any) {
      if (this.$route.name === 'Thread') {
        return true
      } else {
        return false
      }
    },
  },
  computed: {
    newThreadButtonHidden(this: any) {
      if (globalMethods.IsUndefined(this.$store.state.route)) {
        return false
      }
      if (this.$store.state.route.name == 'Board>NewThread') {
        return true
      }
      return false
    },
  },
  mounted() {
    Tooltips.Mount()
  },
  updated() {
    // Tooltips.Mount()
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.board-header {
  a {
    @extend %link-hover-ghost-extenders-disable;
  }
  width: 100%;
  height: 175px; // background-color: $dark-base*0.95;
  border-bottom: 1px solid rgba(0, 0, 0, 0.25);
  display: flex;
  background-color: #1a1f23; //$a-grey-200 * 0.5;
}

.flex-carrier {
  flex: 1;
  display: flex;
}

.board-name-container {
  margin: auto;
  margin-left: 50px;
  display: flex;
  .board-name-image {
    width: 48px;
    height: 48px;
    border-radius: 25px; // background-color: $a-cerulean;
  }
  .board-name-text {
    margin: auto 0;
    padding-left: 15px;
    font-size: 150%;
    max-width: 400px;
    height: 36px;
    @include lineClamp(1);
    // Causes text shift between clamped and nonclamped lines. see why. TODO
  }
}

.board-actions-container {
  margin: auto;
  margin-right: 50px;
  display: flex;
  .button {
    margin-left: 5px;
    text-shadow: 0 1px 1px rgba(0, 0, 0, 0.2);
  }
}

.in-button-icon {
  display: inline-block;
  width: 8px; // height: 15px;
  svg {
    height: 11px;
  }
}

.notifications-button {
  // This button needs special treatment because the sizes of the bell icon and the slashed bell icon are not the same.
  svg {
    width: 18px;
    height: 16px;
  }
}

.unsub-button {
  background-color: $a-transparent;
  color: $a-grey-800;
  &:hover {
    background-color: $a-grey-800;
    color: $mid-base;
  }
}
</style>
