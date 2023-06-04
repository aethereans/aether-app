<template>
  <div
    class="post"
    :class="{
      inflight: isInflightEntity,
      'notification-highlight': isNotificationsHighlight,
      'self-owned': post.selfcreated || isVisibleInflightEntity,
    }"
    :id="'post-' + post.fingerprint"
  >
    <div class="inflight-box" v-if="isVisibleInflightEntity">
      <a-inflight-info
        :status="visibleInflightStatus"
        :refresherFunc="refresh"
      ></a-inflight-info>
    </div>
    <div class="expand-container" @click="toggleCollapse">
      <template v-if="!collapsed">
        <icon name="minus-circle"></icon>
      </template>
      <template v-if="collapsed">
        <icon name="plus-circle"></icon>
      </template>
    </div>
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
    <template v-if="!collapsed">
      <div class="post-content">
        <div class="content-text-container">
          <a-composer
            :spec="postEditExistingSpec"
            v-if="editPaneOpen"
          ></a-composer>
          <template v-if="useBodyPreview">
            <a-markdown
              class="content-text"
              :content="postBodyPreview"
            ></a-markdown>
          </template>
          <template v-else>
            <a-markdown
              class="content-text"
              :content="visiblePost.body"
            ></a-markdown>
          </template>

          <div class="actions-rack" v-if="actionsVisible">
            <a class="action reply" @click="toggleReplyPane"> Reply </a>
            <a
              class="action edit"
              v-show="
                post.selfcreated && !isDeletedByCreator && !inflightInProgress
              "
              @click="toggleEditPane"
            >
              Edit
            </a>
            <a
              class="action delete"
              v-show="
                post.selfcreated && !isDeletedByCreator && !inflightInProgress
              "
              @click="toggleSelfDeleteConfirm"
            >
              Delete
            </a>
            <a
              class="action report"
              v-show="
                !post.selfcreated &&
                !reportSubmitted &&
                !reportSubmittedBefore &&
                !inflightInProgress
              "
              @click="toggleReportPane"
            >
              Report
            </a>
            <a
              class="action moderate"
              v-if="isMod"
              v-show="!inflightInProgress"
              @click="toggleModActionPane"
            >
              Moderate
            </a>

            <a class="action copy-link" @click="copyLink"> Copy link </a>
            <div class="link-copied-text" v-show="linkCopied">
              Link copied to clipboard
            </div>
          </div>
          <div
            class="actions-rack"
            v-if="
              (isinmodactionsview && !post.compiledcontentsignals.modblocked) ||
              uncompiled
            "
          >
            <a class="action copy-link" @click="goTo"> Go to post </a>
          </div>
          <div class="self-delete-confirm-text" v-show="selfDeleteConfirmOpen">
            Are you sure you want to delete this? <a @click="selfDelete">Yes</a>
            <a @click="selfDeleteConfirmOpen = false">Cancel</a>
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
                  Thanks! Your report is
                  <b><em>successfully submitted</em></b> to the mods of this
                  community. You can see it below.
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
                      v-if="
                        report.sourcefp !== $store.state.localUser.fingerprint
                      "
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
                    !(modApprovalPaneOpen || modDeletePaneOpen) &&
                    !modActionTaken
                  "
                >
                  <!-- <a
                    class="button is-success is-outlined"
                    @click="toggleModApprovalPane"
                  >
                    APPROVE
                  </a> -->
                  <a
                    class="button is-success is-outlined"
                    v-show="isinreportsview"
                    @click="submitModIgnore"
                  >
                    IGNORE
                  </a>
                  <a
                    class="button is-success is-outlined"
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
                    Thanks! Your approval is
                    <b><em>successfully submitted.</em></b>
                  </div>
                </div>
                <div
                  class="moddelete-composer-container"
                  v-if="modDeletePaneOpen"
                >
                  <a-composer :spec="modDeleteSpec"></a-composer>
                </div>
                <div
                  class="delete-confirmation-container"
                  v-if="modDeleteSubmitted || modDeleteSubmittedBefore"
                >
                  <div class="inprogress" v-if="modDeleteSubmitted">
                    Thanks! Your deletion request is
                    <b><em>currently being minted</em></b> for submission.
                  </div>
                  <div class="complete" v-if="modDeleteSubmittedBefore">
                    Thanks! Your deletion request is
                    <b><em>successfully submitted.</em></b>
                  </div>
                </div>
                <div
                  class="ignore-confirmation-container"
                  v-if="modIgnoreSubmitted || modIgnoreSubmittedBefore"
                >
                  <div class="inprogress" v-if="modIgnoreSubmitted">
                    Thanks! Your ignore is
                    <b><em>currently being processed</em></b
                    >.
                  </div>
                  <div class="complete" v-if="modIgnoreSubmittedBefore">
                    Thanks! Your ignore is
                    <b><em>successfully submitted.</em></b> You won't be
                    notified of new reports for this content in the Reports tab.
                  </div>
                </div>
              </div>
            </template>
            <template v-if="isinmodactionsview">
              <div class="mod-actions-taken-container">
                <div class="block header-block">
                  <div class="mod-actions-taken-header">
                    Mod actions taken
                    <a-info-marker
                      header="These are the mod actions taken on this entity by various people. "
                      text="<p>These might be mods that are by default enabled, for example, by having created the community, those that are elected, or those you have marked as mods yourself. </p><p>These might be also mods that you have not approved, but those will not affect the final result.</p>"
                    ></a-info-marker>
                  </div>
                </div>
                <div
                  class="block iterator-block mod-actions-taken-iterator"
                  v-for="modActionTaken in post.compiledcontentsignals
                    .modblocksList"
                  :class="{
                    'own-report':
                      modActionTaken.sourcefp ===
                      $store.state.localUser.fingerprint,
                  }"
                >
                  <div class="mod-action-taken-data">
                    <div class="meta-container">
                      <div class="post-datetime">
                        <a-timestamp
                          :creation="modActionTaken.creation"
                          :lastupdate="modActionTaken.lastupdate"
                        ></a-timestamp>
                      </div>
                      <div class="author">
                        <template
                          v-if="
                            modActionTaken.sourcefp !==
                            $store.state.localUser.fingerprint
                          "
                        >
                          <router-link :to="'/user/' + modActionTaken.sourcefp"
                            >See author</router-link
                          >
                        </template>
                        <template v-else>
                          <a-username
                            :owner="modActionTaken.sourcefp"
                          ></a-username>
                        </template>
                      </div>
                    </div>
                    <template v-if="modActionTaken.reason.length > 0">
                      <div class="report-text">{{ modActionTaken.reason }}</div>
                    </template>
                    <template v-else>
                      <div class="report-text no-reason">(no reason given)</div>
                    </template>
                  </div>
                  <div class="mod-action-taken-action">
                    <a class="button is-white is-small is-outlined is-disabled">
                      DELETION REQUEST
                    </a>
                  </div>
                </div>
                <div class="block iterator-block">
                  <div class="final-state-container">
                    <div class="final-state-text">Result &nbsp;</div>
                    <template v-if="post.compiledcontentsignals.modblocked">
                      <div class="final-state-result">
                        <a class="button is-danger is-small is-disabled">
                          DELETED
                        </a>
                      </div>
                    </template>
                    <template v-else>
                      <div class="final-state-result">
                        <a class="button is-white is-small is-disabled">
                          NOT DELETED
                        </a>
                      </div>
                    </template>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
        <a-ballot
          v-if="actionsVisible"
          :contentsignals="post.compiledcontentsignals"
          :boardfp="post.board"
          :threadfp="post.thread"
        ></a-ballot>
      </div>
      <a-composer :spec="postComposerSpec" v-if="replyPaneOpen"></a-composer>
      <a-post
        v-for="iflChild in inflightChildren"
        :key="iflChild.status.requestedtimestamp + iflChild.entity.body"
        :post="iflChild.entity"
        :inflightStatus="iflChild.status"
      ></a-post>
      <a-post
        v-for="child in post.childrenList"
        :key="child.fingerprint"
        :post="child"
        :notificationparent="notificationparent"
        :notificationhighlights="notificationhighlights"
        :notificationfocus="notificationfocus"
      ></a-post>
    </template>
  </div>
