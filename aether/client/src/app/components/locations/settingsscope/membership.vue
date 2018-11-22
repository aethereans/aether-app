<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <!-- <a-markdown :content="intro"></a-markdown> -->
    <a-markdown :content="content"></a-markdown>
    <div class="sponsorship-options">
      <div class="funding-rack">
        <a-patreon-button></a-patreon-button>
        <a-crypto-fund-button></a-crypto-fund-button>
      </div>
    </div>
    <a-markdown :content="content2"></a-markdown>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'membership',
    data() {
      return {
        headline: headline,
        // intro: intro,
        content: content,
        content2: content2,
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = '# Supporter benefits'
  // var intro =
    // `****`
  var content =
    `
*(This page is a summary of [this](https://getaether.net/supporterbenefits))*

## Why?

* Aether is open source and given away for free. It's meant to be ***Internet's own thing***, versus every other public venue on the Internet being owned by corporations.

* Aether is work of one person working full time on it. He would very much prefer to keep doing it. (hey there! 🙂 )

* If you'd like to see Aether thrive, you can fund it via Patreon or crypto (BTC, BCH, ETH, LTC). This will extend the amount of time that I can work full time on it. Regardless of the method, you get the same benefits.

## Supporter benefits

* You get an exclusive, unique username registration.

    **For example**, If the user name \`\`\`@alice\`\`\` is registered, there can be no other registered \`\`\`@alice\`\`\`s. It also turns your user name orange, and it highlights that you're supporting Aether.

* If you have a username you always wanted to have, this is a great way to have it before anyone else does, and support the project in the meanwhile.

`
  var content2 = `
## Company sponsor benefits

* Recurring pledges come with exclusive perks, e.g. having your name listed in the GitHub repository, or have your company logo placed on the website and the app.

* If you're a company and you'd like to sponsor the development of Aether, please [get in touch](mailto:sponsorships@getaether.net).

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

  .sponsorship-options {
    margin-bottom: 30px;
    a {
      // display:inline-block;
      // line-height:0;
      @extend %link-hover-ghost-extenders-disable;
    }
    .funding-rack {
      display:flex;
      .crypto-fund-button {
        margin-left:10px;
      }
    }
  }
</style>