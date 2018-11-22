<template>
  <div class="view-header">
    <!-- <div class="side-container left"></div> -->
    <div class="center-container">
      <div class="header-container side-layout">
        <template v-if="$store.state.route.name === 'Global>NewBoard'">
          <div class="header-text">New community</div>
          <div class="header-subtext">Hey there future admin! Do take a look at
            <router-link class="quickstart-link" to="/adminsquickstart">Admin's quickstart</router-link> before diving in. <br>
            If your new community is safe for work, consider getting it added to the <router-link class="quickstart-link" to="/sfwlist">SFW list</router-link>.
          </div>
        </template>

        <template v-if="$store.state.route.name === 'Global>Subbed'">
          <div class="header-text">Subscribed</div>
          <div class="header-subtext">Communities you're currently subscribed to
            <div class="infomark-container">
              <a-info-marker header="If you subscribe to a community, some of its most popular posts will appear in your home feed. " text="<p></p><p>Your subscriptions are private to you and saved only locally. No one sees what you sub to. (The number of users count visible on the community list is the number of posters, not subscribers.) </p>"></a-info-marker>
            </div>
          </div>
        </template>

        <template v-if="$store.state.route.name === 'Global'">
          <div class="header-text">Browse communities</div>
          <div class="header-subtext">Discover new communities on the network
            <div class="infomark-container">
              <a-info-marker header="If you subscribe to a community, some of its most popular posts will appear in your home feed. " text="<p></p><p>Your subscriptions are private to you and saved only locally. No one sees what you sub to. The number of users count visible on the community list is the number of posters, not subscribers. </p>"></a-info-marker>
            </div>
          </div>
        </template>

      </div>
    </div>
    <div class="side-container right" v-show="!newBoardButtonHidden">
      <div class="actions-container">
        <router-link to="/globalscope/newboard" class="button is-warning is-outlined" title="Create a new community" hasTooltip v-if="!localUserReadOnly" >
          NEW COMMUNITY
        </router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  var Tooltips = require('../services/tooltips/tooltips')
  var mixins = require('../mixins/mixins')
  var globalMethods = require('../services/globals/methods')
  export default {
    name: 'a-global-header',
    mixins: [mixins.localUserMixin],
    data() {
      return {}
    },
    computed: {
      newBoardButtonHidden(this: any) {
        if (globalMethods.IsUndefined(this.$store.state.route)) {
          return false
        }
        if (this.$store.state.route.name == 'Global>NewBoard') {
          return true
        }
        return false
      }
    },
    mounted() {
      Tooltips.Mount()
    },
    updated() {
      Tooltips.Mount()
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/bulmastyles";
  @import "../scss/globals";
  .view-header {
    height: 175px;
    display: flex;
    border-bottom: 1px solid rgba(0, 0, 0, 0.15);
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.25), 1px 0 0 0 rgba(0, 0, 0, 0.2) inset;
    background-color: $a-grey-200*0.5;
    border-left: 5px solid $a-yellow;
    line-height: 200%;

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
        margin-bottom: 32px;
        &.side-layout {
          .header-text,
          .header-subtext {
            text-align: left;
            padding: 0;
            line-height: 140%;
            padding-top: 15px;
          }
          .header-subtext {
            font-family: "SSP Regular"
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
    margin-bottom: 45px;
    a {
      @extend %link-hover-ghost-extenders-disable;
    }
  }

  .quickstart-link {
    font-family: "SSP Bold"
  }

  .quickstart-link:hover {
    color: $a-cerulean;
  }

  .infomark-container {
    margin-left: 3px;
    display: inline-block;
  }
</style>