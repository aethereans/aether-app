<template>
  <div class="global-sublocation">
    <div class="newboard">
      <div class="board-composer">
        <a-composer :spec="boardComposerSpec"></a-composer>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var fe = require('../../../services/feapiconsumer/feapiconsumer')
var mimobjs = require('../../../../../../protos/mimapi/structprotos_pb.js')
export default {
  name: 'newboard',
  data(this: any): any {
    return {
      boardComposerSpec: {
        fixToBottom: true,
        fields: [
          {
            id: 'boardName',
            visibleName: 'Community name',
            description:
              'Short, sweet, without spaces. It will be referred to as <i>b/yourcommunity</i>. This cannot be changed later.',
            placeholder: 'fightclub',
            minCharCount: 2,
            maxCharCount: 32,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: false,
            spaceDisabled: true,
          },
          {
            id: 'boardDescription',
            visibleName: 'Description',
            description:
              'Optional. This is where you put the community rules you want to tell people about, if any. Markdown is available. Click the <em>eye</em> to see a preview.',
            placeholder: "Don't talk about the fight club.",
            maxCharCount: 20480,
            heightRows: 12,
            previewDisabled: false,
            content: '',
            optional: true,
          },
        ],
        draft: {
          parentFp: 'ROOT', // communities have no parents.
          contentType: 'board',
        },
        commitAction: this.submitBoard,
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
    submitBoard(this: any, fields: any) {
      let boardName = ''
      let boardDescription = ''
      for (let val of fields) {
        if (val.id === 'boardName') {
          boardName = val.content
          continue
        }
        if (val.id === 'boardDescription') {
          boardDescription = val.content
          continue
        }
      }
      let board = new mimobjs.Board()
      // Set board, thread, parent, body fields
      board.setName(boardName)
      board.setDescription(boardDescription)
      this.$router.push(
        '/user/' + this.$store.state.localUser.fingerprint + '/boards'
      )
      fe.SendBoardContent('', board, function(resp: any) {
        console.log(resp.toObject())
      })
      fe.SendModModeEnabledStatus(true, function() {
        console.log('Mod mode is enabled due to creation of a board.')
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.board-composer {
  padding: 0 50px;
  padding-top: 25px;
}
</style>

<style lang="scss">
.location.global {
  padding-bottom: 20px;
}
</style>
