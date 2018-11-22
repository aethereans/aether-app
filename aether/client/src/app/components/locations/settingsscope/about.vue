<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content"></a-markdown>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'about',
    data() {
      return {
        headline: headline,
        intro: intro,
        content: content,
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = '# About'
  let versionAndBuild = require('electron').remote.app.getVersion().split('+')
  var intro =
    `**Aether v` + versionAndBuild[0] + `, build ` + versionAndBuild[1] + `**`
  var content =
    `
### What is Aether?

Aether is a network of independent, self-moderated communities. [More info](https://getaether.net)

    `
</script>

<style lang="scss" scoped>
  @import "../../../scss/globals";
  .settings-sublocation {
    color: $a-grey-600;
    .markdowned {
      &:first-of-type {
        margin-bottom: 0;
      }
      margin-bottom: 40px;
    }
  }
</style>