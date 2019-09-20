<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown
      class="settings-bottom-markdown"
      :content="content1"
    ></a-markdown>
    <div class="settings-action-container">
      <div class="current-state">
        <div class="current-state-header">
          Mod mode status
        </div>
        <div
          class="current-state-text"
          v-show="modModeEnabled"
          @click="enableModMode"
        >
          Enabled
        </div>
        <div
          class="current-state-text"
          v-show="!modModeEnabled"
          @click="disableModMode"
        >
          Disabled
        </div>
      </div>
      <div class="flex-spacer"></div>
      <div class="button-carrier">
        <a
          class="button is-success is-outlined"
          @click="enableModMode"
          v-show="!modModeEnabled"
        >
          ENABLE MOD MODE
        </a>
        <a
          class="button is-success is-outlined"
          @click="disableModMode"
          v-show="modModeEnabled"
        >
          DISABLE MOD MODE
        </a>
      </div>
    </div>
    <a-markdown :content="content2"></a-markdown>
  </div>
</template>

<script lang="ts">
var fe = require('../../../services/feapiconsumer/feapiconsumer')
export default {
  name: 'modship',
  data() {
    return {
      headline: headline,
      intro: intro,
      content1: content1,
      content2: content2,
    }
  },
  computed: {
    modModeEnabled(this: any) {
      if (
        this.$store.state.modModeEnabledArrived &&
        this.$store.state.modModeEnabled
      ) {
        return true
      }
      return false
    },
  },
  methods: {
    enableModMode(this: any) {
      console.log('enable mod mode is called')
      fe.SendModModeEnabledStatus(true, function(resp: any) {
        console.log(resp)
      })
    },
    disableModMode(this: any) {
      console.log('disable mod mode is called')
      fe.SendModModeEnabledStatus(false, function(resp: any) {
        console.log(resp)
      })
    },
  },
}
// These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
var headline = '# Mod mode'
var intro = `**This place allows you to set the mod mode, and tells about how moderation works in Aether.**`
var content1 = `
* **Mod mode allows you to act as a mod.** You should enable this mode if you've created any communities, so you can moderate them.

* You should also enable this mode if you want to be an elected mod.

* This mode allows you to take mod actions on ***all of the network***, not just the community you are, or want to be, a mod of.

* Your mod actions will only be effective if you have mod privileges in that community.
`
var content2 = `

### What enabling mod mode does

* Enables the *Reports* tab in the community home page.

  - With this, you can see what threads or posts people are reporting to mods and act on it.

  - It will only show reports that are created after mod mode is enabled.

  - This wil make your app use slightly more disk space, since this list has to be compiled by the frontend separately

* Makes \`\`\`Moddelete\`\`\` and \`\`\`Modapprove\`\`\` actions visible on every thread and post.
  *(\`\`\`Modapprove\`\`\` is still in development.)*

### What enabling mod mode ***doesn't*** do

* It does not give you any privileges that you don't already have

-----

## FAQ

### Things a mod can do

* Delete threads and posts

* Pin posts to top *(This feature is still in development.)*

### Things a mod ***can't*** do

* Edit threads and posts

-----

### Getting mod privileges

* You can get privileges in any community by:

  - Creating that community,

  - Being assigned a mod by the creator of the community,

  - Being elected as a mod by users in that community.

* Your mod actions will also be effective to the specific people who:

  - Voted to elect you as a mod (even if that vote wasn't won).


* Your mod actions will be effective to no one else. They also won't be effective to the specific people that voted to impeach you (even if that vote wasn't won).

-----

### Mod action visibility

* **Your mod actions are visible to everyone.** That's how a person can look at your modding history and decide to vote for you. When you're deleting or approving content, consider putting in why on the *Reason* field.

* If somebody decides that you're not cut out to be a mod, they can vote to impeach you.

* If somebody votes for impeachment, you cease to be a mod in the eyes of that person, and all your mod actions will be reverted.

-----

### Creator mod

* Creator mod is the creator of that community. It is a regular mod, and can be impeached as well as any other mod.

* It is subject to all the same rules, with one additional ability. Creator mod can change the community's info page.

-----

### What to do if you want to be an elected mod

* Enable mod mode and start moderating content.

* Your name will pop up as one of the candidates for modship in that community's elections section. It will also show the number of mod actions you've taken in the recent past.

* People can click your name, and see your mod actions, and can vote for or against you.

* If 5% of the community's population has voted for you, and of those votes, at least 51% is positive, you will be elected as a mod, and your mod actions will be applied to everyone by default.

-----

### Philosophy

* Communities in Aether are sovereign and self-governed.

* If a mod strays too far off what the community wants, he/she will be impeached.

* Impeachment will last so long as the impeach vote stays above 5% of total population and at least 51% negative.

* If there are people in the community that are interested in moderation (or moderating just to curate their own views, even), they can get elected mods.

* Given a conflict (Mod A deletes a piece of content, mod B approves it), Aether will default to making it visible. One \`\`\`modapprove\`\`\` is stronger than an infinite number of \`\`\`moddelete\`\`\`s.

* Generally speaking, if a community elects a lot of moderators, and the original mod team is impeached, it is time to start a new community.

* If you can get elected in a community, you probably have the clout to ask people to move to your new community where you are the creator mod, and have full control, like editing the community info.

* To prevent hostile community takeovers, the population of the community is determined on a two week trailing basis (new posters don't count until their posted content is as least two weeks old), and the mods can decide to temporarily *lock* a community to new entrants until the attempt is over. *(This feature is still in development.)*

`
</script>

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
  .settings-action-container {
    background-color: rgba(0, 0, 0, 0.25);
    padding: 15px 20px; // margin: auto;
    font-family: 'SSP Bold';
    margin-bottom: 20px;
    border-radius: 3px;
    display: flex;

    .current-state {
      font-family: 'SCP Regular';
      margin-bottom: 10px;
      .current-state-header {
        font-family: 'SCP Bold';
      }
    }

    .button-carrier {
      padding-top: 6px;
    }
  }
}

.flex-spacer {
  flex: 1;
}
</style>