</template>

<script lang="ts">
var globalMethods = require('../services/globals/methods')
var mixins = require('../mixins/mixins')
var fe = require('../services/feapiconsumer/feapiconsumer')
var mimobjs = require('../../../../protos/mimapi/mimapi_pb.js')
var vuexStore = require('../store/index').default
export default {
  name: 'a-post',
  mixins: [mixins.localUserMixin],
  // props: ['post', 'inflightStatus', 'uncompiled', 'notificationparent', 'notificationhighlights', 'notificationfocus'],
  props: {
    post: {
      type: Object,
      default: function () {
        return undefined
      },
    },
    inflightStatus: {
      type: Object,
      default: function () {
        return undefined
      },
    },
    uncompiled: {
      type: Boolean,
      default: false,
    },
    notificationparent: {
      type: String,
      default: '',
    },
    notificationhighlights: {
      type: Array,
      default: function () {
        return []
      },
    },
    notificationfocus: {
      type: String,
      default: '',
    },
    isinreportsview: {
      type: Boolean,
      default: false,
    },
    isinmodactionsview: {
      type: Boolean,
      default: false,
    },
    useBodyPreview: {
      type: Boolean,
      default: false,
    },
  },
  data(this: any): any {
    return {
      linkCopied: false,
      copyLinkInProgress: false,
      hasUpvoted: false,
      hasDownvoted: false,
      replyPaneOpen: false,
      editPaneOpen: false,
      selfDeleteConfirmOpen: false,
      reportPaneOpen: false,
      modActionsPaneOpen: false,
      modApprovalPaneOpen: false,
      modDeletePaneOpen: false,
      collapsed: false,
      reportSubmitted: false,
      modApprovalSubmitted: false,
      modDeleteSubmitted: false,
      modIgnoreSubmitted: false,
      postComposerSpec: {
        fields: [
          {
            id: 'postBody',
            visibleName: '',
            description: '',
            placeholder: 'Post a reply',
            maxCharCount: 20480,
            heightRows: 5,
            previewDisabled: false,
            content: '',
            optional: false,
          },
        ],
        draft: {
          parentFp: this.post.fingerprint,
          contentType: 'post',
        },
        commitAction: this.submitPost,
        commitActionName: 'SUBMIT',
        cancelAction: this.toggleReplyPane,
        cancelActionName: 'CANCEL',
        cancelRetainsDraft: true,
        autofocus: false,
        // ^ Why? Because autofocus focuses on the first entity on the page.. Not great when you have multiple possible composers on the same page.
      },
      postEditExistingSpec: {
        fields: [
          {
            id: 'postBody',
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
        commitAction: this.submitEditExistingPost,
        commitActionName: 'SAVE',
        cancelAction: this.toggleEditPane,
        cancelActionName: 'CANCEL',
        autofocus: false,
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
          parentFp: this.post.parent,
          contentType: 'report',
        },
        commitAction: this.submitReport,
        commitActionName: 'REPORT',
        cancelAction: this.toggleReportPane,
        cancelActionName: 'CANCEL',
        autofocus: false,
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
          parentFp: this.post.parent,
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
          parentFp: this.post.parent,
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
    /*----------  Notifications highlight / focus stuff  ----------*/
    isNotificationsParent(this: any) {
      if (globalMethods.IsUndefined(this.post.fingerprint)) {
        return false
      }
      if (this.post.fingerprint === this.notificationparent) {
        return true
      }
      return false
    },
    isNotificationsHighlight(this: any) {
      if (globalMethods.IsUndefined(this.post.fingerprint)) {
        return false
      }
      for (let val of this.notificationhighlights) {
        if (this.post.fingerprint === val) {
          return true
        }
      }
      return false
    },
    isNotificationsFocus(this: any) {
      if (globalMethods.IsUndefined(this.post.fingerprint)) {
        return false
      }
      if (this.post.fingerprint === this.notificationfocus) {
        return true
      }
      return false
    },
    /*----------  END: Notifications highlight / focus stuff  ----------*/
    /*
        This requires a little bit of an explanation, since there is a few layers of stuff here.

        First of all, visibleInflightEntity, or visible[anything] is for edits. Effectively, 'post' is the unedited underlying post, and 'visiblePost' is the edit applied on it, if any. Why is this necessary? Because we want to show the user the edits that it has made immediately, alongside a progress bar for showing when the edit will actually commit in.

        This gets a little confusing to read at first because the post itself can also show recursive children, and each of those children also use the same template. But fundamentally it's simple: what we do is if it's an inflight entity, we disable the further content generation anchors (upvote, downvote, report, create, edit ...) until that entity is fully generated and is in place.

        One trick to keep in mind that I'm doing is I'm only replacing the post's edited body, and I keep all other data using the original post. This simplifies a lot of things because our inflight posts don't actually carry most of the stuff the post has, just the delta (so that it's more efficient), so showing the inflight post fully in exchange to the old post would leave a lot of fields missing.
      */
    postCreation(this: any) {
      // These are necessary because in uncompiled entities, these are in thread.provable.creation, but in compiled ones it's in thread.creation.
      if (this.uncompiled || this.isInflightEntity) {
        if (typeof this.post.provable === 'undefined') {
          return 0
        }
        return this.post.provable.creation
      }
      return this.post.creation
    },
    postLastUpdate(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        if (typeof this.post.updateable === 'undefined') {
          return 0
        }
        return this.post.updateable.lastupdate
      }
      return this.post.lastupdate
    },
    postOwner(this: any) {
      if (this.isInflightEntity) {
        if (
          this.$store.state.localUserArrived &&
          this.$store.state.localUserExists
        ) {
          return this.$store.state.localUser
        }
      }
      return this.post.owner
    },
    postFingerprint(this: any) {
      if (this.uncompiled) {
        return this.post.provable.fingerprint
      }
      return this.post.fingerprint
    },
    isVisibleInflightEntity(this: any): boolean {
      if (typeof this.visibleInflightStatus !== 'undefined') {
        return true
      }
      return false
    },
    inflightInProgress(this: any): boolean {
      if (typeof this.visibleInflightStatus === 'undefined') {
        return false
      }
      if (this.visibleInflightStatus.completionpercent === 100) {
        return false
      }
      return true
    },
    isInflightEntity(this: any): boolean {
      if (typeof this.inflightStatus !== 'undefined') {
        return true
      }
      return false
    },
    inflightChildren(this: any): any[] {
      let iflChildren: any[] = []
      for (let val of this.$store.state.ambientStatus.inflights.postsList) {
        if (this.post.fingerprint !== val.entity.parent) {
          continue
        }
        if (this.post.fingerprint === val.entity.provable.fingerprint) {
          continue
        }
        let isEditOfChild = false
        for (let child of this.post.childrenList) {
          if (child.fingerprint === val.entity.provable.fingerprint) {
            isEditOfChild = true
            break
          }
        }
        if (isEditOfChild) {
          continue
        }
        iflChildren.push(val)
      }
      return iflChildren
    },
    inflightEdits(this: any): any[] {
      /*This one looks for whether the current post itself was edited, and if so, lists all of them. We'll be using the most recent one off of this list.*/
      let iflEdits: any[] = []
      for (let val of this.$store.state.ambientStatus.inflights.postsList) {
        if (this.postFingerprint !== val.entity.provable.fingerprint) {
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
    visiblePost(this: any) {
      if (typeof this.mostRecentInflightEdit !== 'undefined') {
        return this.mostRecentInflightEdit.entity
      }
      return this.post
    },
    visibleInflightStatus(this: any) {
      if (typeof this.mostRecentInflightEdit !== 'undefined') {
        return this.mostRecentInflightEdit.status
      }
      return this.inflightStatus
    },
    /*----------  Visibility checks  ----------*/
    /*
        These are useful because our visible / nonvisible logic is getting too complex to retain in the template itself.
      */
    actionsVisible(this: any) {
      // If in mod actions view, always false
      if (this.isinmodactionsview) {
        return false
      }
      if (this.localUserReadOnly) {
        return false
      }
      if (this.uncompiled) {
        return false
      }
      if (this.isInflightEntity) {
        return false
      }
      return true
    },
    /*----------  Reports  ----------*/
    reportSubmittedBefore(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      if (this.post.compiledcontentsignals.selfreported) {
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
      for (let val of this.post.compiledcontentsignals.reportsList) {
        if (val.sourcefp === this.$store.state.localUser.fingerprint) {
          return val.reason
        }
      }
      return ''
    },
    reportsListVisible(this: any) {
      // If in mod activity pane, don't show reports list
      if (this.isinmodactionsview) {
        return false
      }
      // If uncompiled or inflight, always false
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      // If mod mode is enabled, true if reports are present
      if (
        this.$store.state.modModeEnabledArrived &&
        this.$store.state.modModeEnabled
      ) {
        if (this.post.compiledcontentsignals.reportsList.length > 0) {
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
        return this.post.compiledcontentsignals.reportsList
      }
      let nonModReports: any = []
      // only reports created by self
      if (globalMethods.IsUndefined(this.$store.state.localUser)) {
        return nonModReports
      }
      for (let val of this.post.compiledcontentsignals.reportsList) {
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
      if (this.post.compiledcontentsignals.selfmodblocked) {
        return true
      }
      return false
    },
    modApprovalSubmittedBefore(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      if (this.post.compiledcontentsignals.selfmodapproved) {
        return true
      }
      return false
    },
    modIgnoreSubmittedBefore(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      if (this.post.compiledcontentsignals.selfmodignored) {
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
      // If in mod actions view, always false
      if (this.isinmodactionsview) {
        return false
      }
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
    /*----------  By op  ----------*/
    byOp(this: any) {
      if (this.uncompiled || this.isInflightEntity) {
        return false
      }
      return this.post.compiledcontentsignals.byop
    },
    externalUrl(this: any) {
      return (
        'aether:/' +
        this.$store.state.route.path +
        '?focusSelector=' +
        this.post.fingerprint +
        '&highlightSelectors=["' +
        this.post.fingerprint +
        '"]'
      )
    },
    uncompiledDeepLink(this: any) {
      return (
        '/board/' +
        this.post.board +
        '/thread/' +
        this.post.thread +
        '?focusSelector=' +
        this.post.provable.fingerprint +
        '&highlightSelectors=["' +
        this.post.provable.fingerprint +
        '"]'
      )
    },
    compiledDeepLink(this: any) {
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
    /*----------  Is deleted  ----------*/
    isDeletedByCreator(this: any) {
      return this.visiblePost.body === '*[deleted]*'
    },
    postBodyPreview(this: any) {
      let trailer = this.post.body.length > 250 ? '...' : ''
      return this.post.body.substring(0, 250) + trailer
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
    /*----------  Reply actions  ----------*/
    toggleReplyPane(this: any) {
      console.log('toggle reply pane runs')
      if (this.replyPaneOpen) {
        this.replyPaneOpen = false
      } else {
        this.replyPaneOpen = true
        // this.editPaneOpen = false
      }
    },
    submitPost(this: any, fields: any) {
      let postBody = ''
      for (let val of fields) {
        if (val.id === 'postBody') {
          postBody = val.content
        }
      }
      let post = new mimobjs.Post()
      // Set board, thread, parent, body fields
      post.setBoard(this.post.board)
      post.setThread(this.post.thread)
      post.setParent(this.post.fingerprint)
      post.setBody(postBody)
      let vm = this
      fe.SendPostContent('', post, function (resp: any) {
        console.log(resp.toObject())
        vm.toggleReplyPane()
      })
    },
    /*----------  Edit actions  ----------*/
    toggleEditPane(this: any) {
      if (this.editPaneOpen) {
        this.editPaneOpen = false
      } else {
        // this.replyPaneOpen = false
        this.postEditExistingSpec.fields[0].content = this.visiblePost.body
        this.editPaneOpen = true
      }
    },
    /*----------  Delete actions  ----------*/
    toggleSelfDeleteConfirm(this: any) {
      this.selfDeleteConfirmOpen = !this.selfDeleteConfirmOpen
    },
    selfDelete(this: any) {
      let postBody = '*[deleted]*'
      let post = new mimobjs.Post()
      // Set board, thread, parent, body fields
      let pv = new mimobjs.Provable()
      pv.setFingerprint(this.post.fingerprint)
      post.setProvable(pv)
      post.setBoard(this.post.board)
      post.setThread(this.post.thread)
      post.setParent(this.post.parent) // < heads up, this is different
      post.setBody(postBody)
      this.selfDeleteConfirmOpen = false
      fe.SendPostContent(this.post.fingerprint, post, function (resp: any) {
        console.log(resp.toObject())
      })
    },
    submitEditExistingPost(this: any, fields: any) {
      let postBody = ''
      for (let val of fields) {
        if (val.id === 'postBody') {
          postBody = val.content
        }
      }
      let post = new mimobjs.Post()
      // Set board, thread, parent, body fields
      let pv = new mimobjs.Provable()
      pv.setFingerprint(this.post.fingerprint)
      post.setProvable(pv)
      post.setBoard(this.post.board)
      post.setThread(this.post.thread)
      post.setParent(this.post.parent) // < heads up, this is different
      post.setBody(postBody)
      console.log(this.post.fingerprint)
      let vm = this
      fe.SendPostContent(this.post.fingerprint, post, function (resp: any) {
        console.log(resp.toObject())
        vm.toggleEditPane()
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
        this.post.fingerprint,
        '',
        reportReason,
        this.post.board,
        this.post.thread,
        function (resp: any) {
          console.log(resp.toObject())
          vm.toggleReportPane()
          vm.reportSubmitted = true
        }
      )
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
        this.post.fingerprint,
        '',
        approvalReason,
        this.post.board,
        this.post.thread,
        function (resp: any) {
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
        this.post.fingerprint,
        '',
        deleteReason,
        this.post.board,
        this.post.thread,
        function (resp: any) {
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
        this.post.fingerprint,
        '',
        '',
        this.post.board,
        this.post.thread,
        function (resp: any) {
          console.log(resp.toObject())
          vm.modIgnoreSubmitted = true
        }
      )
    },
    /*----------  Refresh func inflight info  ----------*/
    refresh(this: any) {
      console.log('post refresher is called')
      vuexStore.dispatch('refreshCurrentThreadAndPosts', {
        boardfp: this.$store.state.route.params.boardfp,
        threadfp: this.$store.state.route.params.threadfp,
      })
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
      let vm = this
      setTimeout(function () {
        vm.linkCopied = false
        vm.copyLinkInProgress = false
      }, 1250)
    },
    /*----------  Go to entity (uncompiled)  ----------*/
    goTo(this: any) {
      if (this.uncompiled) {
        console.log(this.uncompiledDeepLink)
        this.$router.push(this.uncompiledDeepLink)
      } else {
        this.$router.push(this.compiledDeepLink)
      }
    },
  },
}
</script>

<style lang="scss">
.post.inflight {
  .markdowned p:last-child {
    margin-bottom: 0;
  }
}
</style>

<style lang="scss" scoped>
@import '../scss/globals';
.post {
  display: flex;
  flex-direction: column;
  padding-left: 20px;
  margin: 20px 20px;
  border-left: 3px solid rgba(255, 255, 255, 0.15);
  position: relative;

  &.inflight {
  }

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
  .ballot {
    padding-top: 0px;
  }
  &:hover {
    > .post-content > .ballot {
      visibility: visible;
    }
  }
}

.actions-rack {
  display: flex;
  .action {
    margin-right: 10px;
    font-family: 'SSP Semibold';
    font-size: 14.4px;
    line-height: 21.6px;
    user-select: none;
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

.notification-parent {
  // background-color: red;
  border-color: $a-cerulean;
}

.self-owned {
  // background-color: red;
  border-color: $a-cerulean;
}

.notification-highlight {
  border-color: $a-orange; // background-color: blue;
}
/*----------  Reports  ----------*/

.report-composer-container {
  padding-top: 30px;
}

.modtools {
  font-size: 16px;
  font-family: 'SSP Bold';
  .modtools-block {
    max-width: 600px;
    margin-top: 5px;
    &:first-of-type {
      margin-top: 20px;
    }
  }
}

.reports-list-container {
  // padding: 15px 20px 20px 20px;
  padding: 20px;
  background-color: rgba(0, 0, 0, 0.25); // margin-top: 20px;
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

.mod-actions-taken-container {
  max-width: 600px;
  .block {
    padding: 20px;
    background-color: rgba(0, 0, 0, 0.25); // margin-top: 20px;
    border-radius: 3px; // font-size: 110%;
    color: $a-grey-600;
  }
  .header-block {
    margin-top: 10px;
  }
  .iterator-block {
    margin-top: 3px;
    padding: 10px 20px;
  }
  .mod-actions-taken-iterator {
    border-left: 3px solid rgba(255, 255, 255, 0.15);
    padding-left: 20px;
    // margin-top: 15px;
    display: flex;
    .mod-action-taken-data {
      flex: 1;
    }
    .mod-action-taken-action {
      display: flex;
      a {
        pointer-events: none;
        user-select: none;
        margin: auto;
        margin-left: 10px;
      }
    }

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
  .mod-action-taken-type {
    display: flex;
    .mod-action-taken-type-text {
      flex: 1;
    }
  }
  .mod-actions-taken-header {
    font-family: 'SSP Bold';
    // margin-bottom: 15px;
    font-size: 110%;
    color: $a-grey-800;
  }
  .final-state-container {
    display: flex;
    .final-state-text {
      flex: 1;
    }
    .final-state-result {
      pointer-events: none;
      user-select: none;
      margin-left: 10px;
    }
  }
}

.mod-actions-container {
  padding: 20px;
  background-color: rgba(0, 0, 0, 0.25); // margin-top: 20px;
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
  line-height: 19px;
  font-size: 16px;
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
.self-delete-confirm-text {
  font-family: 'SSP Regular Italic';
  user-select: none;
  cursor: default;
  line-height: 23px;
  font-size: 16px;
  margin-right: 10px;
  padding: 10px 10px 10px 12px;
  background-color: rgba(0, 0, 0, 0.25);
  display: inline-block;
  margin-top: 10px;
  margin-bottom: 10px;
  border-radius: 2px;
  text-align: center;
  a {
    font-family: 'SSP Regular Italic';
    margin: 0 10px 0 5px;
    &:first-of-type {
      margin: 0 5px 0 10px;
    }
  }
}
</style>
