<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="startAtBootContent"></a-markdown>
    <!--  <a
      class="start-at-boot-toggle"
      @click="startAtBootToggleVisible = !startAtBootToggleVisible"
      >Check start at boot state</a
    > -->
    <div class="start-at-boot-toggle-container">
      <a
        class="button is-small is-outlined is-white start-at-boot-toggle"
        @click="startAtBootToggleVisible = !startAtBootToggleVisible"
        >CHECK START AT BOOT</a
      >
      <p><i>(Your OS might ask you for permissions)</i></p>
    </div>

    <a-settings-block
      v-if="startAtBootToggleVisible"
      name="Start at boot"
      :stateCheckFunc="startAtBootCheckFunc"
      :enableFunc="startAtBootEnableFunc"
      :disableFunc="startAtBootDisableFunc"
    ></a-settings-block>
    <br />
    <a-markdown :content="autoloadContent"></a-markdown>
    <a-settings-block
      name="Autoload whitelisted media"
      :stateCheckFunc="autoloadCheckFunc"
      :enableFunc="autoloadEnableFunc"
      :disableFunc="autoloadDisableFunc"
    ></a-settings-block>
  </div>
</template>

<script lang="ts">
var fe = require('../../../services/feapiconsumer/feapiconsumer')
var AutoLaunch = require('../../../../../node_modules/auto-launch')
var aetherAutoLauncher = new AutoLaunch({
  name: 'Aether',
  isHidden: true,
})
export default {
  name: 'settingsroot',
  data(this: any) {
    return {
      headline: headline,
      intro: intro,
      /*----------  Start at boot  ----------*/
      startAtBootContent: startAtBootContent,
      startAtBootCheckFunc: aetherAutoLauncher.isEnabled,
      startAtBootEnableFunc: aetherAutoLauncher.enable,
      startAtBootDisableFunc: aetherAutoLauncher.disable,
      startAtBootToggleVisible: false,
      /*----------  Resource autoload  ----------*/
      autoloadContent: autoloadContent,
      autoloadCheckFunc: this.checkAutoloadEnabledState,
      autoloadEnableFunc: this.enableAutoload,
      autoloadDisableFunc: this.disableAutoload,
    }
  },
  methods: {
    checkAutoloadEnabledState(this: any) {
      let autoloadDisabled = this.$store.state.externalContentAutoloadDisabled
      console.log('autoload disabled state is:')
      console.log(autoloadDisabled)
      let autoloadDisabledArrived = this.$store.state
        .externalContentAutoloadDisabledArrived
      if (!autoloadDisabledArrived) {
        return false
      }
      return !autoloadDisabled
    },
    enableAutoload(this: any) {
      fe.SendExternalContentAutoloadDisabledStatus(false, function() {})
    },
    disableAutoload(this: any) {
      fe.SendExternalContentAutoloadDisabledStatus(true, function() {})
    },
  },
}
// These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
var headline = '# Preferences'
var intro = `**You can modify the important main settings here.**

  `

var startAtBootContent = `
## Start at boot

* Aether is designed to be running in the background while your computer is open. It is similar to an email client like ***Outlook***, or ***Mail.app***.

* This is useful, because that means you will receive new updates, posts, and replies to your comments as they happen. It also lets your node discover your peers more efficiently, and track the network better.

* If you disable it, we won't be able to help you troubleshoot if something goes wrong, because Aether is designed to remain open in the background like a *Bitcoin* node, so it can keep track of the network, and 'transactions' happening in it.

* ***Recommendation:*** Leave / make this enabled.
`

var autoloadContent = `
## Auto-load media from whitelisted sources within the app

* Aether, by default, auto-loads non-text content automatically from a small list of whitelisted providers. This whitelist currently includes three items: ***Imgur***, ***Giphy***, ***Gfycat***. *(This list will become editable in the future.)*

* Having it enabled is useful, because it allows you to view the images, videos and gifs within the app even if this content is not carried within Aether. *(Remember: Aether only carries text.)*

* However, loading content from a remote server means that they will get your IP address. You might not want that, especially if you don't trust the server.

* If you don't want your client to autoload assets from the whitelisted sites above, you can disable autoload completely. This will make it so that you have to click to see the content, and the content will open in your browser, instead of within Aether.

* If you don't like the previews, or you just want to use it similar to an email client, you might want to disable it. If you're very concerned about privacy from the whitelisted websites above, you might want to disable it. Otherwise having it enabled is nicer.

* It might also help with RAM use if you're RAM constrained, because images and videos being autoloaded does use some memory.

* ***Recommendation:*** Your call. Having this disabled is more email-like. Having it enabled is more similar to Reddit / Slashdot.

`
</script>

<style lang="scss" scoped>
@import '../../../scss/globals';
.settings-sublocation {
  color: $a-grey-600;
  .start-at-boot-toggle-container {
    display: flex;
    margin-bottom: 20px;
    margin-top: -20px;
    .start-at-boot-toggle {
      font-family: 'SSP Bold';
      flex: 1;
      margin-left: 40px;
    }
    p {
      margin: 0;
      margin-left: 10px;
      margin-right: 35px;
      i {
        font-family: 'SSP Regular Italic';
        font-style: normal;
      }
    }
  }

  .markdowned {
    &:first-of-type {
      margin-bottom: 0;
    }
    margin-bottom: 40px;

    &.settings-bottom-markdown {
      margin-bottom: 0;
    }
  }
  .advanced-button {
    font-family: 'SSP Bold';
    margin-bottom: 15px;
  }
  .config-location {
    padding: 16px 20px;
    color: $a-grey-800;
    background-color: #151618; // $mid-base * 0.5
    margin-bottom: 15px;
    border-radius: 3px;
    .config-location-header {
      padding-bottom: 10px;
      font-family: 'SSP Bold';
      font-size: 17.6px;
    }
    .config-location-text {
      font-family: 'SCP Regular';
      word-break: break-all;
    }
  }
}
</style>
<style lang="scss">
.settingsscope .markdowned {
  tr > th {
    text-align: left;
    &:last-of-type {
      text-align: right;
    }
  }
  tr > td {
    text-align: left;
    &:last-of-type {
      text-align: right;
    }
  }
}
</style>
