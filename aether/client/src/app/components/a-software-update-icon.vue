<template>
  <div class="software-update-icon" hasTooltip title="New update ready<br> Click here to restart the app" @click="triggerUpdate">
    <a-header-icon class="notifications-icon" icon="regular/arrow-alt-circle-up" iconunder="circle"></a-header-icon>
  </div>
</template>

<script lang="ts">
var ipc = require('../../../node_modules/electron-better-ipc') // Register IPC caller
var Tooltips = require('../services/tooltips/tooltips')
export default {
  name: 'a-software-update-icon',
  data () {
    return {

    }
  },
  mounted() {
    Tooltips.Mount()
  },
  updated() {
    Tooltips.Mount()
  },
  methods: {
    triggerUpdate(this:any) {
      ipc.callMain('RestartToUpdateTheApp')
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .notifications-icon {
    color: $a-grey-800;
  }
</style>

<style lang="scss">
  @import "../scss/globals";
  .software-update-icon {
    .header-icon-container.notifications-icon {
      svg.fa-icon {
        transform: scale(1);
        animation-duration: 3s;
        animation-name: SCALE_BOUNCE;
        animation-iteration-count: infinite;
        @keyframes SCALE_BOUNCE {
          0% {
            transform: scale(1);
          }
          50% {
            transform: scale(1.1);
          }
          100% {
            transform: scale(1);
          }
        }
        svg.icon-background-color {
          fill:$a-red;
          animation-duration: 5s;
          animation-name: UPDATE_NOTIFY;
          animation-iteration-count: infinite;
          // animation-direction: alternate;
          @keyframes UPDATE_NOTIFY {
            0% {
              fill:$a-red * 0.8;
            }
            25% {
              fill:$a-turquoise * 0.8;
            }
            50% {
              fill:$a-cerulean * 0.8;
            }
            75% {
              fill:$a-purple * 0.8;
            }
            100% {
              fill:$a-red * 0.8;
            }
          }
        }
      }
    }
  }
</style>