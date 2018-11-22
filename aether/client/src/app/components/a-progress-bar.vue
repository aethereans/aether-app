<template>
  <div class="progress-bar">
    <progress class="progress-element" :max="max" :value="currentPercent"></progress>
  </div>
</template>

<script lang="ts">
  var globalMethods = require('../services/globals/methods')
  export default {
    name: 'a-progress-bar',
    props: {
      percent: {
        type: Number,
        default: 0
      },
      max: {
        type: Number,
        default: 100
      },
      animationsDisabled: {
        type: Boolean,
        default: false
      },
    },
    data() {
      return {
        currentPercent: 1,
        animateDeltaRunning: false,
        startingPercentage: 0,
        mountTime: 0,
        firstStepAnimationSkipped: false,
      }
    },
    beforeMount(this: any) {
      this.startingPercentage = this.percent
    },
    mounted(this: any) {
      this.mountTime = globalMethods.NowUnix()
      this.animateDelta(this)
    },
    updated(this: any) {},
    watch: {
      percent(this: any) {
        this.animateDelta(this)
      }
    },
    methods: {
      animateDelta(this: any) {
        if (this.animationsDisabled) {
          this.currentPercent = this.percent
          return
        }
        /*
          The logic here is: If the progress bar appears within 5 seconds of page load, then we have a prior progress still going on, so we do not want to animate between zero and current progress. (But we want to still animate the remainder. If a progress bar starts from 33, 0 to 33 should not be animated, but 33 to 100 should be.)

          But if the bar happens after 5 seconds, it's usually because the user has requested some action, and that action starts from 0, so we do want to animate it.

          Effectively, if the progress bar loads within the first seconds of page load, we skip one 'step' in the animation, and start animating from the next checkpoint.
        */
        if (!this.firstStepAnimationSkipped) {
          if ((this.mountTime - this.$store.state.lastPageLoadTimestamp < 3) || (this.$store.state.lastPageLoadTimestamp === 0)) {
            this.currentPercent = this.percent
            this.firstStepAnimationSkipped = true
          }
        }
        this.animateDeltaRunning = true
        let stepCount = 0
        let vm = this
        let step = function() {
          if (vm.percent === -1) {
            vm.currentPercent = 100
            return
          }
          if (vm.percent === 100) {
            if (vm.startingPercentage === 100) {
              // It started at 100, which means when the page was loaded, this was already complete. We don't animate that - Just skip it.
              vm.currentPercent = 100
            }
            if (vm.currentPercent < 100) {
              vm.currentPercent = vm.currentPercent + 2
              // ^ Complete the rest in 1/3 of a sec
              requestAnimationFrame(step)
            }
            return
          }
          if (stepCount > 300) {
            // after 5 seconds, stop advancement, it'll be basically unmoving
            vm.animateDeltaRunning = false
            return
          }
          stepCount++
          let incr = (vm.percent - vm.currentPercent) / 60 // Assuming 60fps
          vm.currentPercent = vm.currentPercent + incr
          requestAnimationFrame(step)
        }
        requestAnimationFrame(step)
      }
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .progress-bar {
    width: 100%;
    display: flex;
    .progress-element {
      flex: 1;
      -webkit-appearance: none;
      appearance: none;
      height: 3px;
      border-radius: 1.5px;
      overflow: hidden;
      &::-webkit-progress-bar {
        border-radius: 1.5px;
        background-color: $a-grey-400;
      }
      &::-webkit-progress-value {
        background-color: $a-orange;
      }
    }
  }
</style>