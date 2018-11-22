<template>
  <div class="a-settings-block">
    <div class="current-state">
      <div class="current-state-header">
        {{name}}
      </div>
      <div class="current-state-text" v-show="isEnabled">
        Enabled
      </div>
      <div class="current-state-text" v-show="!isEnabled">
        Disabled
      </div>
    </div>
    <div class="flex-spacer"></div>
    <div class="button-carrier">
      <a class="button is-success is-outlined" @click="enable" v-show="!isEnabled">
        ENABLE
      </a>
      <a class="button is-success is-outlined" @click="disable" v-show="isEnabled">
        DISABLE
      </a>
    </div>
  </div>
</template>

<script lang="ts">
export default {

  name: 'a-settings-block',

  data () {
    return {
      isEnabled: false,
    }
  },
  props: {
    name: {
      type: String,
      default: "Name of preference goes here"
    },
    stateCheckFunc: {
      type: Function,
      default: function() { return function(){} },
    },
    enableFunc: {
      type: Function,
      default: function() { return function(){} },
    },
    disableFunc: {
      type: Function,
      default: function() { return function(){} },
    },
  },
  methods: {
    async checkState(this:any) {
      this.isEnabled = await this.stateCheckFunc()
    },
    enable(this:any) {
      if (this.isEnabled) {
        return
      }
      this.enableFunc()
      this.checkState()
    },
    disable(this:any) {
      if (!this.isEnabled) {
        return
      }
      this.disableFunc()
      this.checkState()
    },
  },
  beforeMount(this:any) {
    this.checkState()
  },
  updated(this:any) {
    this.checkState()
  }
}
</script>

<style lang="scss" scoped>
  @import"../scss/bulmastyles";
  @import "../scss/globals";
  .a-settings-block {
    background-color: rgba(0, 0, 0, 0.25);
    padding: 15px 20px; // margin: auto;
    font-family: "SSP Bold";
    margin-bottom: 20px;
    border-radius: 3px;
    display: flex;

    .current-state {
      font-family: "SCP Regular";
      // margin-bottom: 10px;
      .current-state-header {
        font-family: "SCP Bold"
      }
    }

    .button-carrier {
      margin:auto;
    }
  }
  .flex-spacer {
    flex: 1;
  }
</style>