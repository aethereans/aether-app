<template>
  <div class="a-threadname" v-if="isPresent">
    <div class="threadname-container">
      <router-link
        @click.native="$event.stopImmediatePropagation()"
        :to="'/board/' + boardFingerprint + '/thread/' + threadFingerprint"
        >{{ name }}</router-link
      >
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-threadname',
  props: {
    name: {
      type: String,
      default: '',
    },
    boardFingerprint: {
      type: String,
      default: '',
    },
    threadFingerprint: {
      type: String,
      default: '',
    },
  },
  data() {
    return {}
  },
  computed: {
    isPresent(this: any) {
      return this.name && this.boardFingerprint && this.threadFingerprint
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.a-threadname {
  display: inline-block;
  .threadname-container {
    display: flex;
    font-family: 'SSP Bold Italic';
    a {
      color: $a-grey-600;
    }
  }
}
</style>
