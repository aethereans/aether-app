<template>
  <div class="info-marker" hasInfomark :title="infoText">
    <div class="info-anchor">
      <icon class="info-anchor-icon" :style="{'width': width+'px', 'height': height+'px', 'fill': fill}" name="question-circle"></icon>
    </div>
  </div>
</template>

<script lang="ts">
  var Tooltips = require('../services/tooltips/tooltips')
  export default {
    name: 'a-info-marker',
    props: ['width', 'height', 'text', 'header', 'fill'],
    data() {
      return {}
    },
    mounted() {
      Tooltips.MountInfomark()
    },
    updated() {
      // Tooltips.MountInfomark()
    },
    computed: {
      infoText(this: any) {
        if (typeof this.text !== 'undefined' && typeof this.header !== 'undefined') {
          let textfield = '<div class="infomark-header">' + this.header + '</div><div class="infomark-body">' + this.text + '</div>'
          return textfield
        }
        if (typeof this.text !== 'undefined') {
          let textfield = '<div class="infomark-body">' + this.text + '</div>'
          return textfield
        }
        return ''
      }
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .info-marker {
    display: inline-block;
    cursor: pointer;
    .info-anchor {
      display: flex;
      .info-anchor-icon {
        margin: auto;
        width: 12px;
        height: 12px;
        fill: $a-grey-600;
      }
    }
  }
</style>