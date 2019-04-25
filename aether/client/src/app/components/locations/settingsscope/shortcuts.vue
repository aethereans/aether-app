<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content1"></a-markdown>
  </div>
</template>

<script lang="ts">
export default {
  name: 'shortcuts',
  data(this: any) {
    return {
      headline: headline,
      intro: intro,
      content1: content1,
    }
  },
}
// These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
var headline = '# Keyboard shortcuts'
var intro = `**This is a list of keyboard shortcuts that you can use in the app.**
  `

var content1 = `

\`\`\` Mod \`\`\` key mentioned below is:
* \`\`\`Command\`\`\` (⌘) on Mac
* \`\`\`Ctrl\`\`\` on Windows
* \`\`\`Super\`\`\` *(usually also Ctrl)* on Linux.

## Global shortcuts

These shortcuts are valid within the app unless a text field is selected.

| Shortut | Description |
|--- | --- |
| \`\`\` Mod + , \`\`\` | Go back in history ( < )
| \`\`\` Mod + . \`\`\` | Go forward in history ( > )
| \`\`\` Mod + / \`\`\` | Open notifications
| \`\`\` Mod + = \`\`\` | Enlarge UI
| \`\`\` Mod + - \`\`\` | Shrink UI
| \`\`\` Mod + 0 \`\`\` | Return UI to default size
| \`\`\` Mod + Alt + I \`\`\` | Open / Close debug window

## Community shortcuts

These shortcuts work when you're in a community

| Shortut | Description |
|--- | --- |
| \`\`\` Mod + N \`\`\` | Create new thread

## Text entry shortcuts

These shortcuts work in any text entry block of one or multiple fields.

| Shortut | Description |
|--- | --- |
| \`\`\` Tab \`\`\` | Next field in block
| \`\`\` Mod + Enter \`\`\` | Submit
| \`\`\` Mod + Esc \`\`\` | Cancel and close
`
</script>

<!--

(Removing these, because they're too complex to ever be used frequently. We should find better shortcuts for them in the future.)

## Thread shortcuts

These shortcuts work when you're in a thread

| Shortut | Description |
|--- | --- |
| \`\`\` Mod + + (Mod + Shift + =) \`\`\` | Upvote thread
| \`\`\` Mod + _ (Mod + Shift + - ) \`\`\` | Downvote thread

 -->

<style lang="scss" scoped>
@import '../../../scss/globals';
.settings-sublocation {
  color: $a-grey-600;
  .markdowned {
    &:first-of-type {
      margin-bottom: 0;
    }
    margin-bottom: 40px;

    &.settings-bottom-markdown {
      margin-bottom: 0;
    }
  }
  .advanced-button {
    font-family: 'SSP Bold';
    margin-bottom: 15px;
  }
  .config-location {
    padding: 16px 20px;
    color: $a-grey-800;
    background-color: #151618; // $mid-base * 0.5
    margin-bottom: 15px;
    border-radius: 3px;
    .config-location-header {
      padding-bottom: 10px;
      font-family: 'SSP Bold';
      font-size: 17.6px;
    }
    .config-location-text {
      font-family: 'SCP Regular';
      word-break: break-all;
    }
  }
}
</style>
<style lang="scss">
.settingsscope .markdowned {
  tr > th {
    text-align: left;
    &:last-of-type {
      text-align: right;
    }
  }
  tr > td {
    text-align: left;
    &:last-of-type {
      text-align: right;
    }
  }
}
</style>
