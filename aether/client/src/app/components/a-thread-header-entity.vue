<template>
  <div
    class="thread-entity"
    v-if="threadReadyToRender"
    :class="{ uncompiled: uncompiled, 'no-user-present': !localUserPresent }"
  >
    <!-- Without above v-if it fails when currentThread===undefined -->
    <div class="signals-container" v-if="contentSignalsVisible">
      <div class="thread-actions"></div>
      <div class="thread-vote-count">
        <icon name="arrow-up"></icon> {{ threadVoteCount }}
      </div>
      <div class="thread-comment-count">
        <icon name="comment-alt"></icon> {{ currentThread.postscount }}
      </div>
    </div>
    <div
      class="image-container"
      v-show="imageLoadedSuccessfully || videoLoadedSuccessfully"
    >
      <div
        class="image-box"
        :style="'background-image: url(' + sanitisedLink + ')'"
        @click.prevent="openLightbox"
        v-show="imageLoadedSuccessfully"
      ></div>
      <video
        class="video-box"
        :src="sanitisedLink"
        preload="metadata"
        @click="openLightbox"
        v-show="videoLoadedSuccessfully"
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
          autoplay="true"
          loop="true"
          controls
          v-show="lightboxOpen"
        ></video>
      </div>

      <!-- @click="closeLightbox" -->
    </div>
    <div class="main-data-container">
      <div class="inflight-box" v-if="inflightBoxAtTopVisible">
        <a-inflight-info :status="visibleInflightStatus"></a-inflight-info>
      </div>
      <div class="thread-name">
        {{ currentThread.name }}
      </div>
      <div class="thread-link" v-show="sanitisedLink.length > 0">
        <a-link
          :link="currentThread.link"
          :linktext="currentThread.link"
        ></a-link>
      </div>
      <template v-if="!editPaneOpen">
        <a-markdown
          class="thread-body"
          :content="visibleThread.body"
        ></a-markdown>
        <div class="inflight-box" v-if="inflightBoxAtBottomVisible">
          <a-inflight-info :status="visibleInflightStatus"></a-inflight-info>
        </div>
      </template>
      <template v-else>
        <a-composer
          :spec="threadEditExistingSpec"
          v-if="editPaneOpen"
        ></a-composer>
      </template>
      <div class="meta">
        <div class="thread-owner">
          <a-username :owner="currentThread.owner"></a-username>
        </div>
        <div class="thread-datetime">
          <a-timestamp
            :creation="threadCreation"
            :lastupdate="threadLastUpdate"
          ></a-timestamp>
        </div>
        <div class="actions-rack" v-if="actionsVisible">
          <a
            class="action edit"
            v-show="currentThread.selfcreated"
            @click="toggleEditPane"
          >
            Edit
          </a>
          <a
            class="action report"
            v-show="
              !currentThread.selfcreated &&
                !reportSubmitted &&
                !reportSubmittedBefore
            "
            @click="toggleReportPane"
          >
            Report
          </a>
          <a class="action moderate" v-if="isMod" @click="toggleModActionPane">
            Moderate
          </a>
          <a class="action copy-link" @click="copyLink">
            Copy link
          </a>
          <div class="link-copied-text" v-show="linkCopied">
            Link copied to clipboard
          </div>
        </div>
        <div class="actions-rack" v-if="uncompiled">
          <a class="action copy-link" @click="goTo">
            Go to thread
          </a>
        </div>
      </div>
      <div class="modtools">
        <div
          class="reports-list-container modtools-block"
          v-if="reportsListVisible"
        >
          <div
            class="confirmation-container"
            v-if="reportSubmitted || reportSubmittedBefore"
            :class="{ 'no-bottom-margin': reportsList.length === 0 }"
          >
            <div class="inprogress" v-if="reportSubmitted">
              Thanks! Your report is
              <b><em>currently being submitted</em></b> to the mods of this
              community.
            </div>
            <div class="complete" v-if="reportSubmittedBefore">
              Thanks! Your report is <b><em>successfully submitted</em></b> to
              the mods of this community. You can see it below.
            </div>
          </div>
          <div
            class="reports-header"
            v-if="
              reportsList.length > 1 ||
                (!reportSubmittedBefore && !reportSubmitted)
            "
          >
            User reports
          </div>
          <div
            class="report-container"
            v-for="report in reportsList"
            :class="{
              'own-report':
                report.sourcefp === $store.state.localUser.fingerprint,
            }"
          >
            <div class="meta-container">
              <div class="post-datetime">
                <a-timestamp
                  :creation="report.creation"
                  :lastupdate="report.lastupdate"
                ></a-timestamp>
              </div>
              <div class="author">
                <template
                  v-if="report.sourcefp !== $store.state.localUser.fingerprint"
                >
                  <router-link :to="'/user/' + report.sourcefp"
                    >See author</router-link
                  >
                </template>
                <template v-else>
                  <a-username :owner="report.sourcefp"></a-username>
                </template>
              </div>
            </div>
            <template v-if="report.reason.length > 0">
              <div class="report-text">{{ report.reason }}</div>
            </template>
            <template v-else>
              <div class="report-text no-reason">(no reason given)</div>
            </template>
          </div>
        </div>
        <div class="report-composer-container" v-if="reportPaneOpen">
          <a-composer :spec="reportSpec"></a-composer>
        </div>
        <template v-if="isMod">
          <div
            class="mod-actions-container modtools-block"
            v-if="modActionsVisible"
          >
            <div class="mod-actions-header">
              Moderation
              <a-info-marker
                header="You can issue moderation commands here. "
                text="<p><b>Approve</b> marks that content as okay, and it will remain visible even if it receives <em>Delete</em>s from other mods.</p><p><b>Delete</b> marks that content as a rule violation, and it will be removed.</p><p><b>Ignore</b> is only available in <em>Reports</em> page, and it marks the content as not interesting either way. As a result, it will no longer show up on the <em>Reports</em> page.</p> <p><a href='#/modship'><b>More info about moderation</b></a></p>"
              ></a-info-marker>
            </div>
            <div
              class="buttons-row"
              v-show="
                !(modApprovalPaneOpen || modDeletePaneOpen) && !modActionTaken
              "
            >
              <a
                class="button is-success is-outlined"
                @click="toggleModApprovalPane"
              >
                APPROVE
              </a>
              <a
                class="button is-success  is-outlined"
                v-show="isinreportsview"
                @click="submitModIgnore"
              >
                IGNORE
              </a>
              <a
                class="button is-success  is-outlined"
                @click="toggleModDeletePane"
              >
                DELETE
              </a>
            </div>
            <div
              class="modapproval-composer-container"
              v-if="modApprovalPaneOpen"
            >
              <a-composer :spec="modApprovalSpec"></a-composer>
            </div>
            <div
              class="approval-confirmation-container"
              v-if="modApprovalSubmitted || modApprovalSubmittedBefore"
            >
              <div class="inprogress" v-if="modApprovalSubmitted">
                Thanks! Your approval is
                <b><em>currently being minted</em></b> for submission.
              </div>
              <div class="complete" v-if="modApprovalSubmittedBefore">
                Thanks! Your approval is <b><em>successfully submitted.</em></b>
              </div>
            </div>
            <div class="moddelete-composer-container" v-if="modDeletePaneOpen">
              <a-composer :spec="modDeleteSpec"></a-composer>
            </div>
            <div
              class="delete-confirmation-container"
              v-if="modDeleteSubmitted || modDeleteSubmittedBefore"
            >
              <div class="inprogress" v-if="modDeleteSubmitted">
                Thanks! Your deletion is
                <b><em>currently being minted</em></b> for submission.
              </div>
              <div class="complete" v-if="modDeleteSubmittedBefore">
                Thanks! Your deletion is <b><em>successfully submitted.</em></b>
              </div>
            </div>
            <div
              class="ignore-confirmation-container"
              v-if="modIgnoreSubmitted || modIgnoreSubmittedBefore"
            >
              <div class="inprogress" v-if="modIgnoreSubmitted">
                Thanks! Your ignore is <b><em>currently being processed</em></b
                >.
              </div>
              <div class="complete" v-if="modIgnoreSubmittedBefore">
                Thanks! Your ignore is
                <b><em>successfully submitted.</em></b> You won't be notified of
                new reports for this content in the Reports tab.
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
    <a-ballot
      v-if="actionsVisible"
      :contentsignals="currentThread.compiledcontentsignals"
      :boardfp="currentThread.board"
      :threadfp="currentThread.fingerprint"
    ></a-ballot>
  </div>
