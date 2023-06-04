<template>
  <div class="post-listitem" @click="goToPost">
    <div
      class="post"
      :class="{
        'self-owned': post.selfcreated,
      }"
      :id="'post-' + post.fingerprint"
    >
      <div class="meta-container">
        <div class="author">
          <a-username :isop="byOp" :owner="postOwner"></a-username>
        </div>
        <div class="post-datetime">
          <a-timestamp
            :creation="postCreation"
            :lastupdate="postLastUpdate"
          ></a-timestamp>
        </div>
      </div>
      <div class="post-content">
        <div class="content-text-container">
          <a-markdown
            class="content-text"
            :content="this.postBodyPreview"
          ></a-markdown>
        </div>
        <a-ballot
          v-if="actionsVisible"
          :contentsignals="post.compiledcontentsignals"
          :boardfp="post.board"
          :threadfp="post.thread"
        ></a-ballot>
      </div>
      <div class="post-crumb">
        <a-boardname
          :name="post.viewmetaBoardname"
          :fingerprint="post.board"
        ></a-boardname>
        <icon name="chevron-right" scale="0.85" class="post-crumb-caret"></icon>
        <a-threadname
          :name="post.viewmetaThreadname"
          :board-fingerprint="post.board"
          :thread-fingerprint="post.thread"
        ></a-threadname>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var globalMethods = require('../../services/globals/methods')
var mixins = require('../../mixins/mixins')
export default {
  name: 'a-post-listitem',
  mixins: [mixins.localUserMixin],
  // props: ['post', 'inflightStatus', 'uncompiled', 'notificationparent', 'notificationhighlights', 'notificationfocus'],
  props: {
    post: {
      type: Object,
      default: function () {
        return undefined
      },
    },
  },
  data(this: any): any {
    return {}
  },
  computed: {
    postCreation(this: any) {
      return this.post.creation
    },
    postLastUpdate(this: any) {
      return this.post.lastupdate
    },
    postOwner(this: any) {
      return this.post.owner
    },
    postFingerprint(this: any) {
      return this.post.fingerprint
    },
    postBodyPreview(this: any) {
      let trailer = this.post.body.length > 750 ? '...' : ''
      return this.post.body.substring(0, 750) + trailer
    },
    /*----------  By op  ----------*/
    byOp(this: any) {
      return this.post.compiledcontentsignals.byop
    },
    linkToPost(this: any) {
      return (
        '/board/' +
        this.post.board +
        '/thread/' +
        this.post.thread +
        '?focusSelector=' +
        this.post.fingerprint +
        '&highlightSelectors=["' +
        this.post.fingerprint +
        '"]'
      )
    },
    /*----------  Visibility checks  ----------*/
    /*
        These are useful because our visible / nonvisible logic is getting too complex to retain in the template itself.
      */
    actionsVisible(this: any) {
      if (this.localUserReadOnly) {
        return false
      }
      return true
    },
  },
  mounted(this: any) {
    // console.log("post self created")
    // console.log(this.post.selfcreated)
    if (this.isNotificationsFocus) {
      // this.$el.scrollIntoView()
      this.$el.scrollIntoView({
        behavior: 'smooth',
        block: 'start',
      })
    }
  },
  methods: {
    getOwnerName(owner: any): string {
      return globalMethods.GetOwnerName(owner)
    },
    /*----------  Collapse / open state  ----------*/
    toggleCollapse(this: any) {
      if (this.collapsed) {
        this.collapsed = false
        return
      }
      this.collapsed = true
    },
    /*----------  Copy link  ----------*/
    copyLink(this: any) {
      let elc = require('electron')
      elc.clipboard.writeText(this.externalUrl)
    },
    goToPost(this: any) {
      this.$router.push(this.linkToPost)
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../../scss/globals';
.post-listitem {
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  margin: 0px 20px;
  &:hover {
    background-color: rgba(255, 255, 255, 0.05);
  }
  display: flex;
  cursor: pointer;
}
.post {
  display: flex;
  flex-direction: column;
  padding-left: 20px;
  margin: 20px 20px;
  border-left: 3px solid rgba(255, 255, 255, 0.15);
  position: relative;
  flex: 1;

  & > .post {
    margin-left: 0;
    margin-right: 0;
    padding-left: 15px;
  }
  .meta-container {
    display: flex;
    .post-datetime {
      margin-left: 10px;
      font-family: 'SSP Regular Italic';
      color: $a-grey-600;
    }
  }
  .post-content {
    font-family: 'SSP Regular';
    display: flex;
    font-size: 110%;

    .content-text-container {
      flex: 1;

      .composer {
        margin-bottom: 15px;
      }
    }

    .content-text {
      flex: 1;
      p:last-of-type {
        margin-bottom: 5px;
      }
    }
  }
  &:hover {
    > .post-content > .ballot {
      visibility: visible;
    }
  }
}

.expand-container {
  position: absolute;
  top: -8px;
  left: -14px;
  background-color: $mid-base;
  border-radius: 10px;
  display: flex;
  cursor: pointer;
  height: 35px;
  width: 25px;
  &:hover {
    svg {
      fill: $a-grey-600;
    }
  }
  svg {
    fill: rgba(90, 90, 90, 1);
    width: 12px;
    height: 12px;
    margin: 0 auto;
    margin-top: 15px;
  }
}

.self-owned {
  // background-color: red;
  border-color: $a-cerulean;
}

.notification-highlight {
  border-color: $a-orange; // background-color: blue;
}

.post-crumb-caret {
  opacity: 0.8;
  position: relative;
  top: 2px;
  margin: 0 4px 0 5px;
}
</style>
