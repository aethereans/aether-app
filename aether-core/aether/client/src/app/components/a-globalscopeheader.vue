<template>
  <div class="global-scope-header">
    <div
      class="side-container left"
      v-show="$store.state.route.name !== 'Global>NewBoard'"
    ></div>
    <div class="center-container">
      <div
        class="header-container"
        :class="{
          'side-layout': $store.state.route.name === 'Global>NewBoard',
        }"
      >
        <template v-if="$store.state.route.name === 'Global>NewBoard'">
          <div class="header-text">New community</div>
          <div class="header-subtext">
            Hey there future admin, welcome in! <br />Consider taking a look at
            <router-link class="quickstart-link" to="/adminsquickstart"
              >Admin's Quickstart</router-link
            >
            before diving in.
          </div>
        </template>
        <template v-else>
          <div class="header-text">All communities</div>
          <div class="header-subtext">
            Here you can discover new communities and manage your subscriptions.
          </div>
        </template>
      </div>
    </div>
    <div class="side-container right">
      <div class="actions-container">
        <router-link
          to="/globalscope/newboard"
          class="button  is-info is-outlined"
          title="Create a new community"
          v-if="!localUserReadOnly"
          :disabled="newBoardButtonDisabled"
        >
          NEW COMMUNITY
          <!-- <span class="in-button-icon">
            <icon name="plus" ></icon>
          </span> -->
        </router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var mixins = require('../mixins/mixins')
var globalMethods = require('../services/globals/methods')
export default {
  name: 'a-globalscopeheader',
  mixins: [mixins.localUserMixin],
  data() {
    return {}
  },
  computed: {
    newBoardButtonDisabled(this: any) {
      if (globalMethods.IsUndefined(this.$store.state.route)) {
        return false
      }
      if (this.$store.state.route.name == 'Global>NewBoard') {
        return true
      }
      return false
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.global-scope-header {
  height: 175px;
  display: flex;
  border-bottom: 1px solid rgba(0, 0, 0, 0.25);
  .side-container {
    flex: 1;
    display: flex;
  }
  .center-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    .header-container {
      margin: auto;
      text-align: center;
      &.side-layout {
        .header-text,
        .header-subtext {
          text-align: left;
          padding: 0;
        }
        .header-subtext {
          font-family: 'SSP Regular';
        }
        margin-left: 50px;
      }
      .header-text {
        font-size: 200%;
      }
      .header-subtext {
        padding: 0px 50px;
      }
    }
  }
}

.in-button-icon {
  display: inline-block;
  width: 8px; // height: 15px;
  svg {
    height: 11px;
  }
}

.actions-container {
  margin: auto;
  margin-right: 50px;
  a {
    @extend %link-hover-ghost-extenders-disable;
  }
}

.quickstart-link:hover {
  color: $a-cerulean;
}
</style>
