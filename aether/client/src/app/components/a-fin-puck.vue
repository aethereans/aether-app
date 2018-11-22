<template>
  <div class="fin-carrier">
    <div class="fin-text" hasMarker title='This is the end  <br>beautiful friend  <br> this is the end  <br> my only friend, the end <br>— Jim Morrison'>
      END
    </div>
  </div>
</template>

<script lang="ts">
  var Tooltips = require('../services/tooltips/tooltips')
  export default {
    name: 'a-fin-puck',
    data() {
      return {}
    },
    mounted() {
      Tooltips.MountMarker()
    },
    updated() {
      Tooltips.MountMarker()
    }
  }
</script>

<style lang="scss" scoped>
  .fin-carrier {
    display: flex;
    padding: 20px 0 0 0;
    .fin-text {
      margin: auto;
      background-color: rgba(255, 255, 255, 0.25);
      font-family: "SCP Bold";
      padding: 5px;
      line-height: 100%;
      border-radius: 2px;
      color: #292b2f;
      cursor: default;
      user-select: none;
      font-size: 90%;
      letter-spacing: 1px;
    }
  }
</style>