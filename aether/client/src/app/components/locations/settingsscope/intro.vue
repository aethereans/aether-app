<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <!-- <a-markdown :content="intro"></a-markdown> -->
    <img class="header-image" src="../../../ext_dep/images/nasa-89125-unsplash.jpg" alt="">
    <a-markdown :content="content"></a-markdown>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'intro',
    data() {
      return {
        headline: headline,
        intro: intro,
        content: content,
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = `# A Beginner's Guide to the Galaxy`
  var intro =
    `**(You don't have to read it in entirety, but it might be interesting)**
  `
  var content =
    `
## Welcome in!

* Aether is an ephemeral, auditable and privacy-focused network of peer-to-peer, self-moderated communities.

* (BTW, you don't have to read this in entirety to use this thing, it's just a collection of bullet points.)

### How it works

* Every computer participating in the network (that's you!) forwards the content it receives from the network to new users, that way, we don't have a need for centralised servers. It also improves the privacy for every one of us.

* That means, because there are no servers, the data is just shared between the participants, and not stored specifically anywhere. This works well because ***all the content on the network is just text***, no images, no videos, etc. Those can be linked, but if you click on them, they would open in your regular browser.

* For that reason, the app is more like an email app. To be able to get you the new stuff, it stays running in the background (on your taskbar or menubar) and keeps talking to the network, so that it'll be ready whenever you pop it open.

* The longer your computer stays online in the network, the more nodes that will get a chance to connect to you, and the faster the posts you create will be spread to the network. (You might want to shut down the app when you're tethering or when you're on an otherwise metered connection.)

## Good-to-knows

### Usernames and community names

* One thing that's good to know is that the ***usernames are not unique*** by default. So that means there can be two \`\`\`@alice\`\`\`s. If you want to discern one from another, hover over the username and it'll show you a pixelated picture and a Fingerprint. Even if the names are the same, those will be different. Same with boards, their names can be the same, just check the Fingerprint and the picture.

  **What's up with orange usernames?**
  Those are the usernames of the supporters of the project on Patreon. Those usernames are unique, so there can only ever be one single registered \`\`\`@alice\`\`\`. If you like to learn more about this and register your unique name, head to the [Supporter benefits](#/membership).

### Content will eventually disappear

* Aether is an ephemeral network. That means, by default, most things will be deleted around the 6-month mark. It's more like Snapchat than Wikipedia.

* What you post in will eventually be gone unless somebody somewhere decides to keep it. If you are especially interesting, somebody might do that, but for most of us, most of the stuff will be gone eventually.

* This is a *Good Thing*: nobody wants to read their own tweets from when they were 14 years old. To be human is to forget, forgive and move on.

### Privacy oriented by default

* There is no real-name policy, you can use whatever the name you want, including your real name. The underlying idea is that the basic assumption is privacy, so you can choose how identifiable you want to be.

### ... Wait, who's paying for this?

* I am. (This is [\`\`\`@b\`\`\`](https://twitter.com/nehbit) here.) I've been working full time on this since March 2018. If you want to support Aether, you can do so on [Patreon](https://patreon.com/nehbit). It also gives you some [sweet perks](#/membership).
  `
  // var content: string = require('./mdpwnattempt')
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

  .header-image {
    width: 100%;
    border-radius: 2px;
    overflow: hidden;
    margin-bottom: 25px;
  }
</style>

<style lang="scss">
  h1#a-beginner-s-guide-to-the-galaxy {
    margin-bottom: 10px;
  }
</style>