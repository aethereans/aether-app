<template>
  <div class="onboard-sublocation">
    <div class="onboard-carrier">
      <transition name="fade" appear>
        <div class="onboard1-container">
          <a-markdown class="onboard-markdown" :content="content"></a-markdown>
          <div class="developer-preview-notice">
            DEVELOPER PREVIEW
            <a-info-marker
              header="This is a preview build. It might break."
              text="<p>If it does, please let us know on the <a href='https://meta.getaether.net'>meta&nbsp;forum</a>.</p>"
            ></a-info-marker>
          </div>
        </div>
      </transition>
      <div class="continue-box">
        <router-link class="button is-outlined next-button" to="/onboard/2">
          CONTINUE
        </router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var ElcUtil = require('../../../../../node_modules/electron-util')
var AutoLaunch = require('../../../../../node_modules/auto-launch')
if (!ElcUtil.is.development) {
  var aetherAutoLauncher = new AutoLaunch({
    name: 'Aether',
    isHidden: true,
  })
  if (!aetherAutoLauncher.isEnabled) {
    aetherAutoLauncher.enable()
  }
}
export default {
  name: 'onboard1',
  data() {
    return {
      content: textContent,
    }
  },
}
var textContent = `
## Welcome in!

Aether is a network of independent, self-moderated communities. Some are based on shared interests, some are friend groups, and some are just so.
  `
</script>

<style lang="scss">
.onboard-sublocation {
  display: flex;
  flex: 1;
  min-height: 100%;
  .onboard-carrier {
    margin: auto;
    display: flex;
    // width: 518px;
    width: 415px;
    flex-direction: column;
    position: relative;
    bottom: 5%;
  }
}
</style>

<style lang="scss" scoped>
@import '../../../scss/globals';
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
  display: flex;
  margin-top: 15px;
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
  transition-duration: 0.25s;
}

.fade-enter-active {
  transition-delay: 0.25s;
}

.fade-enter,
.fade-leave-active {
  opacity: 0;
}

.developer-preview-notice {
  font-family: 'SCP Regular';
  padding: 10px; // border: 1px solid $a-grey-600;
  color: $a-grey-600;
  border-radius: 5px;
  position: fixed;
  bottom: 40px;
  left: 50%;
  letter-spacing: 3px;
  width: 270px;
  text-align: center;
  margin-left: -135px;
  cursor: default;
  animation-duration: 2s; // animation-delay: 1.5s;
  animation-name: DEV_PREVIEW_FLASH;
}

@keyframes DEV_PREVIEW_FLASH {
  0% {
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  20% {
    opacity: 0.3;
  }
  30% {
    opacity: 1;
  }
  40% {
    opacity: 0.3;
  }
  50% {
    opacity: 1;
  }
  60% {
    opacity: 0.3;
  }
  70% {
    opacity: 1;
  }
  80% {
    opacity: 0.3;
  }
  90% {
    opacity: 1;
  }
}
</style>

<style lang="scss">
.onboard-markdown {
  cursor: default;
  font-size: 120%;
  p {
    font-family: 'SSP Regular';
  }
  em {
    font-family: 'SSP Regular Italic';
    strong {
      font-family: 'SSP Bold Italic';
    }
  }
  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
    em {
      font-family: 'SSP Bold Italic';
      strong {
        font-family: 'SSP Black Italic';
      }
    }
  }
}
</style>
