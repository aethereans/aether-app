<template>
  <router-link
    class="thread-entity"
    :to="linkToThread"
    @click.self
    :class="{ inflight: isInflightEntity }"
  >
    <div class="signals-container">
      <div class="thread-vote-count">
        <icon name="arrow-up"></icon> {{ voteCount }}
      </div>
      <div class="thread-comment-count" v-show="commentCount !== -1">
        <icon name="comment-alt"></icon> {{ commentCount }}
      </div>
    </div>
    <div
      class="image-container"
      v-show="imageLoadedSuccessfully || videoLoadedSuccessfully"
    >
      <div
        class="image-box"
        v-show="imageLoadedSuccessfully && !videoLoadedSuccessfully"
        :style="'background-image: url(' + sanitisedLink + ')'"
        @click.prevent="openLightbox"
      ></div>
      <video
        class="video-box"
        :src="sanitisedLink"
        preload="metadata"
        @click.prevent="openLightbox"
        v-show="videoLoadedSuccessfully && !imageLoadedSuccessfully"
        @loadeddata="handleVideoLoadSuccess"
      ></video>
      <div
        class="lightbox"
        @click.prevent="closeLightbox"
        v-show="lightboxOpen"
      >
        <img
          :src="sanitisedLink"
          alt=""
          @load="handleImageLoadSuccess"
          v-show="imageLoadedSuccessfully"
        />
        <video
          class="video-player"
          v-if="videoLoadedSuccessfully"
          preload="metadata"
          :src="sanitisedLink"
          muted="true"
          controls
          v-show="lightboxOpen"
        ></video>
      </div>
    </div>
    <div class="main-data-container">
      <div class="thread-main-data-centerer">
        <div class="inflight-box" v-if="isInflightEntity">
          <a-inflight-info
            :status="inflightStatus"
            :refresherFunc="refresh"
          ></a-inflight-info>
        </div>

        <div class="thread-header">
          <div class="thread-name">
            {{ thread.name }}
          </div>
        </div>
        <div class="thread-link" v-show="sanitisedLink.length > 0">
          <a-link :link="thread.link" :linktext="thread.link"></a-link>
        </div>
        <div class="thread-body">
          {{ threadBodyPreview }}
        </div>
        <div class="thread-metadata">
          <div class="thread-metadata-container">
            <!-- &nbsp;—&nbsp; -->
            <a-boardname
              :name="thread.viewmetaBoardname"
              :fingerprint="thread.board"
            ></a-boardname>
            <a-timestamp
              class="thread-entity-timestamp"
              :creation="thread.creation"
              :lastupdate="thread.lastupdate"
            ></a-timestamp>
          </div>
        </div>
      </div>
    </div>
    <a-ballot
      v-if="actionsVisible"
      :contentsignals="thread.compiledcontentsignals"
      :boardfp="thread.board"
      :threadfp="thread.fingerprint"
    ></a-ballot>
  </router-link>
</template>

