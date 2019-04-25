<template>
  <div class="settings-sublocation">
    <div class="changelog-successful-update" v-if="firstRunAfterUpdate">
      <div class="changelog-header">
        Hey there! Your app just got new tricks.
      </div>
      <div class="changelog-text">
        You can find what changed in the changelog below. Cheers 🙂<br />
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
        <!-- <a-crypto-fund-button></a-crypto-fund-button> -->
      </div>
      <hr />
    </div>
    <a-markdown
      :content="content"
      class="changelog-markdown-wrapper"
    ></a-markdown>
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
  beforeMount(this: any) {
    this.firstRunAfterUpdate = this.determineFirstRunAfterUpdate()
  },
  updated(this: any) {
    this.firstRunAfterUpdate = this.determineFirstRunAfterUpdate()
  },
  methods: {
    determineFirstRunAfterUpdate(this: any) {
      // Why a separate thing? We want to 'collect' the value and not be affected by its future changes, because we're gonna set it to false immediately if found to be true.
      // console.log('determine first run after update runs.')
      // console.log('fraus start: ',this.$store.state.firstRunAfterUpdate)
      if (this.$store.state.firstRunAfterUpdate) {
        this.$store.dispatch('setFirstRunAfterUpdateState', false)
        // console.log('fraus end: ',this.$store.state.firstRunAfterUpdate)
        return true
      }
      return false
    },
  },
}
// These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
var headline = '# Changelog'
var intro = `**This is the change history between versions. The most current version of the changelog can be found [here](https://getaether.net/docs/changelog).**`
var content = `
## 2.0.0.dev.12
*Apr 24, 2019*

This update fixes several UI bugs, adds hashtag support and provides some quality-of-life improvements. It also updates some dependencies and lays the groundwork for future functionality such as @username detection and mentions.

### Features

- **Adds new markdown renderer.** This is what allows us to start building features like detection of certain pieces of text in user posted content.

- **Hashtag detection.** Now, when you post a #hashtag, such as #test, it becomes a search query that provides a live feed of that hashtag's use.

- **Aether:// link detection in body text.** When you link to another part of the app within the app, the links are now detected, just like regular http:// links. This makes linking discussions easier.

- **Syntax highlighting in 185 programming languages.** This adds syntax definitions for most major programming languages in the world to make Aether more welcoming for programming discussions.

- **Better Fingerprint / linking UI.** Link copy UI now provides feedback that a link is copied, and provides an explainer on what's happening. In the case the user tries to copy the link manually with copy and paste, it now makes sure that the full link is copied without missing digits.

- **BUG FIX** Aether links on the thread headers are now properly linked without needing a roundtrip to the browser. This was an issue that sometimes re-routed users to their browser when an aether:// link was clicked.


## 2.0.0.dev.11
*Feb 22, 2019*

This is the update that was built over the last few days after reaching 500 concurrent users in the network. This is an entirely backend-based, network-stability oriented update, with few user-facing changes — those will come later after getting these new backend improvements to be stable.

### General comments

This update should generally improve connectivity and speed of delivery. It's the first in a series of updates aimed at the same goal.

***You might see a drop or slowdown in your inbound or outbound connection counts. This doesn't mean the connection frequency has slowed.***

This is because the app now *only* counts and shows you connections that it can ***confirm*** as successful. There are connections *(e.g. inbound syncs happening with versions of the app older than this one)* that succeed and deliver full data, but won't ever be confirmed, because sending confirmations is a feature that is landing with this version.

In other words, it is possible to have your inbounds at zero, marked yellow, but still correctly functioning, because all your peers might be older versions of Aether.

As more of the network moves to this version and beyond, this should become less common.

### Backend

- Making network scans a process that runs every 10 minutes, and not something that can be called in from the scout (outreach) process.

- Streamline scout to the point that it will repeatedly try to connect to other nodes, with the knowledge that scout sync attempts no longer trigger network scans directly.

- Make reverse opens formalised as the reverse scout, in the same form of *"try again until succeeds"* logic as the normal scout, but with a lower number of maximum attempts.

- Create a third type of lease called *"ping"*, specifically for network scan requests marked as such.

- Create a new /ping endpoint bucket that other nodes can call to ping you, and this no longer occupies slots in the sync leases.

- Providing some data about the state of the reverse connections back to the requesting node, so the requester can know whether it succeeded or failed, and can act accordingly.

- Have a reverse dispatcher exclusions so that reverse opens will operate on their own exclusions list.

- Handles a case in the client where we might have to raise a notification about a response whose user key has not arrived yet.

- Fixes a bug where badlist refresh cycle was accidentally using the same name as the cache generation cycle.

- We now only show the user the successful syncs in the status page of the user interface. That way, failed syncs won't increment the number and give the wrong impression.

## 2.0.0.dev.10
*Feb 19, 2019*

This release is the first release that comes with Windows code signing with Authenticode. This should help with antivirus errors if you've been experiencing any. This also reverts a change that landed in dev.9 that was meant to improve file path compatibility with Windows, however, it appears that it has caused connectivity issues with some users. This fix is reverted for the time being out of abundance of caution.

### Improvements

- Windows binaries are now code-signed. This should make the app easier to install on Windows, eventually.

- Potential connectivity issue in some computers that came with a change in dev.9 is reverted.

- Search is now moved to the sidebar from the header. It should be more visible now.

- A bug in the user info pane where the unique username hover status of the username would end up stale in showing its canonicality status is now fixed.

- Misc. design fixes for consistency.

## 2.0.0.dev.9
*Feb 13, 2019*

This is a minor update that fixes some design bugs coming from the search feature of the prior release.

## 2.0.0.dev.8
*Feb 11, 2019*

This release comes with a major feature, search, and a few minor improvements and bug fixes.

### Features

#### Search

Search allows you to look for communities, content and people that you are looking for. This was a major user request, so I'm very happy to announce its availability.

- For communities, it searches in the community name and description
- For content, it searches in the titles and body (for threads and posts)
- For people, it searches the username and info.

- For all of these, it also searches for the fingerprint of each

So if someone gives you a fingerprint of a community, piece of content or user, you can find it that way, as well.

(This is the first release of this feature, bugs are to be expected. Please report bugs if you find any at b/meta, or at meta.getaether.net. Thank you!)

**Heads up:** This will mangle your notification order for one time — since it has to regenerate the frontend database to allow for it to be indexed for the first time. This is a one-time thing. It will also add a search index file, which is fairly sizeable - 200-300mb is normal. If the size becomes a problem in the future, we'll add a 'disable search' switch.)

#### Improvements

- Adds a spell checker. Now the misspelled text will underline itself in red. This also comes with an automatic language detector, so it should work in most languages. Mind that this is very minimalistic, it does not offer to correct the spelling for you. If that becomes an issue, we'll add it later.

- Made the lightbox mode for images and videos a little darker. Also, now the lightbox can be closed by pressing esc as well.

- Improved the webpack configuration of the app, so that the bundle size is smaller, and the app tries to load fewer CSS files per page. This should make the app feel generally snappier.

- Clicking the name 'aether' on all versions now makes the app go to 'Home' view. This is to aid navigation while the sidebar is collapsed, especially on narrower screen sizes.

- The versions for Electron and other software dependencies are updated.

#### Bug fixes

- Improves the fix to a bug where in Windows, in the case where the app is prevented from writing a cache index, it would generate a large amount of caches. This bug fix does two things, in one, it closes a few more of the possible places where this write denial can happen, and in the other, it prevents multiple caches from being generated in the case the write denial happens. Aether will receive an application signature from Microsoft shortly, whenever that happens, Windows (and antiviruses) should start to trust Aether enough to not prevent it working — that will be the ultimate fix.

- Fixes a bug where the text composer would accept text fields with too short entries, even after warning the user that the text is too short. Now it correctly declines these fields.

## 2.0.0.dev.7
*Dec 27, 2018*

Happy holidays everyone!

In this last release of the year, there's a bunch of goodies, feature requests from the community, and a bunch of bugfixes.

### Features

#### Keyboard Shortcuts

There are quite a few new keyboard shortcuts implemented, and they're now documented in 'Shortcuts' in Preferences.

#### Control over images and videos loading

You can now disable or enable the images and videos being auto loaded into the app. The app always had a whitelist of services to load from (Imgur, Giphy, Gfycat), so it never loaded everything. Now, you can completely disable this image and video autoloading feature to get more privacy, or simply, just a little more of a professional look.

### Bug fixes & miscellania

- In the Windows version, there was a bug that sometimes created a large number of cache files due to the app not being able to properly generate cache indexes, which made it think that it had not generated caches before. This is now fixed.

- The versions of Electron, Vue, and a whole host of dependencies are now updated to the most recent versions. This should make everything slightly faster.

- When two posts are the same in rankings, now the more recent post will surface to the top.

- If the UI is too small to contain some UI text (i.e. Preferences becomes Pref...), the full text of the UI item is now visible as hover.

- Design fixes so that the code blocks will look more attractive, and they won't break the composer. This makes is much easier to post source code intact.

- The app now retains the window size from the last run.

- Some wonkiness in the size of the status cards is now fixed.

- A bug where clicking the notification on Mac did not result in the app surfacing is now fixed.

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
@import '../../../scss/globals';
.settings-sublocation {
  color: $a-grey-600;
  .markdowned {
    &:first-of-type {
      margin-bottom: 0;
    }
    margin-bottom: 40px;
  }
  .changelog-successful-update {
    background-color: rgba(0, 0, 0, 0.4);
    padding: 15px 20px;
    margin-bottom: 15px;
    border-radius: 3px;
    .changelog-header {
      font-family: 'SCP Bold';
      margin-bottom: 15px;
    }
    .changelog-text {
      font-family: 'SCP Regular';
    }
  }
  .patreon-segment {
    margin-bottom: 20px;
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
      display: flex;
      .crypto-fund-button {
        margin-left: 10px;
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
  & > h2 + p {
    // All direct descendants of above that are immediately after a h2. i.e. this is the date line.
    margin-top: 0;
    margin-bottom: 0;
  }
}
</style>