</template>

<script lang="ts">
// var Vuex = require('../../../node_modules/vuex').default
var globalMethods = require('../services/globals/methods')
var mixins = require('../mixins/mixins')
var fe = require('../services/feapiconsumer/feapiconsumer')
var mimobjs = require('../../../../protos/mimapi/structprotos_pb.js')
var Mousetrap = require('../../../node_modules/mousetrap')
export default {
  name: 'a-thread-header-entity',
  mixins: [mixins.localUserMixin],
  // props: ['inflightStatus', 'thread', 'uncompiled'],
  props: {
    inflightStatus: {
      type: Object,
      default: function() {
        return undefined
      },
    },
    uncompiled: {
      type: Boolean,
      default: false,
    },
    thread: {
      type: Object,
      default: function() {
        return undefined
      },
    },
    isinreportsview: {
      type: Boolean,
      default: false,
    },
  },
  // ^ Unused as of now, since we do not allow inflight threads to be opened, but could be useful in the future.
  data(this: any): any {
    return {
      linkCopied: false,
      copyLinkInProgress: false,
      imageLoadedSuccessfully: false,
      videoLoadedSuccessfully: false,
      hasUpvoted: false,
      hasDownvoted: false,
      editPaneOpen: false,
      reportPaneOpen: false,
      modActionsPaneOpen: false,
      modApprovalPaneOpen: false,
      modDeletePaneOpen: false,
      lightboxOpen: false,
      reportSubmitted: false,
      modApprovalSubmitted: false,
      modDeleteSubmitted: false,
      modIgnoreSubmitted: false,
      threadEditExistingSpec: {
        fields: [
          {
            id: 'threadBody',
            visibleName: '',
            description: '',
            placeholder: '',
            maxCharCount: 20480,
            heightRows: 5,
            previewDisabled: false,
            content: '',
            optional: false,
          },
        ],
        commitAction: this.submitEditExistingThread,
        commitActionName: 'SAVE',
        cancelAction: this.toggleEditPane,
        cancelActionName: 'CANCEL',
      },
      reportSpec: {
        fields: [
          {
            id: 'reportReason',
            emptyWarningDisabled: true,
            visibleName: 'Report to mods',
            description:
              'This report will go to the mods. (Heads up - the reports are publicly visible to everyone.)',
            placeholder: "What's the reason?",
            maxCharCount: 256,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: false,
          },
        ],
        draft: {
          parentFp: this.thread.board,
          contentType: 'report',
        },
        commitAction: this.submitReport,
        commitActionName: 'REPORT',
        cancelAction: this.toggleReportPane,
        cancelActionName: 'CANCEL',
      },
      modApprovalSpec: {
        fields: [
          {
            id: 'approvalReason',
            emptyWarningDisabled: true,
            visibleName: 'Approve',
            description:
              "Optional. Enter the reason why you're approving this. ",
            placeholder: 'Approval reason',
            maxCharCount: 256,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: true,
          },
        ],
        draft: {
          parentFp: this.thread.board,
          contentType: 'modApproval',
        },
        commitAction: this.submitModApproval,
        commitActionName: 'APPROVE',
        cancelAction: this.toggleModApprovalPane,
        cancelActionName: 'CANCEL',
      },
      modDeleteSpec: {
        fields: [
          {
            id: 'deleteReason',
            emptyWarningDisabled: true,
            visibleName: 'Delete',
            description:
              "Optional. Enter the reason why you're deleting this. ",
            placeholder: 'Deletion reason',
            maxCharCount: 256,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: true,
          },
        ],
        draft: {
          parentFp: this.thread.board,
          contentType: 'modDelete',
        },
        commitAction: this.submitModDelete,
        commitActionName: 'DELETE',
        cancelAction: this.toggleModDeletePane,
        cancelActionName: 'CANCEL',
      },
    }
  },
  computed: {
    localUserPresent(this: any) {
      if (this.localUserArrived && !this.localUserExists) {
        return false
      }
      return true
    },
    // ...Vuex.mapState(['currentThread']),
    currentThread(this: any) {
      return this.thread
    },
    threadCreation(this: any) {
      // These are necessary because in uncompiled entities, these are in thread.provable.creation, but in compiled ones it's in thread.creation.
      if (this.uncompiled) {
        if (typeof this.thread.provable === 'undefined') {
          return 0
        }
        return this.thread.provable.creation
      }
      if (typeof this.mostRecentInflightEdit !== 'undefined') {
        if (this.isVisibleInflightEntity) {
          if (
            typeof this.mostRecentInflightEdit.entity.provable === 'undefined'
          ) {
            return 0
          }
          return this.mostRecentInflightEdit.entity.provable.creation
        }
      }
      return this.thread.creation
    },
    threadLastUpdate(this: any) {
      if (this.uncompiled) {
        if (typeof this.thread.updateable === 'undefined') {
          return 0
        }
        return this.thread.updateable.lastupdate
      }
      if (typeof this.mostRecentInflightEdit !== 'undefined') {
        if (this.isVisibleInflightEntity) {
          if (
            typeof this.mostRecentInflightEdit.entity.updateable === 'undefined'
          ) {
            return 0
          }
          return this.mostRecentInflightEdit.entity.updateable.lastupdate
        }
      }
      return this.thread.lastupdate
    },
    threadFingerprint(this: any) {
      if (this.uncompiled) {
        return this.thread.provable.fingerprint
      }
      return this.thread.fingerprint
    },
    threadVoteCount(this: any) {
      if (typeof this.currentThread.compiledcontentsignals === 'undefined') {
        return 0
      }
      return (
        this.currentThread.compiledcontentsignals.upvotes -
        this.currentThread.compiledcontentsignals.downvotes
      )
    },
    threadReadyToRender(this: any) {
      if (this.uncompiled) {
        // If uncompiled and this code is running, the thread data is already there, and the compiled data won't ever be there since it's uncompiled, which means we're ready to go.
        return true
      }
      if (this.isVisibleInflightEntity) {
        return true
      }
      if (typeof this.currentThread.compiledcontentsignals === 'undefined') {
        return false
      }
      return true
    },
    sanitisedLink(this: any) {
      if (typeof this.currentThread.link === 'undefined') {
        return ''
      }
      if (this.currentThread.link.length === 0) {
        return ''
      }
      if (
        this.currentThread.link.substring(0, 8) === 'https://' ||
        this.currentThread.link.substring(0, 7) === 'http://'
      ) {
        return this.currentThread.link
      }
      return 'http://' + this.currentThread.link
    },
    /*----------  Inflight computeds  ----------*/
    inflightEdits(this: any) {
      let iflEdits = []
      for (let val of this.$store.state.ambientStatus.inflights.threadsList) {
        if (this.threadFingerprint !== val.entity.provable.fingerprint) {
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
    visibleThread(this: any) {
      if (typeof this.mostRecentInflightEdit !== 'undefined') {
        return this.mostRecentInflightEdit.entity
      }
      return this.currentThread
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
    /*----------  Visibility calculations  ----------*/
    /*
        This matters because we're starting to have more states. One state that we have is inflight, and in this state, we don't have results of the compilation, nor do we have things such as a fingerprint. As a result, both signals and actions are disabled.

        The second state we have is the uncompiled state. This is how we show entities that we have pulled directly from the backend. A major use case for this is to show a user's profile, to show posts, threads, boards that the user has created, and so on.
      */
    contentSignalsVisible(this: any) {
      if (this.localUserReadOnly) {
        return false
      }
      if (this.uncompiled) {
        return false
      }
      return true
    },
    actionsVisible(this: any) {
      if (this.localUserReadOnly) {
        return false
      }
      // if inflight and compiled =  true
      // if inflight and uncompiled = false
      // if just uncompiled = false
      if (this.isVisibleInflightEntity) {
        if (!this.uncompiled) {
          return true
        }
        return false
      }
      if (this.uncompiled) {
        return false
      }
      return true
    },
    inflightBoxAtTopVisible(this: any) {
      // We use this in two different places. Box at top is suitable for lists, and box at bottom is suitable for full page views.
      if (this.uncompiled && this.isVisibleInflightEntity) {
        return true
      }
      return false
    },
    inflightBoxAtBottomVisible(this: any) {
      // We use this in two different places. Box at top is suitable for lists, and box at bottom is suitable for full page views.
      if (this.isVisibleInflightEntity && !this.uncompiled) {
        return true
      }
      return false
    },
    /*----------  Reports  ----------*/
    reportSubmittedBefore(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      if (this.thread.compiledcontentsignals.selfreported) {
        return true
      }
      return false
    },
    selfReportText(this: any): string {
      if (this.uncompiled || this.isInflightEntity) {
        return ''
      }
      if (globalMethods.IsUndefined(this.$store.state.localUser)) {
        return ''
      }
      for (let val of this.thread.compiledcontentsignals.reportsList) {
        if (val.sourcefp === this.$store.state.localUser.fingerprint) {
          return val.reason
        }
      }
      return ''
    },
    reportsListVisible(this: any) {
      // If uncompiled or inflight, always false
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      // If mod mode is enabled, true if reports are present
      if (
        this.$store.state.modModeEnabledArrived &&
        this.$store.state.modModeEnabled
      ) {
        if (this.thread.compiledcontentsignals.reportsList.length > 0) {
          return true
        }
      }
      // If mod mode is not enabled, and not uncompiled and not inflight, if a report has been submitted now or before, true
      if (this.reportSubmitted || this.reportSubmittedBefore) {
        return true
      }
      return false
    },
    reportsList(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return []
      }
      if (
        this.$store.state.modModeEnabledArrived &&
        this.$store.state.modModeEnabled
      ) {
        return this.thread.compiledcontentsignals.reportsList
      }
      let nonModReports: any = []
      // only reports created by self
      if (globalMethods.IsUndefined(this.$store.state.localUser)) {
        return nonModReports
      }
      for (let val of this.thread.compiledcontentsignals.reportsList) {
        if (val.sourcefp === this.$store.state.localUser.fingerprint) {
          nonModReports.push(val)
        }
      }
      return nonModReports
    },
    /*----------  Mod mode and actions  ----------*/
    modDeleteSubmittedBefore(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      if (this.thread.compiledcontentsignals.selfmodblocked) {
        return true
      }
      return false
    },
    modApprovalSubmittedBefore(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      if (this.thread.compiledcontentsignals.selfmodapproved) {
        return true
      }
      return false
    },
    modIgnoreSubmittedBefore(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      if (this.thread.compiledcontentsignals.selfmodignored) {
        return true
      }
      return false
    },
    isMod(this: any) {
      if (
        this.$store.state.modModeEnabledArrived &&
        this.$store.state.modModeEnabled
      ) {
        return true
      }
      return false
    },
    modActionSubmitted(this: any) {
      if (
        this.modApprovalSubmitted ||
        this.modDeleteSubmitted ||
        this.modIgnoreSubmitted
      ) {
        return true
      }
      return false
    },
    modActionSubmittedBefore(this: any) {
      if (
        this.modApprovalSubmittedBefore ||
        this.modDeleteSubmittedBefore ||
        this.modIgnoreSubmittedBefore
      ) {
        return true
      }
      return false
    },
    modActionsVisible(this: any) {
      if (this.reportsListVisible) {
        return true
      }
      if (this.modActionsPaneOpen) {
        return true
      }
      if (
        this.modActionSubmitted ||
        this.modActionSubmittedBefore ||
        this.modIgnoreSubmittedBefore
      ) {
        return true
      }
      return false
    },
    modActionTaken(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      if (
        this.modApprovalSubmitted ||
        this.modApprovalSubmittedBefore ||
        this.modDeleteSubmitted ||
        this.modDeleteSubmittedBefore ||
        this.modIgnoreSubmitted ||
        this.modIgnoreSubmittedBefore
      ) {
        return true
      }
      return false
    },
    externalUrl(this: any) {
      return 'aether:/' + this.$store.state.route.path
    },
    uncompiledDeepLink(this: any) {
      return (
        '/board/' +
        this.thread.board +
        '/thread/' +
        this.thread.provable.fingerprint
      )
    },
  },
  methods: {
    /*----------  Lightbox open/close  ----------*/
    openLightbox(this: any) {
      let vm = this
      Mousetrap.bind('esc', function() {
        vm.closeLightbox()
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
    getOwnerName(owner: any): string {
      return globalMethods.GetOwnerName(owner)
    },
    /*----------  Edit thread actions  ----------*/
    toggleEditPane(this: any) {
      if (this.editPaneOpen) {
        this.editPaneOpen = false
      } else {
        this.threadEditExistingSpec.fields[0].content = this.visibleThread.body
        this.editPaneOpen = true
      }
    },
    submitEditExistingThread(this: any, fields: any) {
      let threadBody = ''
      for (let val of fields) {
        if (val.id === 'threadBody') {
          threadBody = val.content
        }
      }
      let thread = new mimobjs.Thread()
      // Set board, thread, parent, body fields
      thread.setBody(threadBody)
      thread.setBoard(this.currentThread.board)
      let vm = this
      vm.toggleEditPane()
      fe.SendThreadContent(this.currentThread.fingerprint, thread, function(
        resp: any
      ) {
        console.log(resp.toObject())
      })
    },
    /*----------  Report actions  ----------*/
    toggleReportPane(this: any) {
      if (this.reportPaneOpen) {
        this.reportPaneOpen = false
      } else {
        this.reportPaneOpen = true
      }
    },
    submitReport(this: any, fields: any) {
      let reportReason = ''
      for (let val of fields) {
        if (val.id === 'reportReason') {
          reportReason = val.content
        }
      }
      let vm = this
      fe.ReportToMod(
        this.currentThread.fingerprint,
        '',
        reportReason,
        this.currentThread.board,
        this.currentThread.fingerprint,
        function(resp: any) {
          console.log(resp.toObject())
          vm.toggleReportPane()
          vm.reportSubmitted = true
        }
      )
    },
    handleImageLoadSuccess(this: any) {
      this.imageLoadedSuccessfully = true
    },
    handleVideoLoadSuccess(this: any) {
      // console.log('video loaded successfully')
      this.videoLoadedSuccessfully = true
    },
    /*----------  Moderation actions  ----------*/
    toggleModActionPane(this: any) {
      if (this.modActionsPaneOpen) {
        this.modActionsPaneOpen = false
      } else {
        this.modActionsPaneOpen = true
      }
    },
    /*----------  Approval modaction  ----------*/
    toggleModApprovalPane(this: any) {
      if (this.modApprovalPaneOpen) {
        this.modApprovalPaneOpen = false
      } else {
        this.modApprovalPaneOpen = true
      }
    },
    submitModApproval(this: any, fields: any) {
      let approvalReason = ''
      for (let val of fields) {
        if (val.id === 'approvalReason') {
          approvalReason = val.content
        }
      }
      let vm = this
      fe.ModApprove(
        this.currentThread.fingerprint,
        '',
        approvalReason,
        this.currentThread.board,
        this.currentThread.fingerprint,
        function(resp: any) {
          console.log(resp.toObject())
          vm.toggleModApprovalPane()
          vm.modApprovalSubmitted = true
        }
      )
    },
    /*----------  Delete modaction  ----------*/
    toggleModDeletePane(this: any) {
      if (this.modDeletePaneOpen) {
        this.modDeletePaneOpen = false
      } else {
        this.modDeletePaneOpen = true
      }
    },
    submitModDelete(this: any, fields: any) {
      let deleteReason = ''
      for (let val of fields) {
        if (val.id === 'deleteReason') {
          deleteReason = val.content
        }
      }
      let vm = this
      fe.ModDelete(
        this.currentThread.fingerprint,
        '',
        deleteReason,
        this.currentThread.board,
        this.currentThread.fingerprint,
        function(resp: any) {
          console.log(resp.toObject())
          vm.toggleModDeletePane()
          vm.modDeleteSubmitted = true
        }
      )
    },
    /*----------  Ignore modaction  ----------*/
    submitModIgnore(this: any) {
      let vm = this
      fe.ModIgnore(
        this.currentThread.fingerprint,
        '',
        '',
        this.currentThread.board,
        this.currentThread.fingerprint,
        function(resp: any) {
          console.log(resp.toObject())
          vm.modIgnoreSubmitted = true
        }
      )
    },
    /*----------  Copy link  ----------*/
    copyLink(this: any) {
      if (this.copyLinkInProgress) {
        return
      }
      this.copyLinkInProgress = true
      let elc = require('electron')
      elc.clipboard.writeText(this.externalUrl)
      this.linkCopied = true
      // Prevent repeated clicks while in progress
      let vm = this
      setTimeout(function() {
        vm.linkCopied = false
        vm.copyLinkInProgress = false
      }, 1250)
    },
    /*----------  Go to entity (uncompiled)  ----------*/
    goTo(this: any) {
      console.log(this.uncompiledDeepLink)
      this.$router.push(this.uncompiledDeepLink)
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
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

.thread-entity {
  display: block;
  display: flex;
  padding: 15px 5px;
  margin: 0 20px;
  margin-top: 20px;
  color: $a-grey-800;
  &.no-user-present {
    padding-bottom: 25px;
    border-bottom: 3px solid rgba(255, 255, 255, 0.05);
  }
  .ballot {
    // padding-top: 75px;
  }
  &:hover {
    .ballot {
      visibility: visible;
    }
  }

  &.uncompiled {
    .main-data-container .thread-name {
      cursor: default;
    }
  }

  .signals-container {
    width: 64px;
    font-size: 110%;
    margin: auto;
    margin-top: 0;
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

  .image-container {
    width: 80px;
    padding: 6px 8px;

    .image-box {
      height: 100%;
      max-height: 100px;
      overflow: hidden;
      border-radius: 2px;
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
    min-width: 0;

    .composer {
      margin-bottom: 15px;
    }
    .thread-name {
      font-size: 120%;
      cursor: pointer;
    }

    .thread-link {
      font-size: 95%;
      a {
        word-break: break-all;
      }
    }

    .thread-body {
      margin-top: 25px;
      font-family: 'SSP Regular';
      font-size: 110%;
    }
  }
  .meta {
    display: flex;
    .thread-owner {
    }
    .thread-datetime {
      margin-left: 10px;
      font-family: 'SSP Regular Italic';
      color: $a-grey-600;
    }
  }
}

.actions-rack {
  display: flex;
  margin-left: 12px;
  .action {
    margin: auto;
    margin-bottom: 0;
    margin-right: 10px;
    font-family: 'SSP Semibold';
    font-size: 90%;
    user-select: none;
  }
}

.modtools {
  font-size: 16px;
  font-family: 'SSP Bold';
  .modtools-block {
    margin-top: 5px;
    &:first-of-type {
      margin-top: 20px;
    }
  }
}

.report-composer-container {
  padding-top: 30px;
}

.reports-list-container {
  // padding: 15px 20px 20px 20px;
  padding: 20px;
  background-color: rgba(0, 0, 0, 0.25);
  margin-top: 20px;
  border-radius: 3px; // font-size: 105%;
  color: $a-grey-600;

  .reports-header {
    font-family: 'SSP Bold';
    font-size: 110%;
    color: $a-grey-800;
  }
  .report-container {
    border-left: 3px solid rgba(255, 255, 255, 0.15);
    padding-left: 20px;
    margin-top: 15px;

    &.own-report {
      border-color: $a-cerulean;
    }

    .report-text {
      font-family: 'SSP Regular';
      font-size: 110%;
      &.no-reason {
        font-family: 'SSP Regular Italic';
      }
    }

    .meta-container {
      display: flex;
      .post-datetime {
        margin-right: 10px;
        margin-left: 0;
        font-family: 'SSP Regular Italic';
        color: $a-grey-600;
      }
    }
  }
}

.mod-actions-container {
  padding: 20px;
  background-color: rgba(0, 0, 0, 0.25);
  margin-top: 20px;
  border-radius: 3px; // font-size: 110%;
  color: $a-grey-600;
  .mod-actions-header {
    font-family: 'SSP Bold';
    margin-bottom: 15px;
    font-size: 110%;
    color: $a-grey-800;
  }
  .buttons-row {
    margin-top: 25px;
  }
  .button {
    margin-right: 5px;
  }

  .modapproval-composer-container,
  .moddelete-composer-container {
    margin-top: 15px;
  }

  .approval-confirmation-container,
  .delete-confirmation-container,
  .ignore-confirmation-container {
    margin-top: 15px;
    margin-bottom: 0;
    font-size: 110%;
  }
}

.confirmation-container,
.approval-confirmation-container,
.delete-confirmation-container,
.ignore-confirmation-container {
  font-family: 'SSP Regular';
  font-size: 110%;
  margin-bottom: 15px;
  &.no-bottom-margin {
    margin-bottom: 0;
  }
  b em {
    font-family: 'SSP Bold Italic';
  }
}
.link-copied-text {
  font-family: 'SSP Regular Italic';
  user-select: none;
  cursor: default;
  opacity: 0;
  animation-duration: 1.25s;
  animation-name: DELAY_INVISIBLE;
}

@keyframes DELAY_INVISIBLE {
  0% {
    opacity: 1;
  }
  60% {
    opacity: 1;
  }
  100% {
    opacity: 0;
  }
}
</style>

<!--
 creation:1535323102
 lastupdate:0
 reason:"This is too coooll. Waay too cool. It shouldn't be allowed."
 sourcefp:"85e03dda365a42caf48129e5832b80305868c470c4afccb400190b73e9c959b5"

  -->
