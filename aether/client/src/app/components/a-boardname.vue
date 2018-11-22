<template>
<div class="a-boardname" v-if="isPresent">
  <div class="boardname-container">
    <!-- <a-hashimage :hash="fingerprint" height="16px"></a-hashimage> -->
    <router-link :to="'/board/'+fingerprint">b/{{name}}</router-link>
  </div>
</div>
</template>

<script lang="ts">
export default {
  name: 'a-boardname',
  props: {
    name: {
      type: String,
      default: ""
    },
    fingerprint: {
      type: String,
      default: ""
    }
  },
  data () {
    return {

    }
  },
  computed: {
    isPresent(this:any) {
      return this.name && this.fingerprint
    }
  }
}
</script>

<style lang="scss" scoped>
@import "../scss/globals";
.a-boardname {
  display:inline-block;
  .boardname-container {
    display:flex;
    font-family: "SSP Bold Italic";
    a {
      color: $a-grey-600;
    }
  }
}
</style>