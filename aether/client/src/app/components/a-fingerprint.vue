<template>
  <div class="a-fingerprint">
    <div v-for="rowblocks in rows" class="rowblock">
      <div v-for="rowblock in rowblocks" class="block" :style="{fontSize:fontSize+'px'}">
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
        default: "",
      },
      fontSize: {
        type: Number,
        default: 14,
      }
    },
    data() {
      return {}
    },
    computed: {
      rows(this: any): any {
        let vm = this
        // row1
        let row1 = vm.fp.substr(0, vm.fp.length / 2)
        let row1blocks: any = []
        for (var i = 0; i < row1.length; i++) {
          (function(i) {
            if (i % 4 === 0) {
              row1blocks.push(row1[i])
            } else {
              row1blocks[row1blocks.length - 1] += row1[i]
            }
          })(i)
        }
        for (var i = 0; i < row1blocks.length; i++) {
          (function(i) {
            row1blocks[i] = row1blocks[i].replace(/(\d+)/g, '<span class="number">$1</span>')
          })(i)
        }
        // row2
        let row2 = vm.fp.substr(vm.fp.length / 2, vm.fp.length)
        let row2blocks: any = []
        for (var i = 0; i < row2.length; i++) {
          (function(i) {
            if (i % 4 === 0) {
              row2blocks.push(row2[i])
            } else {
              row2blocks[row2blocks.length - 1] += row2[i]
            }
          })(i)
        }
        for (var i = 0; i < row2blocks.length; i++) {
          (function(i) {
            row2blocks[i] = row2blocks[i].replace(/(\d+)/g, '<span class="number">$1</span>')
          })(i)
        }
        // console.log(row1blocks)
        // console.log(row1)
        // console.log(row2blocks)
        // console.log(row2)
        return [row1blocks, row2blocks]
      }
    }
  }
</script>
<style lang="scss">
  @import "../scss/globals";
  .a-fingerprint .rowblock .block .number {
    color: $a-yellow;
  }
</style>
<style lang="scss" scoped>
  @import "../scss/globals";
  .a-fingerprint {
    font-family: "SCP Bold";
    background-color: rgba(0, 0, 0, 0.25); // padding: 5px 10px;
    padding: 0.4em 0.4em;
    border-radius: 3px; // line-height: 140%;
    display: inline-block;
    .rowblock {
      display: flex;
      .block {
        padding: 0 0.3em;
        .number {
          color: blue;
        }
      }
    }
  }
</style>