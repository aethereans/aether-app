<template>
  <div class="a-fingerprint" @copy.stop.prevent="copyFp">
    <div class="meta-box">
      <a
        class="button is-outlined is-small is-link copy-link-button no-outline copy-link-button"
        @click="copyLink"
        title="Copy link"
      >
        <icon name="link" class="copy-link-icon"></icon> Copy link
      </a>
      <div class="link-copied-text" v-if="linkCopied">
        <span>Link copied to clipboard</span>
      </div>
      <div class="flex-spacer"></div>
      <div class="infomark-container">
        <a-info-marker
          header="Fingerprint above is the unique URL of this content or user."
          text="<p></p><p>Copying the link will put the link to this entity to your clipboard. This link works inside, and outside of Aether (i.e. posted on the web) if the recipients have Aether installed.</p>"
        ></a-info-marker>
      </div>
    </div>
    <div v-for="rowblocks in rows" class="rowblock">
      <div v-for="rowblock in rowblocks" class="block">
        <!-- :style="{ fontSize: fontSize + 'px' }" -->
        <span v-html="rowblock"></span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-fingerprint',
  props: {
    fp: {
      type: String,
      default: '',
    },
    fontSize: {
      type: Number,
      default: 14,
    },
    link: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      linkCopied: false,
      copyLinkInProgress: false,
    }
  },
  computed: {
    rows(this: any): any {
      let vm = this
      // row1
      let row1 = vm.fp.substr(0, vm.fp.length / 2)
      let row1blocks: any = []
      for (var i = 0; i < row1.length; i++) {
        ;(function (i) {
          if (i % 4 === 0) {
            row1blocks.push(row1[i])
          } else {
            row1blocks[row1blocks.length - 1] += row1[i]
          }
        })(i)
      }
      for (var i = 0; i < row1blocks.length; i++) {
        ;(function (i) {
          row1blocks[i] = row1blocks[i].replace(
            /(\d+)/g,
            '<span class="number">$1</span>'
          )
        })(i)
      }
      // row2
      let row2 = vm.fp.substr(vm.fp.length / 2, vm.fp.length)
      let row2blocks: any = []
      for (var i = 0; i < row2.length; i++) {
        ;(function (i) {
          if (i % 4 === 0) {
            row2blocks.push(row2[i])
          } else {
            row2blocks[row2blocks.length - 1] += row2[i]
          }
        })(i)
      }
      for (var i = 0; i < row2blocks.length; i++) {
        ;(function (i) {
          row2blocks[i] = row2blocks[i].replace(
            /(\d+)/g,
            '<span class="number">$1</span>'
          )
        })(i)
      }
      // console.log(row1blocks)
      // console.log(row1)
      // console.log(row2blocks)
      // console.log(row2)
      return [row1blocks, row2blocks]
    },
  },
  methods: {
    copyFp(this: any) {
      let elc = require('electron')
      elc.clipboard.writeText(this.fp)
    },
    copyLink(this: any) {
      if (this.copyLinkInProgress) {
        return
      }
      this.copyLinkInProgress = true
      let elc = require('electron')
      elc.clipboard.writeText(this.link)
      this.linkCopied = true
      // Prevent repeated clicks while in progress
      let vm = this
      setTimeout(function () {
        vm.linkCopied = false
        vm.copyLinkInProgress = false
      }, 1250)
    },
  },
}
</script>
<style lang="scss">
@import '../scss/globals';
.a-fingerprint .rowblock .block .number {
  color: $a-yellow;
}
</style>
<style lang="scss" scoped>
@import '../scss/globals';
.a-fingerprint {
  font-family: 'SCP Bold';
  background-color: rgba(0, 0, 0, 0.25); // padding: 5px 10px;
  padding: 0.4em 0.4em;
  border-radius: 3px; // line-height: 140%;
  display: inline-block;
  font-size: 100%;
  .rowblock {
    display: flex;
    .block {
      padding: 0 0.3em;
      .number {
        color: blue;
      }
    }
  }
  .meta-box {
    display: flex;
    padding: 5px 5px 10px 3px;
    .copy-link-button {
      padding: 0 3px;
      .copy-link-icon {
        width: 12px;
        height: 12px;
        margin-right: 8px;
      }
    }
    .flex-spacer {
      flex: 1;
    }
    .infomark-container {
      display: flex;
      .info-marker {
        margin: auto;
      }
    }
    .button.no-outline {
      border-color: transparent;
    }
    .copy-link-button {
      color: $a-grey-800;
    }
    .link-copied-text {
      font-size: 0.75rem;
      height: 2.25em;
      display: flex;
      animation-duration: 1.25s;
      animation-name: DELAY_INVISIBLE;
      padding-left: 10px;
      span {
        margin: auto;
        user-select: none;
        cursor: default;
      }
    }
    @keyframes DELAY_INVISIBLE {
      0% {
        opacity: 1;
      }
      60% {
        opacity: 1;
      }
      100% {
        opacity: 0;
      }
    }
  }
}
</style>
