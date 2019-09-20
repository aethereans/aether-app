<template>
  <div class="a-settings-block">
    <div class="current-state">
      <div class="current-state-header">
        {{ name }}
      </div>
      <div class="current-state-text" v-show="isEnabled">
        Enabled
      </div>
      <div class="current-state-text" v-show="!isEnabled">
        Disabled
      </div>
    </div>
    <div class="flex-spacer"></div>
    <div class="button-carrier" v-if="!readOnly">
      <template v-if="isLoadingState">
        <a class="button is-success is-outlined is-loading">
          DISABLE
          <!-- This text above determines button width. Set it to larger text. -->
        </a>
      </template>
      <template v-else>
        <a
          class="button is-success is-outlined"
          @click="enable"
          v-show="!isEnabled"
        >
          ENABLE
        </a>
        <a
          class="button is-success is-outlined"
          @click="disable"
          v-show="isEnabled"
        >
          DISABLE
        </a>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-settings-block',

  data() {
    return {
      isEnabledState: false,
    }
  },
  props: {
    name: {
      type: String,
      default: 'Name of preference goes here',
    },
    stateCheckFunc: {
      type: Function,
      default: function() {
        return function() {}
      },
    },
    manualState: {
      type: Boolean,
      default: false,
    },
    useManualState: {
      type: Boolean,
      default: false,
    },
    isLoadingState: {
      type: Boolean,
      default: false,
    },
    enableFunc: {
      type: Function,
      default: function() {
        return function() {}
      },
    },
    disableFunc: {
      type: Function,
      default: function() {
        return function() {}
      },
    },
    readOnly: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isEnabled(this: any) {
      if (this.useManualState) {
        return this.manualState
      }
      return this.isEnabledState
    },
  },
  methods: {
    async checkState(this: any) {
      this.isEnabledState = await this.stateCheckFunc()
    },
    async enable(this: any) {
      if (this.isEnabled) {
        return
      }
      await this.enableFunc()
      if (!this.useManualState) {
        await this.sleep(100)
        await this.checkState()
      }
    },
    async disable(this: any) {
      if (!this.isEnabled) {
        return
      }
      await this.disableFunc()
      if (!this.useManualState) {
        await this.sleep(100)
        await this.checkState()
      }
    },
    async sleep(ms: any) {
      /*
        This is a hack - however it's necessary because we don't actually know if the function that is going to be given to this block is async or not. If it's async, it'll do its own wait. If it's synchronous, and will make an external roundtrip, we don't have any choice but to wait a set amount of time before we poll for an update to the state.
      */
      return new Promise(function(resolve: any) {
        setTimeout(resolve, ms)
      })
    },
  },
  async beforeMount(this: any) {
    await this.checkState()
  },
  async updated(this: any) {
    await this.checkState()
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.a-settings-block {
  background-color: rgba(0, 0, 0, 0.25);
  padding: 15px 20px; // margin: auto;
  font-family: 'SSP Bold';
  margin-bottom: 20px;
  border-radius: 3px;
  display: flex;

  .current-state {
    font-family: 'SCP Regular';
    // margin-bottom: 10px;
    .current-state-header {
      font-family: 'SCP Bold';
    }
  }

  .button-carrier {
    margin: auto;
  }
}
.flex-spacer {
  flex: 1;
}
</style>
