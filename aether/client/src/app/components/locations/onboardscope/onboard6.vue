<template>
  <div class="onboard-sublocation">
    <div class="onboard-carrier">
      <transition name="fade" appear>
        <a-markdown class="onboard-markdown last-p-spacer" :content="content"></a-markdown>
      </transition>
      <div class="continue-box">
        <div class="button is-outlined next-button" @click="toNewUser">
          YES
        </div>
        <div class="button is-outlined next-button" @click="toApp">
          NO
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  var fe = require('../../../services/feapiconsumer/feapiconsumer')
  export default {
    name: 'onboard6',
    data() {
      return {
        content: textContent
      }
    },
    methods: {
      toNewUser(this: any) {
        this.$router.push('/newuser')
        fe.SetOnboardComplete(function() {})
      },
      toApp(this: any) {
        this.$router.push('/popular')
        fe.SetOnboardComplete(function() {})
      }
    }
  }
  var textContent = `
## You're ready

Aether has a few more nice things, but you've got the basics covered. Enjoy your stay, and hope I'll see you&nbsp;around!

This is [\`\`\`@nehbit\`\`\`](https://twitter.com/nehbit), by the way, I'm the person behind this&nbsp;thing.

***Heads up:*** *This is a pre-release version. Things will break. Pre-release versions also do collect non-identifying metrics to make sure they work right.*

**Do you want to create a user?**
*You can always create one later.*
  `
</script>

<style lang="scss" scoped>
  @import"../../../scss/bulmastyles";
  @import"../../../scss/globals";
  // .onboard-sublocation {
  //   display: flex;
  //   flex: 1;
  //   min-height: 100%;
  //   .onboard-carrier {
  //     margin: auto;
  //     display: flex;
  //     width: 518px;
  //     flex-direction: column;
  //     position: relative;
  //     bottom: 5%;
  //   }
  // }

  .continue-box {
    display: flex; // margin-top: 5px;
    .next-button {
      @extend %link-hover-ghost-extenders-disable;
      background-color: $a-transparent;
      color: $a-grey-800;
      &:hover {
        background-color: $a-grey-800;
        color: $mid-base;
      }
      animation-duration: 1s;
      animation-name: DELAY_VISIBLE;
      margin-right:10px;
    }
  }

  @keyframes DELAY_VISIBLE {
    0% {
      opacity: 0;
    }
    60% {
      opacity: 0;
    }
    100% {
      opacity: 1;
    }
  }

  .fade-enter-active,
  .fade-leave-active {
    transition-property: opacity;
    transition-duration: .25s;
  }

  .fade-enter-active {
    transition-delay: .25s;
  }

  .fade-enter,
  .fade-leave-active {
    opacity: 0
  }
</style>

<style lang="scss">
  .onboard-markdown {
    font-size: 120%;
    &.last-p-spacer {
      p:last-of-type {
        margin-top: 40px;
      }
    }
    p {
      font-family: "SSP Regular";
    }
    .spacer {
      height: 10px;
      width: 1px;
    }
  }
</style>