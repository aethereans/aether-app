<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content"></a-markdown>
    <div class="sponsorship-options">
      <div class="option-container">
        <div class="funding-rack">
          <a-patreon-button></a-patreon-button>
          <a-crypto-fund-button></a-crypto-fund-button>
        </div>
      </div>
    </div>
    <a-markdown :content="mintingSection"></a-markdown>
    <template v-if="mintingState==='idle'">
    <div class="unique-name-composer">
      <a-composer :spec="uniqueNameComposerSpec"></a-composer>
    </div>
    </template>
    <template v-if="mintingState==='inprogress'">
        <a-spinner :hidetext="true"></a-spinner>
    </template>
    <template v-if="mintingState==='complete'">
      <div>
        <h3>Minting complete.</h3>
        The newly-created object pushed to the database is as follows: <br>
        <code>
          {{mintedUniqueNameResp}}
        </code>
      </div>
    </template>
    <template v-if="mintingState==='failed'">
      <div>
        <h3>Minting failed.</h3>
        Error message: <br>
        <code>
          {{mintedUniqueNameResp}}
        </code>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
const nmsupervisor = require('../../../services/nmsupervisor/nmsupervisor')
const feapiconsumer = require('../../../services/feapiconsumer/feapiconsumer')
export default {
  name: 'namemint',
  data (this:any) {
    return {
      headline: headline,
      intro: intro,
      content: content,
      mintingSection: mintingSection,
      mintingState: "idle", //"inprogress", "complete", "failed"
      mintedUniqueNameResp: "",
      uniqueNameComposerSpec: {
        fixToBottom: true,
        fields: [{
          id: "uniqueUsername",
          visibleName: "Unique username",
          description: "Put the desired username here.",
          placeholder: "",
          maxCharCount: 24,
          heightRows: 1,
          previewDisabled: true,
          content: "",
          optional: false,
          spaceDisabled:true,
        }, {
          id: "targetKeyFp",
          visibleName: "Target key fingerprint",
          description: "This is the fingerprint of the key that's going to receive the unique username.",
          placeholder: "",
          maxCharCount: 64,
          heightRows: 1,
          previewDisabled: true,
          content: "",
          optional: false,
        },{
          id: "expiryTimestamp",
          visibleName: "Expiry timestamp",
          description: "Put a timestamp when it should expire. If no timestamp is given, it's assumed to be expiring in 6 months.",
          placeholder: "",
          maxCharCount: 24,
          heightRows: 1,
          previewDisabled: true,
          content: "0",
          optional: true,
        },{
          id: "password",
          visibleName: "Password",
          description: "Put in the password for the CA private key here.",
          placeholder: "",
          maxCharCount: 256,
          heightRows: 1,
          previewDisabled: true,
          content: "",
          optional: false,
          passwordField: true,
        }],
        commitAction: this.submitUniqueUsernameRequest,
        commitActionName: "SUBMIT",
        cancelAction: function() { history.back() },
        cancelActionName: "CANCEL",
        autofocus: true,
      },
    }
  },
  methods: {
    submitUniqueUsernameRequest(this:any, fields:any) {
      let vm = this
      vm.mintingState = "inprogress"
      console.log('enters submit unique username request')
      nmsupervisor.MintNewUniqueUsername(fields[0].content, fields[1].content, fields[2].content, fields[3].content, function(resp:string){
        vm.mintedUniqueNameResp = resp
        if (resp.indexOf("Command failed:") !== -1) {
          vm.mintingState = "failed"
          return
        }
        feapiconsumer.SendMintedUsername(resp, function() { //resp:any
          vm.mintingState = "complete"
          // console.log('minted username sent, response: ')
          // console.log(resp)
        })
      })
    }
  }
}
var headline = `# Name minter`
var intro = `**This is not for public use, and it will not work without the name minter, which is a separate app.**`
var content = `
## Quick info

- This is the part of the app that I (@b) use to mint unique names for supporters of the project.

- This page facilitates minting of unique names based on the CA private key. It requires that you know the CA password, CA private key, the name to be minted, and the key fingerprint of the recipient.

- This page is just a shell, it relies on an external app to mint the actual unique name. Being able to see the page does not mean you will be able to mint unique names.

- If you'd like to get unique names, you should support Aether on Patreon.
`
var mintingSection = `
  ## New unique name
`

</script>

<style lang="scss" scoped>
  @import "../../../scss/globals";

  .sponsorship-options {
    margin-top: 30px;
    margin-bottom: 30px;
    a {
      // display:inline-block;
      // line-height:0;
      @extend %link-hover-ghost-extenders-disable;
    }
    .funding-rack {
      display:flex;
      .crypto-fund-button {
        margin-left:10px;
      }
    }
  }

  .unique-name-composer {
    margin-top:20px;
    font-family: "SSP Bold";
  }

  code {
    word-break: break-all;
  }
</style>