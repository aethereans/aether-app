<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content"></a-markdown>
  </div>
</template>

<script lang="ts">
export default {
  name: 'adminsquickstart',
  data() {
    return {
      headline: headline,
      intro: intro,
      content: content,
    }
  },
}
// These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
var headline = "# Admin's Quickstart"
var intro = `**This is a quick introduction to founding your own community on Aether.**`
var content = `
Thanks for being here! Aether is better because of you. The more communities with diverse views inhabiting this place, the more interesting, insightful, interesting it ends up being. We have space for cats, corgis, humans, motorcycles, \`\`\`python\`\`\`, pythons, cryptocurrencies, **Ć̢̘̳̺͉̯̘̣̠̂̓̔̓͗ͭ͊̿̿́͜ţ̛̫̗̱̗̖̬̱̜͖̬̐͊̌̃̀̅̒̋͛͝͡ͅh̷̤̦̝̟̜̓ͦͥ̿ͬ̐̔͘͜ụ̺̻͍̱̟̦͇̗̦̳̤͑̌͐̆̔͊ͧ̔̈͊̀͞l̸̴͔̣͚͇̯̭ͦ̋̈͌͊͂̇͂̚͟͝h̤̰͕͓̪͓͗ͭ̆̾ͮ̄ͪ͒́̆̉͊̍͜ͅͅȗ̸̮̦̦̣̠̝̫̹̟͈͇̯͚̏ͯ͂̎̚͞ͅ**, and anything in between.

***Heads up***
These points below are intended as an introduction to moderating communities on Aether, and how it differs from other places. These aren't rules — you are free to moderate as you like.

## The prime directive

Be excellent to each other.

## The ***real*** prime directive

Act fast and chomp down on posts that just haven't gotten the basic idea.

* Why? Vast majority of people are genuinely good, pleasant humans. But a very small number of misbehavers can make any place not so great, even when almost all of the people in the place is still well-meaning.

* Vast majority of people also have no concept of just how much effort and vigilance it takes to keep a community a pleasant place, or how efficiently an awful post can ruin a lot of people's days. The work you do is valuable.

* Aether attempts to fix that by making your actions visible, and in the front and center. People **will** see just how much work you're doing to keep your corner of the world an enjoyable place.

## Err on the side of action

* A false negative is much more acceptable than a false positive. Deleting a post that could have stayed is much more benign than not deleting a post that should have been gone.

## People have control, as well as you

* The reason you can be much more swift in moderating in Aether communities is that ***people can see what you're doing*** and give you feedback.

* This feedback loop is missing in other types of communities, because if something gets deleted, it's gone - users can't be aware of the absence of something. In Aether, if they flip up the right switch, they can see deleted items, and why you deleted them.

* That means if you hear more than the usual amount of complaints about a specific deletion, it might be worth a second look.

* The mod queue (where reported content shows up) is visible to everyone.

* If someone doesn't like a specific mod action you took, they can revert that action for themselves.

## Aether communities are democracies...

* As the creator of a community, you are the creator mod. That means only you can control what gets into the community info.

* However, other people can elect regular mods in that community. If this happens, they will be regular mods so long as the vote that elected them stays valid.

* You can assign mods as well. These assigned mods don't need to be elected, and they stay as long as you want.

## ... and it cuts both ways

* Any mod can get impeached if it collects enough impeachment votes based on the percentage of population of that community. If that happens, that mod stays impeached as long as the vote that impeached them stays valid. When the impeachment vote falls below the threshold, the mod will regain your privileges automatically.

* This probably will never happen to you, considering that you're reading this.

## Further reading

The good folks over at Discourse has also written a helpful article on how to build and maintain [communities](https://blog.discourse.org/2014/08/building-a-discourse-community/).

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
  }
}
</style>
