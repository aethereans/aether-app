<template>
  <div
    class="hashimage"
    :style="
      'background-image: ' +
      dataurl +
      ';' +
      ' height: ' +
      height +
      ';' +
      ' width: ' +
      height +
      ';'
    "
  ></div>
</template>

<script lang="ts">
const hashimage = require('../services/hashimage/hashimage')
export default {
  name: 'a-hashimage',
  props: ['hash', 'isUser', 'height'],
  data() {
    return {}
  },
  computed: {
    dataurl(this: any) {
      return 'url(' + hashimage(this.hash, this.isUser) + ')'
    },
  },
}
</script>

<style lang="scss" scoped>
.hashimage {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background-size: cover;
  background-repeat: no-repeat;
}
</style>