<script lang="ts">
var vuexStore = require('../../store/index').default
var mixins = require('../../mixins/mixins')
var Mousetrap = require('../../../../node_modules/mousetrap')
var globalMethods = require('../../services/globals/methods')
export default {
  name: 'a-thread-listitem',
  mixins: [mixins.localUserMixin],
  props: ['thread', 'inflightStatus'],
  data(this: any) {
    return {
      imageLoadedSuccessfully: false,
      videoLoadedSuccessfully: false,
      lightboxOpen: false,
    }
  },
  computed: {
    // ...Vuex.mapState(['currentBoard']),
    isInflightEntity(this: any) {
      if (typeof this.inflightStatus !== 'undefined') {
        return true
      }
      return false
    },
    voteCount(this: any) {
      if (typeof this.thread.compiledcontentsignals === 'undefined') {
        return 0
      }
      return (
        this.thread.compiledcontentsignals.upvotes -
        this.thread.compiledcontentsignals.downvotes
      )
    },
    commentCount(this: any) {
      if (typeof this.thread.compiledcontentsignals === 'undefined') {
        return 0
      }
      return this.thread.postscount
    },
    sanitisedLink(this: any) {
      if (this.thread.link.length === 0) {
        return ''
      }
      if (typeof this.thread.link === 'undefined') {
        return ''
      }
      if (
        this.thread.link.substring(0, 8) === 'https://' ||
        this.thread.link.substring(0, 7) === 'http://'
      ) {
        return this.thread.link
      }
      return 'http://' + this.thread.link
    },
    linkToThread(this: any) {
      if (this.isInflightEntity) {
        return ''
      }
      return (
        '/board/' + this.thread.board + '/thread/' + this.thread.fingerprint
      )
    },
    /*----------  Visibility  ----------*/
    actionsVisible(this: any) {
      if (this.isInflightEntity) {
        return false
      }
      if (this.localUserReadOnly) {
        return false
      }
      return true
    },
    threadBodyPreview(this: any) {
      let trailer = this.thread.body.length > 250 ? '...' : ''
      return this.thread.body.substring(0, 250) + trailer
    },
  },
  methods: {
    openLightbox(this: any) {
      let vm = this
      Mousetrap.bind('esc', function() {
        vm.closeLightbox()
        console.log('mousetrap bind works')
      })
      let vid: any = this.$el.getElementsByClassName('video-player')[0]
      this.lightboxOpen = true
      if (!globalMethods.IsUndefined(vid)) {
        vid.play()
      }
    },
    closeLightbox(this: any) {
      Mousetrap.unbind('esc')
      let vid: any = this.$el.getElementsByClassName('video-player')[0]
      this.lightboxOpen = false
      if (!globalMethods.IsUndefined(vid)) {
        vid.pause()
      }
    },
    handleImageLoadSuccess(this: any) {
      // console.log('image loaded successfully')
      this.imageLoadedSuccessfully = true
    },
    handleVideoLoadSuccess(this: any) {
      // console.log('video loaded successfully')
      this.videoLoadedSuccessfully = true
    },
    refresh(this: any) {
      let sortByNew = false
      if (this.$store.state.route.name === 'Board>ThreadsNewList') {
        sortByNew = true
      }
      vuexStore.dispatch('refreshCurrentBoardAndThreads', {
        boardfp: this.thread.board,
        sortByNew: sortByNew,
      })
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../../scss/globals';
.thread-entity {
  display: block;
  @extend %link-hover-ghost-extenders-disable;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  padding: 15px 5px;
  margin: 0 20px; // cursor: default;
  color: $a-grey-800;

  &.inflight {
    cursor: default;
    .main-data-container .thread-name {
      cursor: default;
    }
  }

  .lightbox {
    display: flex;
    outline: none;
    position: fixed;
    z-index: 999;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    background-color: rgba(0, 0, 0, 0.85);

    img,
    video {
      margin: auto;
      max-height: 90%;
      max-width: 90%;
    }
  }

  &:hover {
    background-color: rgba(255, 255, 255, 0.05);
    .actions-container .upvote {
      visibility: visible;
      &.voted::before {
        background-color: $a-cerulean-100;
      }
      svg {
        // fill: $a-grey-800;
      }
    }
    .ballot {
      visibility: visible;
    }
  }
  .actions-container .upvote.voted {
    visibility: visible;
    &::before {
      box-shadow: 0 0 0 2px $a-transparent inset;
      background-color: $a-cerulean-80;
    }
    svg {
      fill: $mid-base;
    }
  }
  .signals-container {
    width: 64px;
    font-size: 110%;
    margin: auto;
    display: flex;
    flex-direction: column;
    color: $a-grey-600;

    .thread-vote-count {
      flex: 1;
      margin: auto;
      font-family: 'SSP Black';
      display: flex;
      svg {
        margin: auto;
        margin-right: 3px;
      }
    }

    .thread-comment-count {
      @extend .thread-vote-count;
      svg {
        margin-right: 5px;
        height: 13px;
        width: 13px;
      }
    }
  }

  .preview-image {
    width: 0px;
    height: 0px; // Because we just want to use this as a way to get the load failed event in the case it fails - background-image load does not raise failure events.
  }

  .image-container {
    width: 80px;
    padding: 6px 8px;
    padding-left: 0;

    .image-box {
      // width: 80px;
      height: 100%;
      overflow: hidden;
      border-radius: 2px;
      background-size: cover;
      background-color: $a-cerulean;
      background-size: cover;
      background-position: center center;
      cursor: pointer;

      .thread-image {
        object-fit: cover;
        height: inherit;
      }
    }
    .video-box {
      max-height: 150px;
      max-width: 100%;
      overflow: hidden;
      border-radius: 2px;
      overflow: hidden;
      cursor: pointer;
    }
  }

  .main-data-container {
    flex: 1;
    padding-right: 15px;
    display: flex;
    min-width: 0;

    .thread-main-data-centerer {
      margin: auto 0;
      flex: 1;
      min-width: 0;
    }
    .thread-metadata {
      display: inline-block;
      .thread-metadata-container {
        display: flex;
        color: $a-grey-800;
        margin-top: 8px;
        // font-family: "SSP Bold Italic";
        .thread-entity-timestamp {
          display: inline-block;
        }
        .a-boardname {
          margin-right: 10px;
        }
      }
    }
    .thread-header {
      .thread-name {
        font-size: 120%;
        cursor: pointer;
        display: inline;
      }
    }

    .thread-link {
      font-size: 95%;
      overflow: hidden;
      text-overflow: ellipsis;
      position: relative;
      left: -5px;
      padding-left: 5px;
      color: $a-cerulean;
      a {
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
        display: inline;
        word-break: break-word;
      } // color: $a-cerulean;
      // @include lineClamp(1);
    }

    .thread-body {
      font-family: 'SSP Regular';
      color: $a-grey-500;
      font-size: 95%;
      @include lineClamp(3);
      word-break: break-word;
    }
  }
  .actions-container {
    width: 75px;
    display: flex;
    margin-left: 10px;
    .upvote {
      visibility: hidden;
      width: 50px;
      height: 50px;
      margin: auto;
      margin-left: 0;
      display: flex;
      cursor: pointer;
      position: relative;

      &::before {
        // nonvoted normal state
        content: '';
        width: 54px;
        height: 54px;
        left: -2px;
        top: -2px;
        display: inline-block;
        border-radius: 27px;
        position: absolute;
        opacity: 1;
        box-shadow: 0 0 0 2px $a-grey-500 inset;
      }

      svg {
        padding: 10px;
        margin: auto;
        width: 100%;
        height: 100%;
        fill: $a-grey-500;
        position: relative; // z-index: 2;
      }

      &:hover {
        // non-voted hover state
        &::before {
          box-shadow: 0 0 0 2px $a-grey-800 inset; // background-color: $a-cerulean;
        }

        svg {
          fill: $a-grey-800;
        }
      } // &:active::before { // background-color: $a-cerulean-100; // } // &.voted {
      //   &::before {
      //     // voted normal state
      //     // background-color: $a-cerulean-100;
      //     box-shadow: 0 0 0 2px $a-grey-800 inset;
      //   }
      //   svg {
      //     fill: $a-grey-800;
      //   }
      // }
    }
  }
}
</style>
