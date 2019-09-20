<template>
  <div class="board-sublocation">
    <div class="board-info">
      <div class="user-data-field-container">
        <div class="user-data-field-header">
          Description
        </div>
        <div class="user-data-field">
          <template v-if="!editPaneOpen">
            <a-markdown :content="visibleEntity.description"></a-markdown>
            <div class="inflight-box" v-if="isVisibleInflightEntity">
              <a-inflight-info
                :status="visibleInflightStatus"
              ></a-inflight-info>
            </div>
            <div class="actions-rack" v-if="actionsVisible">
              <a
                class="action edit"
                v-if="$store.state.currentBoard.selfcreated"
                @click="toggleEditPane"
              >
                <!-- todo: make it only visible when self, otherwise the button is meaningless -->
                Edit
              </a>
            </div>
          </template>
          <template v-else>
            <a-composer :spec="boardEditExistingSpec"></a-composer>
          </template>
        </div>
      </div>
      <!-- <div class="user-data-field-container">
        <div class="user-data-field-header">
          Share
        </div>
        <div class="user-data-field">
          <a class="field-action copy-link" @click="copyLink">
            Copy link
          </a>
        </div>
      </div> -->
      <div class="user-data-field-container">
        <div class="user-data-field-header">
          URL
        </div>
        <div class="user-data-field fingerprint">
          <!-- <div class="fingerprint-text">
            {{$store.state.currentBoard.fingerprint}}
          </div> -->
          <a-fingerprint
            :fp="$store.state.currentBoard.fingerprint"
            :link="this.externalUrl"
          ></a-fingerprint>
        </div>
      </div>
      <div class="user-data-field-container">
        <div class="user-data-field-header">
          Created
        </div>
        <div class="user-data-field">
          {{ timeString($store.state.currentBoard.creation) }}
        </div>
      </div>
      <div class="user-data-field-container">
        <div class="user-data-field-header">
          Created by
        </div>
        <div class="user-data-field">
          <a-username :owner="$store.state.currentBoard.owner"></a-username>
        </div>
      </div>
      <div class="user-data-field-container">
        <div class="user-data-field-header">
          Last Updated
        </div>
        <div class="user-data-field">
          {{ timeString($store.state.currentBoard.lastupdate) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var fe = require('../../../services/feapiconsumer/feapiconsumer')
var mimobjs = require('../../../../../../protos/mimapi/mimapi_pb.js')
var mixins = require('../../../mixins/mixins')
// var globalMethods = require('../../../services/globals/methods')
export default {
  name: 'boardinfo',
  mixins: [mixins.localUserMixin],
  data(this: any): any {
    return {
      editPaneOpen: false,
      boardEditExistingSpec: {
        fields: [
          {
            id: 'boardDescription',
            visibleName: '',
            description: '',
            placeholder: '',
            maxCharCount: 20480,
            heightRows: 12,
            previewDisabled: false,
            content: '',
            optional: true,
          },
        ],
        commitAction: this.submitEditExistingBoard,
        commitActionName: 'SAVE',
        cancelAction: this.toggleEditPane,
        cancelActionName: 'CANCEL',
      },
    }
  },
  methods: {
    timeString(input: any) {
      // All int64 values are provided to JS as strings by gRPC.. which means our timestamps are strings.
      if (parseInt(input) === 0) {
        return 'Never'
      }
      let d = new Date(input * 1000)
      return d.toDateString() + ', ' + d.toLocaleTimeString()
    },
    /*----------  Edit board actions  ----------*/
    toggleEditPane(this: any) {
      if (this.editPaneOpen) {
        this.editPaneOpen = false
      } else {
        this.boardEditExistingSpec.fields[0].content = this.visibleEntity.description
        this.editPaneOpen = true
      }
    },
    submitEditExistingBoard(this: any, fields: any) {
      let boardDescription = ''
      for (let val of fields) {
        if (val.id === 'boardDescription') {
          boardDescription = val.content
        }
      }
      let board = new mimobjs.Board()
      // Set board, thread, parent, body fields
      board.setDescription(boardDescription) // < heads up, this is different
      let vm = this
      fe.SendBoardContent(
        this.$store.state.currentBoard.fingerprint,
        board,
        function(resp: any) {
          vm.toggleEditPane()
          console.log(resp.toObject())
        }
      )
    },
    /*----------  Copy link  ----------*/
    copyLink(this: any) {
      let elc = require('electron')
      elc.clipboard.writeText(this.externalUrl)
    },
  },
  computed: {
    /*----------  Inflight computeds  ----------*/
    inflightEdits(this: any) {
      let iflEdits = []
      for (let val of this.$store.state.ambientStatus.inflights.boardsList) {
        if (
          this.$store.state.currentBoard.fingerprint !==
          val.entity.provable.fingerprint
        ) {
          continue
        }
        if (val.status.eventtype !== 'UPDATE') {
          continue
        }
        iflEdits.push(val)
      }
      return iflEdits
    },
    mostRecentInflightEdit(this: any) {
      let mostRecentTs = 0
      let mostRecent = undefined
      for (let val of this.inflightEdits) {
        if (val.status.requestedtimestamp >= mostRecentTs) {
          mostRecentTs = val.status.requestedtimestamp
          mostRecent = val
        }
      }
      return mostRecent
    },
    /*----------  Visibility  ----------*/
    visibleEntity(this: any) {
      if (typeof this.mostRecentInflightEdit !== 'undefined') {
        return this.mostRecentInflightEdit.entity
      }
      return this.$store.state.currentBoard
    },
    visibleInflightStatus(this: any) {
      if (typeof this.mostRecentInflightEdit !== 'undefined') {
        return this.mostRecentInflightEdit.status
      }
      return this.inflightStatus
    },
    isVisibleInflightEntity(this: any) {
      if (typeof this.visibleInflightStatus !== 'undefined') {
        return true
      }
      return false
    },
    actionsVisible(this: any) {
      if (this.localUserReadOnly) {
        return false
      }
      return true
    },
    externalUrl(this: any) {
      return 'aether://board/' + this.$store.state.route.params.boardfp
      // ^ This is a little different than the others, because it's not visible directly in its root page, it's visible on /info. so route.path ends up /board/[boardfp]/info while we're looking for /board/boardfp.
    },
    // isSelf(this: any) {
    //   if (globalMethods.IsUndefined(this.$store.state.localUser)) {
    //     return false
    //   }
    //   if (globalMethods.IsUndefined(this.$store.state.currentBoard)) {
    //     return false
    //   }
    //   if (this.$store.state.localUser.fingerprint !== this.$store.state.currentBoard.owner.fingerprint) {
    //     return false
    //   }
    //   return true
    // },
  },
}
</script>

<style lang="scss" scoped>
@import '../../../scss/globals';
.board-info {
  padding: 20px;
}

.inflight-box {
  margin-top: 10px;
}

.user-data-field-container {
  &:first-of-type {
    margin-top: 0;
  }
  margin-top: 17px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.25);
  padding-bottom: 12px;

  .user-data-field-header {
    font-family: 'SCP Bold';
    background-color: rgba(255, 255, 255, 0.05);
    padding: 1px 5px;
    display: inline-block;
    margin-bottom: 5px;
    position: relative;
    left: -5px;
  }

  .user-data-field .markdowned {
    font-family: 'SSP Regular';
    font-size: 110%;
  }

  .user-data-field.fingerprint {
    position: relative;
    left: -6px;
    margin-top: 7px;
  }

  .markdowned {
    p:last-of-type {
      margin-bottom: 0;
    }
  }

  .fingerprint-text {
    font-family: 'SCP Bold';
    word-wrap: break-word;
    width: 258px;
    font-size: 80%;
    margin-top: 10px;
    margin-bottom: 6px;
    color: $a-grey-600;
    background-color: rgba(255, 255, 255, 0.075);
    padding: 2px 6px;
    margin-left: -5px;
  }
}
</style>
