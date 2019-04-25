<template>
  <div class="user-sublocation">
    <div class="user-root">
      <div class="user-data-field-container">
        <!-- <div class="user-data-field-header">
          Info
        </div> -->
        <div class="user-data-field">
          <template v-if="!editPaneOpen">
            <a-no-content
              no-content-text="This user hasn't written a bio."
              quoteDisabled="true"
              v-if="visibleEntity.info.length === 0"
            >
            </a-no-content>
            <a-markdown :content="visibleEntity.info"></a-markdown>
            <div class="inflight-box" v-if="isVisibleInflightEntity">
              <a-inflight-info
                :status="visibleInflightStatus"
              ></a-inflight-info>
            </div>
            <div class="actions-rack" v-if="actionsVisible && isSelf">
              <a class="action edit" @click="toggleEditPane">
                <!-- todo: make it only visible when self, otherwise the button is meaningless -->
                Edit
              </a>
            </div>
          </template>
          <template v-else>
            <a-composer :spec="userEditExistingSpec"></a-composer>
          </template>
        </div>
      </div>
      <div class="user-data-field-container">
        <div class="user-data-field-header">
          URL
        </div>
        <div class="user-data-field location">
          <!-- <a class="field-action copy-link" @click="copyLink">
              Copy link
            </a> -->
          <a-fingerprint
            :fp="$store.state.currentUserEntity.fingerprint"
            :fontSize="11.5"
            :link="externalUrl"
          ></a-fingerprint>
        </div>
      </div>
      <div class="user-data-field-container">
        <div class="user-data-field-header">
          Created
        </div>
        <div class="user-data-field">
          <!-- {{timeSince($store.state.currentUserEntity.Creation)}} {{$store.state.currentUserEntity.Creation}} -->{{
            timeString($store.state.currentUserEntity.creation)
          }}
        </div>
      </div>
      <div class="user-data-field-container">
        <div class="user-data-field-header">
          Last Updated
        </div>
        <div class="user-data-field">
          {{ timeString($store.state.currentUserEntity.lastupdate) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var globalMethods = require('../../../services/globals/methods')
var fe = require('../../../services/feapiconsumer/feapiconsumer')
var mimobjs = require('../../../../../../protos/mimapi/structprotos_pb.js')
var mixins = require('../../../mixins/mixins')
export default {
  name: 'userroot',
  mixins: [mixins.localUserMixin],
  data(this: any): any {
    return {
      editPaneOpen: false,
      userEditExistingSpec: {
        fields: [
          {
            id: 'userInfo',
            visibleName: '',
            Info: '',
            placeholder: '',
            maxCharCount: 20480,
            heightRows: 5,
            previewDisabled: false,
            content: '',
            optional: true,
          },
        ],
        commitAction: this.submitEditExistingUser,
        commitActionName: 'SAVE',
        cancelAction: this.toggleEditPane,
        cancelActionName: 'CANCEL',
      },
    }
  },
  computed: {
    /*----------  Inflight computeds  ----------*/
    inflightEdits(this: any) {
      let iflEdits = []
      for (let val of this.$store.state.ambientStatus.inflights.keysList) {
        if (
          this.$store.state.currentUserEntity.fingerprint !==
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
      return this.$store.state.currentUserEntity
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
    isSelf(this: any) {
      if (globalMethods.IsUndefined(this.$store.state.currentUserEntity)) {
        return false
      }
      if (
        this.$store.state.currentUserEntity.fingerprint !==
        this.$store.state.localUser.fingerprint
      ) {
        return false
      }
      return true
    },
    externalUrl(this: any) {
      return 'aether:/' + this.$store.state.route.path
    },
  },
  methods: {
    timeSince(input: any) {
      return globalMethods.TimeSince(input)
    },
    timeString(input: any) {
      // All int64 values are provided to JS as strings by gRPC.. which means our timestamps are strings.
      if (parseInt(input) === 0) {
        return 'Never'
      }
      let d = new Date(input * 1000)
      return d.toDateString() + ', ' + d.toLocaleTimeString()
    },
    toggleEditPane(this: any) {
      console.log('toggle edit pane runs')
      if (this.editPaneOpen) {
        this.editPaneOpen = false
      } else {
        if (typeof this.$store.state.currentUserEntity.info !== 'undefined') {
          this.userEditExistingSpec.fields[0].content = this.visibleEntity.info
        }
        this.editPaneOpen = true
      }
    },
    /*----------  Edit user actions  ----------*/
    submitEditExistingUser(this: any, fields: any) {
      let userInfo = ''
      for (let val of fields) {
        if (val.id === 'userInfo') {
          userInfo = val.content
        }
      }
      let user = new mimobjs.Key()
      // Set board, thread, parent, body fields
      user.setInfo(userInfo) // < heads up, this is different
      let vm = this
      fe.SendUserContent(
        this.$store.state.currentUserEntity.fingerprint,
        user,
        function(resp: any) {
          vm.toggleEditPane()
          console.log(resp.toObject())
        }
      )
    },
    /*----------  Copy link  ----------*/
    // copyLink(this: any) {
    //   let elc = require('electron')
    //   elc.clipboard.writeText(this.externalUrl)
    // },
  },
}
</script>

<style lang="scss" scoped>
.user-root {
  padding: 20px;
}

.inflight-box {
  margin-top: 10px;
}

.user-data-field-container {
  &:first-of-type {
    margin-top: 0;
  }
  .no-content {
    margin: 0;
    margin-bottom: 10px;
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

  .user-data-field.location {
    padding-top: 7px;
    margin-left: -5px;
  }

  .markdowned {
    p:last-of-type {
      margin-bottom: 0;
    }
  }
}

.actions-rack {
  display: flex;
  .action {
    margin-right: 10px;
    font-family: 'SSP Semibold';
    cursor: pointer;
    user-select: none;
  }
}
</style>
<style lang="scss">
.user-data-field-container {
  .markdowned {
    p:last-of-type {
      margin-bottom: 0;
    }
  }
}
</style>
