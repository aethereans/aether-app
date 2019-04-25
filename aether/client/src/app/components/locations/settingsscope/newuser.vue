<template>
  <div class="settings-sublocation create-user">
    <template v-if="localUserExists && !mintingStarted">
      <!-- ^ Minting started conditional because we don't want to show this the first time the user actually completes the process.  -->
      <!-- Local user already exists -->
      <a-markdown :content="userAlreadyExistsHeadline"></a-markdown>
      <a-markdown :content="userAlreadyExistsIntro"></a-markdown>
      <a-markdown :content="userAlreadyExistsContent"></a-markdown>
      <router-link class="button is-warning is-outlined" to="/popular"
        >GO TO POPULAR</router-link
      >
    </template>
    <template v-else>
      <template v-if="initialFormVisible">
        <a-markdown :content="initialHeadline"></a-markdown>
        <a-markdown :content="initialIntro"></a-markdown>
        <a-composer id="userComposer" :spec="createNewUserSpec"></a-composer>
      </template>
      <template v-if="inProgressVisible">
        <a-markdown :content="mintingInProgressHeadline"></a-markdown>
        <a-markdown :content="mintingInProgressContent"></a-markdown>
        <a-spinner :hidetext="true"></a-spinner>
      </template>
      <template v-if="completionVisible">
        <a-markdown :content="completionHeadline"></a-markdown>
        <a-markdown :content="completionContent"></a-markdown>
        <router-link to="/popular" class="button is-warning is-outlined">
          GO TO POPULAR
        </router-link>
      </template>
    </template>
  </div>
</template>

<script lang="ts">
// var globalMethods = require('../../../services/globals/methods')
var mixins = require('../../../mixins/mixins')
var fe = require('../../../services/feapiconsumer/feapiconsumer')
var mimobjs = require('../../../../../../protos/mimapi/structprotos_pb.js')
export default {
  name: 'newuser',
  mixins: [mixins.localUserMixin],
  data(this: any): any {
    return {
      /*----------  Initial state  ----------*/
      initialHeadline: initialHeadline,
      initialIntro: initialIntro,
      /*----------  minting state  ----------*/
      mintingInProgressHeadline: mintingInProgressHeadline,
      mintingInProgressContent: mintingInProgressContent,
      /*----------  completion state  ----------*/
      completionHeadline: completionHeadline,
      completionContent: completionContent,
      /*----------  user already exists state  ----------*/
      userAlreadyExistsHeadline: userAlreadyExistsHeadline,
      userAlreadyExistsIntro: userAlreadyExistsIntro,
      userAlreadyExistsContent: userAlreadyExistsContent,
      mintingStarted: false,
      createNewUserSpec: {
        fields: [
          {
            id: 'userName',
            visibleName: 'Pick a name',
            description: `These names are <p class="em">not unique</p>, there can be multiple users with the same name. However, the block avatars of two different users won't ever be the same. When in doubt, check name <i>and</i> the picture. <div id="postscript">By the way, if you become a supporter, you can get an exclusive, unique username in recognition of your support. Check <a href="#/membership">Membership</a> for more information.</div>`,
            placeholder: 'deanmoriarty',
            maxCharCount: 24,
            minCharCount: 3,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: false,
            spaceDisabled: true,
            preInfoEnabled: true,
            preInfoText: '@',
          },
          {
            id: 'userInfo',
            visibleName: 'Info',
            description:
              'Optional, can be changed later. Markdown is available.',
            placeholder: 'rebel without a cause / new york - san francisco',
            maxCharCount: 20480,
            heightRows: 5,
            previewDisabled: false,
            content: '',
            optional: true,
          },
        ],
        draft: {
          parentFp: 'ROOT', // users have no parents.
          contentType: 'key',
        },
        commitActionName: 'CREATE',
        commitAction: this.submitNewUser,
        cancelActionName: '',
        cancelAction: function() {},
        fixToBottom: true,
        autofocus: true,
      },
    }
  },
  computed: {
    initialFormVisible(this: any) {
      if (this.mintingStarted) {
        return false
      }
      return true
    },
    inProgressVisible(this: any) {
      if (this.$store.state.localUserExists) {
        return false
      }
      if (!this.mintingStarted) {
        return false
      }
      return true
    },
    completionVisible(this: any) {
      if (this.$store.state.localUserExists) {
        return true
      }
      return false
      // if (globalMethods.IsUndefined(this.$store.state.localUser)) {
      //   return false
      // }
      // if (globalMethods.IsEmptyObject(this.$store.state.localUser)) {
      //   return false
      // }
      // return true
    },
  },
  methods: {
    submitNewUser(this: any, fields: any) {
      this.mintingStarted = true
      let userName = ''
      let userInfo = ''
      for (let val of fields) {
        if (val.id === 'userName') {
          userName = val.content
          continue
        }
        if (val.id === 'userInfo') {
          userInfo = val.content
          continue
        }
      }
      let user = new mimobjs.Key()
      user.setName(userName)
      user.setInfo(userInfo)
      // let vm = this
      fe.SendUserContent('', user, function(resp: any) {
        console.log('user create request sent in.')
        console.log(resp.toObject())
      })
    },
  },
}
/*<br><br>(PS. Supporters of the work on this app can get unique  names and flair in recognition of their support. Check the Membership tab to the left if you're interested in that.)*/
// These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
/*----------  Initial state  ----------*/
var initialHeadline = '# Create new user'
var initialIntro = `**Hey there! 👋 &nbsp; Let's get you set up.**`
// The main body of the initial state is a form, provided above.
/*----------  Minting in progress state  ----------*/
var mintingInProgressHeadline = '# Key creation in progress...'
var mintingInProgressContent = `
Minting the proof-of-work for your user key. This can take from a few seconds to a couple minute depending on your computer. If you want to get updated when it is done, stay on this screen.`
/*----------  Completion state  ----------*/
var completionHeadline = '# Successfully created'
var completionContent = `
  * **Your user is now ready.** You can now vote, post, create communities, and moderate.

  * Your user name is not unique, but your generated picture is.

  * You can get a unique username no one else can have by becoming a [supporting member](#/membership).

  * Consider backing up your user profile somewhere safe. You can do it by following the instructions at the [Preferences](#/settings).

  `
/*----------  User already exists state  ----------*/
var userAlreadyExistsHeadline = '# User already exists'
var userAlreadyExistsIntro = `**There's already a user present on this app. (You probably just created it.)**`
var userAlreadyExistsContent = `If this is not you, or if you want to create a new user, you can delete the user by deleting the user profile. The location of the user profile can be found in *Preferences*. All settings and data will be wiped alongside the profile.`
</script>

<style lang="scss">
/* <<--global, not scoped */

@import '../../../scss/globals';
#userComposer {
  font-family: 'SSP Bold';
  p.em {
    font-family: 'SSP Black';
    display: inline;
  }
  #postscript {
    font-family: 'SSP Regular Italic';
    padding-top: 10px;
    letter-spacing: 0.3px;
  }
  .visible-name {
    color: $a-grey-600;
  }
}
</style>

<style lang="scss" scoped>
@import '../../../scss/globals';
.settings-sublocation {
  color: $a-grey-600;

  &.create-user {
    // font-size: 16px;
  }
  .markdowned {
    &:first-of-type {
      margin-bottom: 0;
    }
    margin-bottom: 40px;
  }
}

.button {
  font-family: 'SSP Bold';
}
</style>
