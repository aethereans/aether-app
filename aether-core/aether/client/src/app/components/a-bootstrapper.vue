<template>
  <div
    class="bootstrapper"
    v-if="
      containerVisible &&
      (bootstrapInProgress ||
        bootstrapRefreshInProgress ||
        bootstrapRefreshComplete)
    "
  >
    <div class="bootstrapper-carrier">
      <div class="bootstrap-message" v-if="bootstrapInProgress">
        <div class="spinner-container">
          <a-spinner :hidetext="true"></a-spinner>
        </div>
        <div class="text-carrier">
          <b>Catching up:</b> Downloading fresh content from Aether... (1-5
          min.)
        </div>
      </div>
      <div class="bootstrap-message" v-if="bootstrapRefreshInProgress">
        <div class="spinner-container">
          <a-spinner :hidetext="true"></a-spinner>
        </div>
        <div class="text-carrier">
          <b>Catching up:</b> Compiling the entity graph...
        </div>
      </div>
      <div class="bootstrap-message" v-if="bootstrapRefreshComplete">
        <div class="text-carrier">
          <b>Catching up:</b> Complete. You now have the most recent content.
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var Vuex = require('../../../node_modules/vuex').default
export default {
  name: 'a-bootstrapper',
  data() {
    return {
      containerVisible: false,
    }
  },
  computed: {
    ...Vuex.mapState(['ambientStatus']),
    bootstrapInProgress(this: any) {
      return this.ambientStatus.backendambientstatus.bootstrapinprogress
    },
    lastBootstrapTimestamp(this: any) {
      return this.ambientStatus.backendambientstatus.lastbootstraptimestamp
    },
    bootstrapRefreshInProgress(this: any) {
      return this.ambientStatus.frontendambientstatus.bootstraprefreshinprogress
    },
    bootstrapRefreshComplete(this: any) {
      return this.ambientStatus.frontendambientstatus.bootstraprefreshcomplete
    },
  },
  methods: {
    isVisible(this: any): boolean {
      // Runs once at the beginning
      if (this.bootstrapInProgress) {
        return true
      }
      return false
    },
  },
  beforeMount(this: any) {
    if (this.isVisible()) {
      this.containerVisible = true
    }
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.bootstrapper {
  font-family: 'SCP Regular';
  display: flex;
  border-radius: 3px;
  margin: 0 20px;
  margin-top: 20px;
  background-color: rgba(0, 0, 0, 0.25);
  .bootstrapper-carrier {
    padding: 15px 20px; // width: 100%;
  }
  .bootstrap-message {
    display: flex;
    .text-carrier {
      margin: auto;
      margin-left: 15px;
      b {
        font-family: 'SCP Bold';
      }
    }
  }
}

.flex-spacer {
  flex: 1;
}

.spinner-container {
  display: flex;
  .spinner {
    margin: auto;
  }
}
</style>
