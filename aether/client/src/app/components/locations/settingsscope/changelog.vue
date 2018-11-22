<template>
  <div class="settings-sublocation">
    <div class="changelog-successful-update" v-if="firstRunAfterUpdate">
      <div class="changelog-header">
        Hey there! Your app just got new tricks.
      </div>
      <div class="changelog-text">
        You can find what changed in the changelog below. Cheers 🙂<br>
        - B
      </div>
    </div>

    <a-markdown :content="headline"></a-markdown>
    <!-- <a-markdown :content="intro"></a-markdown> -->
    <div class="patreon-segment">
      <div class="patreon-segment-text">
        Want to keep it improving? You should fund Aether.
      </div>
      <div class="funding-rack">
        <a-patreon-button></a-patreon-button>
        <a-crypto-fund-button></a-crypto-fund-button>
      </div>
      <hr>
    </div>
    <a-markdown :content="content" class="changelog-markdown-wrapper"></a-markdown>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'changelog',
    data() {
      return {
        headline: headline,
        intro: intro,
        content: content,
        firstRunAfterUpdate: false,
      }
    },
    beforeMount(this:any) {
      this.firstRunAfterUpdate = this.determineFirstRunAfterUpdate()
    },
    updated(this:any) {
      this.firstRunAfterUpdate = this.determineFirstRunAfterUpdate()
    },
    methods: {
      determineFirstRunAfterUpdate(this:any) {
        // Why a separate thing? We want to 'collect' the value and not be affected by its future changes, because we're gonna set it to false immediately if found to be true.
        // console.log('determine first run after update runs.')
        // console.log('fraus start: ',this.$store.state.firstRunAfterUpdate)
        if (this.$store.state.firstRunAfterUpdate) {
          this.$store.dispatch('setFirstRunAfterUpdateState', false)
          // console.log('fraus end: ',this.$store.state.firstRunAfterUpdate)
          return true
        }
        return false
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = '# Changelog'
  var intro =
    `**This is the change history between versions. The most current version of the changelog can be found [here](https://getaether.net/docs/changelog).**`
  var content =
    `
## 2.0.0.dev.6
*Nov 17, 2018*

* This is the first public, *announced* release of Aether. Though this distinction is somewhat lost because Aether's 'soft' launch did not end up very 'soft', with 8000+ people. [Hacker News discussion](https://news.ycombinator.com/item?id=18370208)

* You can get access to the rough roadmap of Aether as I build here: [Aether Rough Roadmap](https://www.patreon.com/posts/aether-rough-22568167). This also allows you to put your suggestions in.

* This release, and likely future releases will be made available to supporters before the public. This is both useful in testing for regressions, and it's also just a nice thing to do for backers, as a thank-you for their continued support.

### Improvements


#### Deep links in user profile
In user profiles, posts and threads are now deep linked to where they come from, so you can go to a user's profile, and navigate to the threads that s/he responded in.

#### Sort by subscriber count
The communities are now sorted by user count in both the communities list, and on the subscriptions list on the sidebar.

#### Visible communities and dates
On Popular and Home views, the community and the posting date is now visible on each thread.

#### Drafts
When you want to check something else for a sec while you're writing something, your post is now automatically saved as a draft. Whenever you come back, your draft will be there. Drafts persist while the app is open in the tray / menubar. If you want to clear all the drafts, just restart the app.

#### Accommodation for much smaller screens
Design improvements in the Preferences, Status, Onboarding process and in post, thread and community views to ensure that Aether is functional in much smaller screens.

#### Vote states now visible
When a vote is locked (you change your vote once), it is now visible with a lock icon. [Why can I change my vote once?](https://getaether.net/docs/faq/changing_votes/)

#### Better notifications

- Notifications now show the name of the person who responded to you.

- Notifications icon will now stay lit up so long as there are notifications to read. (It used to go into a no-notification state when it was clicked, even when there were notifications present. This caused 'forgotten' notifications in some cases.)


#### Other improvements in brief
- You can now see who created every community in the community info.

- The app now retains the window size from the last run.

- Settings menu and dropdown updates to make things a little more clear.

- Breadcrumbs and tabs now show their full title if hovered. This is good especially on smaller screens where they might now have the space for the full text.

- Thread body previews are now clamped to load faster.

- The code blocks in Markdown fields now create a 'code box' block instead of just singular lines of code text.

- Add a console toggle, so that if something goes wrong in the UI, the users can open the console, and send the error message to help with bug fixes.

- Frontend status now has a 'content' field that shows when the content was last updated.

- Minter in the status box on the sidebar is replaces with the last refresh.

- Added a @ sign to the beginning of the username selection field, so that users won't accidentally add the @ themselves to end up with @@username fields.


### Bug fixes

- A potential fix is implemented for some people not seeing Windows native notifications.

- In rare circumstances, if you clicked vote twice on the same direction, the second vote would also mint an update that was superfluous, and discarded.

- When you edited a post, under certain rare circumstances, it showed like it had two copies on the same thread.

- If you had a unique username, the minting-in-progress state of a post incorrectly showed your username as non-unique while the minting process was ongoing.

- Having a very long (60.000+ characters long) body text on a thread caused performance degradation in certain views, because it was loaded into the DOM in entirety when all that was needed was the first few lines. This is now fixed by only using the first few hundred characters to generate the preview.

- In onboarding, with very small computer screens, the 'next' buttons were left below the fold, preventing people from moving forward.

- Popular and Home views were not correctly applying the moddelete commands when they were generating themselves.

### Backend improvements

- Implemented an incremental, exponentially sized caches and cache consolidation, so that the bootstrap node can be fully static, yet still current.

- Reduced the neighbourhood size defaults for new users to reduce network chatter.

- Reverse opens now use a second, separate port, which is your Aether port minus one. This improves the success rate of reverse opens and it prevents the blocking of TCP.Accept(), which is single-threaded.

- Fixed a socket leak in the bootstrap server related to the issue above, so that the servers are no longer dying every few hours due to socket exhaustion.

### Documentation improvements

- Added a FAQ entry on why votes can change direction once. [See it here.](https://getaether.net/docs/faq/changing_votes/)

- Put the Mim protocol documentation online. [See it here.](https://getaether.net/docs/developers/)

- Added a troubleshooting entry on how to delete the user profile, which is useful in the case a user needs to start from scratch. [See it here.](https://getaether.net/docs/troubleshooting/resetting/)

- Added a few more FAQ and Troubleshooting entries about common questions and issues. [FAQ index](http://localhost:1313/docs/faq/), [Troubleshooting index](https://getaether.net/docs/troubleshooting/)

- Changelogs are now publicly available at the website. [Changelog](https://getaether.net/docs/changelog)

## 2.0.0.dev.5
*Nov 2, 2018*

* This is the first public release of Aether.

* Fixed a bug where the cache generation timestamp was set on attempts, and not successful cache generations. This caused some deliverability issues, since this timestamp controlled the oldest entity that can be provided in a POST response.

* Fixed the bug where the user counts on boards were incorrect for new users.

* Made logs quieter. If you want to see logs, you can set the logginglevel to 1 or 2, and start the app from the command line.

* First-ever cache generation will now happen after bootstrap. This should ease the load on new users coming into the network.

* Preferences and status locations are now centered in the app.

* Miscellaneous stability fixes.

## 2.0.0.dev.5
*Nov 2, 2018*

* This is the first public release of Aether.

* Fixed a bug where the cache generation timestamp was set on attempts, and not successful cache generations. This caused some deliverability issues, since this timestamp controlled the oldest entity that can be provided in a POST response.

* Fixed the bug where the user counts on boards were incorrect for new users.

* Made logs quieter. If you want to see logs, you can set the logginglevel to 1 or 2, and start the app from the command line.

* First-ever cache generation will now happen after bootstrap. This should ease the load on new users coming into the network.

* Preferences and status locations are now centered in the app.

* Miscellaneous stability fixes.


## 2.0.0.dev.4
*Oct 27, 2018*

* Auto update is implemented. This is the first release auto update should work on Windows and Mac.

* The build kill switch is removed. This is the first release that can be used without a time restriction.

* There are now 'copy link' support on communities, threads, posts, and users. Those links will produce aether:// links and can be posted online.

* Other minor changes.

* This release is not public. Future developer previews will be, but this one is too likely to break.

-------
## 2.0.0.dev.2 / 2.0.0.dev.3
*October, 2018*

* Too many changes to count.

-------

## 2.0.0.dev.1
*September 29, 2018*

* Private developer preview release.

* Developer preview releases are not intended for normal, non-technical users.

* This is a testing platform. Anything can go away, anything can disappear, anything can crash, anything can break. If you write your next great novel in Aether, remember to save a copy.

* This release is not public. Future developer previews will be, but this one is too likely to break.
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
    .changelog-successful-update {
      background-color: rgba(0,0,0,0.4);
      padding: 15px 20px;
      margin-bottom: 15px;
      border-radius:3px;
      .changelog-header {
        font-family: "SCP Bold";
        margin-bottom:15px;
      }
      .changelog-text {
        font-family: "SCP Regular";
      }
    }
    .patreon-segment {
      margin-bottom:20px;
      .patreon-segment-text {
        margin-bottom: 20px;
      }
      .patreon-button {
        margin-bottom: 20px;
      }
      hr {
        background-color: rgba(255, 255, 255, 0.25);
        height: 3px;
        border: none;
      }
      .funding-rack {
        display:flex;
        .crypto-fund-button {
          margin-left:10px;
        }
      }
    }
  }
</style>
<style lang="scss">
  .changelog-markdown-wrapper.markdowned {
    h2 {
      margin-bottom: 0;
    }
    &>h2+p { // All direct descendants of above that are immediately after a h2. i.e. this is the date line.
      margin-top:0;
      margin-bottom:0;
    }
  }
</style>