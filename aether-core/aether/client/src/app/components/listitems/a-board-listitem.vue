<template>
  <router-link
    class="board-entity"
    :to="linkToBoard"
    :class="{ inflight: isInflightEntity, uncompiled: uncompiled }"
  >
    <div class="flex-carrier signals" v-if="signalsVisible">
      <div class="signals-container">
        <div class="nsfw-box" v-show="nsfwBoard">
          <span class="nsfw-text">NSFW</span>
        </div>
        <div class="population-count">
          <icon name="regular/user-circle"></icon> {{ userCount }}
        </div>
        <div class="threads-count">
          <icon name="regular/comments"></icon> {{ threadsCount }}
        </div>
      </div>
    </div>
    <div class="flex-carrier image" v-if="imageVisible">
      <div class="image-container">
        <a-hashimage :hash="boardFingerprint" height="72px"></a-hashimage>
      </div>
    </div>
    <div class="flex-carrier main">
      <div class="main-data-container">
        <div class="inflight-box" v-if="isInflightEntity">
          <a-inflight-info
            :status="inflightStatus"
            :refresherFunc="refresh"
          ></a-inflight-info>
        </div>
        <div class="board-name">
          {{ board.name }}
        </div>
        <div class="board-description">
          {{ board.description }}
        </div>
      </div>
    </div>
    <div class="flex-carrier actions" v-if="actionsVisible">
      <div class="actions-container">
        <div class="flex-spacer"></div>
        <a
          class="button is-danger"
          hasTooltip
          title="Subscribe"
          @click.prevent="subToBoard({ fp: boardFingerprint, notify: true })"
          v-show="!board.subscribed"
        >
          SUB
          <span class="in-button-icon">
            <icon name="plus"></icon>
          </span>
        </a>

        <a
          class="button is-outlined unsub-button"
          @click.prevent="unsubFromBoard({ fp: boardFingerprint })"
          hasTooltip
          title="Unsubscribe"
          v-show="board.subscribed"
        >
          SUBBED
        </a>
      </div>
    </div>
  </router-link>
</template>

<script lang="ts">
var Vuex = require('../../../../node_modules/vuex').default
var vuexStore = require('../../store/index').default
export default {
  name: 'a-board-listitem',
  props: ['board', 'inflightStatus', 'uncompiled', 'refresher'],
  data(this: any) {
    return {
      nsfwBoard: !this.board.sfwlisted,
    }
  },
  methods: {
    ...Vuex.mapActions([
      'subToBoard',
      'unsubFromBoard',
      'silenceBoard',
      'unsilenceBoard',
    ]),
    refresh(this: any) {
      this.refresher(this.$store.state.currentUserEntity.fingerprint)
      vuexStore.dispatch('pruneInflights')
    },
  },
  computed: {
    boardCreation(this: any) {
      // These are necessary because in uncompiled entities, these are in thread.provable.creation, but in compiled ones it's in thread.creation.
      if (this.uncompiled || this.isInflightEntity) {
        if (typeof this.board.provable === 'undefined') {
          return 0
        }
        return this.board.provable.creation
      }
      return this.board.creation
    },
    boardLastUpdate(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        if (typeof this.board.updateable === 'undefined') {
          return 0
        }
        return this.board.updateable.lastupdate
      }
      return this.board.lastupdate
    },
    boardFingerprint(this: any) {
      if (this.uncompiled) {
        return this.board.provable.fingerprint
      }
      return this.board.fingerprint
    },
    signalsVisible(this: any) {
      if (this.uncompiled) {
        return false
      }
      return true
    },
    actionsVisible(this: any) {
      if (this.uncompiled) {
        return false
      }
      if (this.isInflightEntity) {
        return false
      }
      return true
    },
    isInflightEntity(this: any) {
      if (typeof this.inflightStatus !== 'undefined') {
        return true
      }
      return false
    },
    userCount(this: any) {
      if (typeof this.board.usercount !== 'undefined') {
        return this.board.usercount
      }
      return 0
    },
    threadsCount(this: any) {
      if (typeof this.board.threadscount !== 'undefined') {
        return this.board.threadscount
      }
      return 0
    },
    linkToBoard(this: any) {
      // if (this.uncompiled) {
      //   return ''
      // }
      if (this.isInflightEntity) {
        return ''
      }
      return '/board/' + this.boardFingerprint
    },
    imageVisible(this: any) {
      // if (this.uncompiled) {
      //   return false
      // }
      if (this.isInflightEntity) {
        return false
      }
      return true
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../../scss/globals';
a:hover {
  color: inherit;
}

.board-entity {
  display: block;
  @extend %link-hover-ghost-extenders-disable;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  padding: 15px 5px;
  margin: 0 20px; // cursor: default;
  color: $a-grey-800;
  word-break: break-word;
  &:hover {
    color: $a-grey-800;
    background-color: rgba(255, 255, 255, 0.05);
  }
  a {
    @extend %link-hover-ghost-extenders-disable;
    &:hover {
      color: inherit;
    }
  }
  &.inflight {
    cursor: default;
    .main-data-container .board-name {
      cursor: default;
    }
    border-bottom: none;
  }
}

.flex-carrier {
  display: flex;
  &.signals {
  }
  &.main {
    flex: 1;
  }
}

.signals-container {
  min-width: 64px;
  padding: 0px 10px;
  white-space: nowrap;
  font-size: 110%;
  margin: auto;
  display: flex;
  flex-direction: column;
  color: $a-grey-600;
  .nsfw-box {
    border: 1px solid $a-red;
    text-align: center;
    padding: 1px;
    white-space: nowrap;
    margin: auto;
    display: flex;
    .nsfw-box-text {
      color: $a-red;
      vertical-align: middle;
    }
  }
  .threads-count,
  .population-count {
    margin: auto;
    display: flex;
    svg {
      margin: auto;
      margin-right: 3px;
    }
  }
}

.main-data-container {
  flex: 1;
  padding-right: 15px;
  .board-name {
    font-size: 120%;
    cursor: pointer;
  }
  .board-description {
    font-family: 'SSP Regular';
    color: $a-grey-500;
    font-size: 95%;
    @include lineClamp(3);
  }
}

.image-container {
  margin: auto;
  padding-right: 15px;
}

.actions-container {
  margin: auto;
  width: 100px;
  display: flex;
  margin-right: 25px;
  .flex-spacer {
    flex: 1;
  }
  .in-button-icon {
    display: inline-block;
    width: 8px; // height: 15px;
    svg {
      height: 11px;
    }
  }
}

a.unsub-button {
  background-color: $a-transparent;
  color: $a-grey-800;
  &:hover {
    background-color: $a-grey-800;
    color: $mid-base;
  }
}
</style>
