<template>
  <div class="vote-action">
    <div class="already-voted-highlighter">
      <icon name="arrow-up"></icon>
    </div>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'a-vote-action',
    data() {
      return {}
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";

  .vote-action {
    display: flex;
    margin: auto;
  }

  .already-voted-highlighter {
    margin: auto;
    background-color: $a-cerulean;
    border-radius: 25px;
    display: flex;
    padding: 2px;
    svg {
      flex: 1;
    }
  }
</style>