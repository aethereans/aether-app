<template>
  <div class="board-sublocation">
    <div class="newthread">
      <div class="thread-composer">
        <a-composer :spec="threadComposerSpec"></a-composer>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var fe = require('../../../services/feapiconsumer/feapiconsumer')
var mimobjs = require('../../../../../../protos/mimapi/structprotos_pb.js')
export default {
  name: 'newthread',
  data(this: any): any {
    return {
      threadComposerSpec: {
        fixToBottom: true,
        fields: [
          {
            id: 'threadName',
            visibleName: 'Title',
            description:
              'Keep it short, consider using the original title if present.',
            placeholder: '',
            maxCharCount: 256,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: false,
          },
          {
            id: 'threadLink',
            visibleName: 'Link',
            description:
              "Optional. If it's an image, GIF or video, <b><em>and</em></b> you want it to have a preview, upload to a whitelisted host (Imgur, Gfycat, Giphy) and link the file directly.",
            placeholder: 'https://imgur.com/random.jpg',
            maxCharCount: 2048,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: true,
          },
          {
            id: 'threadBody',
            visibleName: 'Body',
            description:
              'Optional. Markdown is available. Click the <em>eye</em> to see a preview.',
            placeholder: '',
            maxCharCount: 20480,
            heightRows: 10,
            previewDisabled: false,
            content: '',
            optional: true,
          },
        ],
        draft: {
          parentFp: this.$route.params.boardfp,
          contentType: 'thread',
        },
        commitAction: this.submitThread,
        commitActionName: 'SUBMIT',
        cancelAction: function() {
          history.back()
        },
        cancelActionName: 'CANCEL',
        autofocus: true,
      },
    }
  },
  methods: {
    submitThread(this: any, fields: any) {
      let threadName = ''
      let threadLink = ''
      let threadBody = ''
      for (let val of fields) {
        if (val.id === 'threadName') {
          threadName = val.content
          continue
        }
        if (val.id === 'threadLink') {
          threadLink = val.content
          continue
        }
        if (val.id === 'threadBody') {
          threadBody = val.content
          continue
        }
      }
      let thread = new mimobjs.Thread()
      // Set board, thread, parent, body fields
      thread.setBoard(this.$store.state.currentBoard.fingerprint)
      thread.setName(threadName)
      thread.setLink(threadLink)
      thread.setBody(threadBody)
      let vm = this
      fe.SendThreadContent('', thread, function(resp: any) {
        console.log(resp.toObject())
        vm.$router.push(
          '/board/' + vm.$store.state.currentBoard.fingerprint + '/new'
        )
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.thread-composer {
  padding: 0 50px;
  padding-top: 25px;
}
</style>

<style lang="scss">
.location.board {
  padding-bottom: 20px;
}
</style>
