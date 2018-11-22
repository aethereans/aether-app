<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content"></a-markdown>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'sfwlist',
    data() {
      return {
        headline: headline,
        intro: intro,
        content: content,
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = "# Safe-for-work communities list"
  var intro =
    `**Preferences related to the safe-for-work markings for communities.**`
  var content =
    `
* SFW communities are communities that are chosen to be work-safe, and they contain some of the higher-quality discussion and content the network has available.

* A board being not SFW is not a negative signal. It might mean that the board was just formed, or that it is interesting and high quality but not SFW, or that it does not have the volume or active moderation.

## Disabling the safe-for-work list

* If you want to ***always*** see non-SFW communities on the communities list, set \`\`\`SFWListDisabled\`\`\` flag to true in frontend configuration. This will be saved to your profile and will carry over to any device you move your profile to.

* **Leaving the safe-for-work list enabled does not prevent you from seeing or interacting with non-SFW communities.** It just keeps them in a separate list at the bottom of the Communities tab. This list does not affect your subs. It does not affect anything else except the order of the communities list, and what goes into the Popular page, for communities you don't subscribe to.

* If you subscribe to a non-SFW community, you will always see content from it regardless of your safe-for-work list settings.

------
## FAQ


### I just started a new community. It is SFW. How do I get it added to the list?

\`\`\`
(For a limited time, you can get your community added to the SFW list instantly. Just email sfwlist@getaether.net with a link to your community. Mind that your community should actually be SFW.)
\`\`\`

Your community should to be created at least 30 days ago, and it should not be empty or near-empty.

Send an email to <b>sfwlist@getaether.net</b> with the following details:

1. Your community's fingerprint (available at YourCommunity > Info)

2. A valid email address from a reputable email host that you can be reached from within a reasonable timeframe.


* Safe-for-work list inclusion is not a value judgment. It's just a way for you to certify that you want your community to be SFW, and you intend to keep it so.

------

### What exactly does safe-for-work list do?

* It makes the listed communities to be visible at the top of the communities list. *(Non-listed communities will still be visible as a second list at the bottom, with a toggle.)*

* It controls which communities will be considered for addition in Popular view alongside your subbed communities. (Home view is your subs, Popular view is your subs *and* the safe-for-work list.)

* If it is disabled, the app acts as if the whole network is safe-for-work listed.

------

### Why does it exist?

* Since there is no way to know if a community is non-SFW, and there is no way to assign that tag via someone else who's not a mod (because communities are sovereigns), it's best to treat all communities as non-SFW at first, and add the safe-for-work ones to a known list.

* **Aether aims to be for everyone.** It is much better user experience to keep the default to safe-for-work and allow people who want to see more to explicitly enable it. Chances are, even if you're OK with non-SFW content, you might not necessarily want to see it on your work computer within the work hours.
------
### Is this a part of the protocol and/or the network?

No. The safe-for-work is not a part of any protocol (Mim) nor any subprotocol (c0) specification. It is an optional feature of the Aether app, which is a client app for c0 subprotocol of the Mim protocol.

If you are building your own client for c0, you are welcome to use this safe-for-work list, but you are not required to. In the case you use it, it is recommended that you leave it optional, and not rely on it as a value judgment.

**The safe-for-work list is provided as is on a best-effort basis as a service to the community, and it does not come with any guarantee of accuracy, impartiality or freshness.**

------
### Who makes the decision on what to include?

* **Currently:** The person behind the Aether app.

* **Future:** The lists feature becomes generic, and it allows import and update of custom community-maintained community lists for many purposes (interest based, country based, etc), and safe-for-work list is one of the many different lists that can be chosen from, and it becomes maintained by impartial observers.
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