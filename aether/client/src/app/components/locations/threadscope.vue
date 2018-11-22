<template>
  <div class="location" v-if="$store.state.currentThreadLoadComplete">
    <template v-if="entityNotFound">
      <a-notfound></a-notfound>
    </template>
    <template v-else>
      <div class="threadscope">
        <a-thread-header-entity :thread="currentThread"></a-thread-header-entity>
        <a-no-content no-content-text="There are no posts in this thread yet. You should write something." v-if="posts.length === 0 && inflightChildren.length === 0"></a-no-content>
        <!-- <div class="composer-box">
        <a-composer :spec="composerSpecOne"></a-composer>
      </div> -->
        <div class="composer-box" v-if="!localUserReadOnly">
          <a-composer :spec="postComposer"></a-composer>
        </div>
        <a-post v-for="iflChild in inflightChildren.slice().reverse()" :post="iflChild.entity" :key="iflChild.status.requestedtimestamp+iflChild.entity.body" :inflightStatus="iflChild.status"></a-post>
        <a-post v-for="p in posts" :post="p" :notificationparent="route_parentSelector" :notificationhighlights="route_highlightSelectors" :notificationfocus="route_focusSelector"></a-post>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
  var globalMethods = require('../../services/globals/methods')
  var Vuex = require('../../../../node_modules/vuex').default
  var GetPlaceholder = require('../../services/phpicker/phpicker')
  var mimobjs = require('../../../../../protos/mimapi/structprotos_pb.js')
  var fe = require('../../services/feapiconsumer/feapiconsumer')
  var mixins = require('../../mixins/mixins')
  export default {
    name: 'threadscope',
    mixins: [mixins.localUserMixin],
    // props: [, 'route_focusSelector', 'route_parentSelector', 'route_highlightSelectors'],
    props: {
      route_focusSelector: {
        type: String,
        default: ""
      },
      route_parentSelector: {
        type: String,
        default: ""
      },
      route_highlightSelectors: {
        type: Array,
        default: function() { return [] }
      },
    },
    data(this: any): any {
      return {
        postComposer: {
          fields: [{
            id: "postBody",
            emptyWarningDisabled: true,
            visibleName: "",
            description: "",
            placeholder: "",
            maxCharCount: 20480,
            heightRows: 5,
            previewDisabled: false,
            content: "",
            optional: false,
          }],
          draft: {
            parentFp: this.$route.params.threadfp,
            contentType: "post"
          },
          commitAction: this.submitPost,
          commitActionName: "SUBMIT",
          cancelAction: function() {},
          // ^ It doesn't have a cancel action because it does not have a cancel button.
          cancelActionName: "",
        },
      }
    },
    computed: {
      ...Vuex.mapState(['currentThread']),
      inflightChildren(this: any) {
        let iflChildren = []
        for (let val of this.$store.state.ambientStatus.inflights.postsList) {
          if (!(this.$store.state.currentThread.fingerprint === val.entity.parent && val.status.eventtype == 'CREATE')) {
            continue
          }
          iflChildren.push(val)
        }
        return iflChildren
      },
      entityNotFound(this: any) {
        return this.$store.state.currentThread.fingerprint.length === 0
      },
      posts(this: any) {
        return this.currentThread.childrenList
      }
    },
    methods: {
      submitPost(this: any, fields: any) {
        let postBody = ""
        for (let val of fields) {
          if (val.id === 'postBody') {
            postBody = val.content
          }
        }
        let post = new mimobjs.Post
        // Set board, thread, parent, body fields
        post.setBoard(this.$store.state.currentThread.board)
        post.setThread(this.$store.state.currentThread.fingerprint)
        post.setParent(this.$store.state.currentThread.fingerprint)
        post.setBody(postBody)
        fe.SendPostContent('', post, function(resp: any) {
          console.log(resp.toObject())
        })
      },
      postVisible(this: any, post: any) {
        // If not compiled or inflight, visible
        if (post.uncompiled || post.isInflightEntity) {
          return true
        }
        // If modapproved, regardless of whether it's deleted, visible
        if (post.compiledcontentsignals.modapproved) {
          return true
        }
        if (post.compiledcontentsignals.modblocked) {
          return false
        }
        // If not approved and not deleted, visible
        // Default case, no if clause needed
        return true
      }
    },
    beforeMount(this: any) {
      if (globalMethods.IsUndefined(this.posts)) {
        return
      }
      if (this.posts.length === 0) {
        // Blank textbox without a placeholder if there's no content.
        return
      }
      this.postComposer.fields[0].placeholder = GetPlaceholder().Placeholder
    },
    mounted(this: any) {
      // if (this.$route.query.focusSelector.length === 0) {
      //   return
      // }
      // let selector = this.$store.state.focusSelector
      // let selector = this.$route.query.focusSelector
      // // this.$store.dispatch('setFocusSelector', "")
      // let iterCount = 0
      // let checker = function() {
      //   console.log('checker runs')
      //   let el: any = document.getElementById(selector)
      //   if (el !== null) {
      //     el.scrollIntoView()
      //   } else {
      //     if (iterCount < 100) {
      //       setTimeout(checker, 50)
      //       iterCount++
      //     }
      //   }
      // }
      // setTimeout(function() {
      //   checker()
      // }, 25)
    },
    beforeUpdate(this: any) {
      if (this.posts.length === 0) {
        this.postComposer.fields[0].placeholder = ""
        return
      }
      if (this.postComposer.fields[0].placeholder.length === 0) {
        this.postComposer.fields[0].placeholder = GetPlaceholder().Placeholder
      }
    }
  }
</script>
<style lang="scss">
  .location {
    &.notfound {
      height: 100%;
    }
  }
</style>
<style lang="scss" scoped>
  .composer-box {
    padding: 20px;
    display: flex;
  }
</style>