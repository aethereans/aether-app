<template>
  <div class="composer">
    <template v-for="(block, index) in baseData.fields">
      <div class="block-meta-container">
        <div
          class="visible-name"
          v-show="block.visibleName"
          v-html="block.visibleName"
        ></div>
        <div
          class="description"
          v-show="block.description"
          v-html="block.description"
        ></div>
      </div>
      <!-- Search block (can only contain one search field, nothing else) -->
      <template v-if="block.searchField">
        <div class="composer-block search">
          <div class="single-row-text-container">
            <textarea
              :rows="block.heightRows"
              v-show="!block._previewIsVisible"
              :ref="'textarea_' + index"
              :placeholder="block.placeholder"
              class="composer-text-entry"
              v-model.trim="block.content"
              @keydown="textareaKeydown(index, $event)"
              @input="textareaInput(index)"
            ></textarea>
          </div>
          <div class="search-actions">
            <div class="flex-spacer"></div>
            <div class="actions">
              <a
                class="button is-outlined is-white"
                v-show="baseData.cancelActionName"
                @click="maybeCancel"
                >{{ baseData.cancelActionName }}</a
              >
              <a
                class="button is-outlined is-info"
                @click="maybeCommit"
                :disabled="!formIsValid"
                >{{ baseData.commitActionName }}</a
              >
            </div>
          </div>
        </div>
      </template>
      <!-- A non-search block -->
      <template v-else>
        <div
          class="composer-block"
          :class="{
            'preview-is-visible': block._previewIsVisible,
            'single-row': block.heightRows === 1,
          }"
        >
          <div class="composer-top-rack">
            <div class="errors-container" v-show="!block._previewIsVisible">
              <div
                class="error"
                v-show="
                  block._touched &&
                  block.content.length === 0 &&
                  !block.optional &&
                  !block.emptyWarningDisabled
                "
              >
                <div class="error-text">You need something in here</div>
              </div>
              <div
                class="error"
                v-show="
                  block._touched && block.maxCharCount < block.content.length
                "
              >
                <div class="error-text">This is way too long</div>
              </div>
              <div
                class="error"
                v-show="
                  block._touched &&
                  block.minCharCount > block.content.length &&
                  block.content.length !== 0
                "
              >
                <div class="error-text">This is too short</div>
              </div>
            </div>
            <div class="flex-spacer"></div>

            <div class="entry-actions">
              <div class="action-tag">
                <div
                  class="preview-tag"
                  @click="togglePreview(block)"
                  v-show="block._previewIsVisible"
                >
                  <a
                    href="https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet"
                    >MARKDOWN</a
                  >
                  PREVIEW
                </div>
              </div>
              <div
                class="preview-button-in-composer"
                hasTooltip
                title="Markdown preview"
                :class="{ enabled: block._previewIsVisible }"
                v-if="!block.previewDisabled"
                @click="togglePreview(block)"
              >
                <icon name="eye"></icon>
              </div>
            </div>
            <div
              class="info"
              :class="{
                toolong: block.maxCharCount - block.content.length < 0,
              }"
            >
              {{ block.maxCharCount - block.content.length }}
            </div>
          </div>

          <!-- Multiline template -->
          <template v-if="block.heightRows !== 1">
            <textarea
              v-show="!block._previewIsVisible"
              :ref="'textarea_' + index"
              class="composer-text-entry"
              v-model.trim="block.content"
              @keydown="textareaKeydown(index, $event)"
              :rows="block.heightRows"
              :placeholder="block.placeholder"
            ></textarea>
            <a-markdown
              class="composer-markdown-block"
              v-if="block._previewIsVisible"
              :content="block.content"
              :style="{ 'min-height': getMdHeight(index) }"
            ></a-markdown>
          </template>
          <!-- Multiline template end -->

          <!-- Single line template -->
          <!-- Single line password template -->
          <template v-if="block.heightRows === 1 && block.passwordField">
            <div class="single-row-text-container">
              <input
                type="password"
                v-show="!block._previewIsVisible"
                :ref="'textarea_' + index"
                class="composer-text-entry"
                v-model.trim="block.content"
                @keydown="textareaKeydown(index, $event)"
                @input="textareaInput(index)"
              />
            </div>
          </template>

          <template v-if="block.heightRows === 1 && !block.passwordField">
            <div class="single-row-text-container">
              <div class="pre-info-container" v-show="block.preInfoEnabled">
                <div class="text">
                  {{ block.preInfoText }}
                </div>
              </div>
              <textarea
                v-show="!block._previewIsVisible"
                :ref="'textarea_' + index"
                class="composer-text-entry"
                v-model.trim="block.content"
                @keydown="textareaKeydown(index, $event)"
                @input="textareaInput(index)"
                :rows="block.heightRows"
                :placeholder="block.placeholder"
              ></textarea>
              <div class="info-container">
                <div class="errors-container" v-show="!block._previewIsVisible">
                  <div
                    class="error"
                    v-show="
                      block._touched &&
                      block.content.length === 0 &&
                      !block.optional &&
                      !block.emptyWarningDisabled
                    "
                  >
                    <!-- Empty warning can be disabled specifically. For example, when you're putting this box on a persistent place that is always visible, you don't want to keep showing an annoying warning if somebody writes something and then decides to not post it and delete it. In text fields that can be closed (like reply to post) that is not a problem, but in ones that cannot be hidden without navigating away, the empty warning is annoying, and disabling the submit button communicates the same. -->
                    <div class="error-text">
                      You probably need something in here
                    </div>
                  </div>
                  <div
                    class="error"
                    v-show="
                      block._touched &&
                      block.maxCharCount < block.content.length
                    "
                  >
                    <div class="error-text">This is way too long</div>
                  </div>
                  <div
                    class="error"
                    v-show="
                      block._touched &&
                      block.minCharCount > block.content.length &&
                      block.content.length !== 0
                    "
                  >
                    <div class="error-text">This is too short</div>
                  </div>
                </div>
                <div
                  class="info-text"
                  :class="{
                    toolong: block.maxCharCount - block.content.length < 0,
                  }"
                >
                  {{ block.maxCharCount - block.content.length }}
                </div>
              </div>
            </div>
          </template>
          <!-- Single line template end -->

          <div
            class="composer-bottom-rack"
            v-if="!(baseData.fields.length > 1)"
          >
            <div class="flex-spacer"></div>
            <div class="actions">
              <a
                class="button is-outlined is-small is-white"
                v-show="baseData.cancelActionName"
                @click="maybeCancel"
                >{{ baseData.cancelActionName }}</a
              >
              <a
                class="button is-outlined is-info is-small"
                @click="maybeCommit"
                :disabled="!formIsValid"
                >{{ baseData.commitActionName }}</a
              >
            </div>
          </div>
        </div>
        <!-- Visible only if there are multiple form fields, otherwise  -->
      </template>
    </template>
    <div class="actions-rack" v-if="baseData.fields.length > 1">
      <div class="errors-container" v-if="validationErrorVisible">
        <div class="error">
          <div class="error-text">
            {{ baseData.validationError.Error }}
          </div>
        </div>
      </div>
      <div class="flex-spacer"></div>
      <div class="actions">
        <a
          class="button is-outlined is-white"
          v-show="baseData.cancelActionName"
          @click="maybeCancel"
          >{{ baseData.cancelActionName }}</a
        >
        <a
          class="button is-outlined is-info"
          :disabled="!formIsValid"
          @click="maybeCommit"
          >{{ baseData.commitActionName }}</a
        >
      </div>
    </div>
    <template v-if="baseData.fixToBottom">
      <div class="scrolltarget"></div>
    </template>
  </div>
</template>

<script lang="ts">
const isDev = require('electron-is-dev')
if (isDev) {
  var Vue = require('../../../node_modules/vue/dist/vue.js') // Production
} else {
  var Vue = require('../../../node_modules/vue/dist/vue.min.js') // Production
}
var Autosize = require('../../../node_modules/autosize')
var Mousetrap = require('../../../node_modules/mousetrap')
var Tooltips = require('../services/tooltips/tooltips')
export default {
  name: 'a-composer',
  props: ['spec'],
  data() {
    return {
      baseData: {},
      touchingBottom: false,
      scrollTarget: {},
      fullView: {},
      fullViewHeight: 0,
    }
  },
  created(this: any) {
    // let intermediateObj: any = {}
    // Copy it over to data
    // intermediateObj = this.spec
    for (let i = 0; i < this.spec.fields.length; i++) {
      // This needs to be set via Vue.set because that's the only way to add something after the fact and still get reactivity.
      Vue.set(this.spec.fields[i], '_touched', false)
      Vue.set(this.spec.fields[i], '_currentHeight', 0)
      Vue.set(this.spec.fields[i], '_previewIsVisible', false)
    }
    this.baseData = this.spec
  },
  computed: {
    formIsValid(this: any) {
      for (let val of this.baseData.fields) {
        if (!val.optional) {
          if (val.content.length === 0) {
            return false
          }
          if (val.content.length > val.maxCharCount) {
            return false
          }
          if (val.content.length < val.minCharCount) {
            return false
          }
        }
      }
      return true
    },
    validationErrorVisible(this: any) {
      return (
        typeof this.baseData.validationError !== 'undefined' &&
        this.baseData.validationError.Error.length > 0
      )
    },
  },
  mounted(this: any) {
    // console.log(this.baseData.content)
    this.startAutosize(this)
    Tooltips.Mount()
    // This is specific to new thread, and it allows the screen to be fixed to the bottom as the user types more and more into the thread body. In absence of this, the bottom of the page goes out of the screen and only scrolls down when the text itself crosses into the unseen zone. This instead keeps the bottom of the page fixed to the bottom, so long as it has touched bottom once before.
    if (this.baseData.fixToBottom) {
      let vm = this
      var options = {
        root: document.getElementsByClassName('main-block')[0],
        rootMargin: '0px',
        threshold: 0,
      }
      var observer: any = new IntersectionObserver(function (entries: any) {
        vm.touchingBottom = entries[0].isIntersecting
      }, options)
      this.scrollTarget = document.getElementsByClassName('scrolltarget')[0]
      this.fullView = document.getElementById('scrolltarget-container-target')
      this.fullViewHeight = this.fullView.getBoundingClientRect().height
      observer.observe(this.scrollTarget)
    }
    if (this.baseData.autofocus) {
      // Select first text area and bring it into focus
      document.getElementsByTagName('textarea')[0].focus()
      // ^ Heads up, this only works when the web inspector is closed.
    }
    /*----------  Mousetrap bindings  ----------*/
    let vm = this
    Mousetrap(vm.$el).bind('mod+enter', function () {
      vm.maybeCommit()
    })
    Mousetrap(vm.$el).bind('mod+esc', function () {
      vm.maybeCancel()
    })
  },
  updated(this: any) {
    // Tooltips.Mount()
  },
  beforeDestroy(this: any) {
    this.saveDraft()
  },
  beforeMount(this: any) {
    this.insertDraftsToFields()
  },
  methods: {
    textareaKeydown(this: any, index: any, event: any) {
      if (
        event.key !== 'Control' &&
        event.key !== 'Alt' &&
        event.key !== 'Meta' &&
        event.key !== 'Shift' &&
        event.key !== 'Tab' &&
        event.key !== 'CapsLock' &&
        event.key !== 'Backspace' &&
        event.key !== 'Super' &&
        event.key !== 'Hyper'
      ) {
        //event.key !== 'Monster')
        /*
          No joke, apparently Hyper key is a thing. No one has it on their keyboards though, so it's probably used to invoke Cthulhu.
        */
        this.markTouched(index)
        this.startAutosize(index)
        this.enforceNoSpaces(index, event)
        if (this.baseData.fixToBottom) {
          if (this.touchingBottom) {
            let vm = this
            setTimeout(function () {
              if (
                vm.fullView.getBoundingClientRect().height > vm.fullViewHeight
              ) {
                vm.scrollTarget.scrollIntoViewIfNeeded()
                vm.fullViewHeight = vm.fullView.getBoundingClientRect().height
              } else {
                vm.fullViewHeight = vm.fullView.getBoundingClientRect().height
              }
            }, 0)
          }
        }
      }
    },
    /*
        Input function runs on input - this scans for space if space is not allowed, and removes it.
      */
    textareaInput(this: any, index: any) {
      if (this.baseData.fields[index].spaceDisabled) {
        this.baseData.fields[index].content = this.baseData.fields[
          index
        ].content.replace(/\s/g, '')
      }
    },
    enforceNoSpaces(this: any, index: any, event: any) {
      if (this.baseData.fields[index].spaceDisabled) {
        if (event.keyCode === 32) {
          event.preventDefault()
        }
      }
    },
    startAutosize(this: any) {
      // This runs autosize on all boxes of this form at every keypress. This shouldn't cause much pain, since I think any run but the first is idempotent. The cool thing is that it does no longer need the index of the thing, so that this can become a thing that you can stick to mount or update events.
      if (this.baseData.sizing === 'fixed') {
        return
      }
      for (let i = 0; i < Object.keys(this.$refs).length; i++) {
        if (Object.keys(this.$refs)[i].indexOf('textarea_') !== -1) {
          Autosize(this.$refs[Object.keys(this.$refs)[i]][0])
        }
      }
    },
    // markTouched marks this specific base data as user _touched, so that we can actually show warnings.
    markTouched(this: any, index: any) {
      if (!this.baseData.fields[index]._touched) {
        this.baseData.fields[index]._touched = true
      }
    },
    togglePreview(this: any, block: any) {
      if (
        block._previewIsVisible ||
        typeof block._previewIsVisible === 'undefined'
      ) {
        block._previewIsVisible = false
      } else {
        block._previewIsVisible = true
      }
    },
    getMdHeight(this: any, index: any): string {
      /*
          There is something weird happening: when you enable preview on one of the panes, and go on to enable preview on a second pane, the first pane's getmdheight function is being called again. And it being zero (since the text field is hidden), it resets its min height to zero. I'm not sure why that is being called, because it doesn't share any state, there should not be anything that is able to cross-trigger it across two separate instances of <a-markdown>.

          Nevertheless, I made it so that the function returns its previous saved value if the new value is found to be zero, so that when you open a new preview, the first preview does not collapse into oblivion if it's empty.
        */
      let height = this.$refs['textarea_' + index][0].clientHeight
      if (height !== 0) {
        this.baseData.fields[index]._currentHeight = height
      }
      return this.baseData.fields[index]._currentHeight + 'px'
    },
    // The point of having these maybeCommit and maybeCancel actions is that wrapping those calls allows us to do maintenance, like clearing drafts whenever we need.
    maybeCommit(this: any) {
      if (this.formIsValid) {
        this.baseData.commitAction(
          this.baseData.fields,
          this.baseData.validationError
        )
        // Silence warnings so after the delete, it won't end up in a warning state. This can happen if the action is delayed to close the window but still has committed.
        if (!this.baseData.preventClearAfterSuccessfulCommit) {
          for (let val of this.baseData.fields) {
            val.content = ''
            val._touched = false
          }
        }
        this.clearDrafts()
      }
    },
    maybeCancel(this: any) {
      this.baseData.cancelAction()
      if (!this.baseData.cancelRetainsDraft) {
        for (let val of this.baseData.fields) {
          val.content = ''
          val._touched = false
        }
        this.clearDrafts()
      }
    },
    /*----------  Drafts  ----------*/
    saveDraft(this: any) {
      // If basedata or drafts inside it doesn't exist, return
      if (
        typeof this.baseData === 'undefined' ||
        typeof this.baseData.draft == 'undefined'
      ) {
        return
      }
      let dr: any = {}
      // Add fields to the fields array
      dr.fields = new Map()
      for (let field of this.baseData.fields) {
        dr.fields.set(field.id, field.content)
      }
      dr.parentFp = this.baseData.draft.parentFp
      dr.contentType = this.baseData.draft.contentType
      this.$store.dispatch('saveDraft', dr)
    },
    insertDraftsToFields(this: any) {
      // If basedata or drafts inside it doesn't exist, return
      if (
        typeof this.baseData === 'undefined' ||
        typeof this.baseData.draft == 'undefined'
      ) {
        return
      }
      // If drafts have no draft saved for this parentFp, return
      if (!this.$store.state.drafts.has(this.baseData.draft.parentFp)) {
        return
      }
      // We have a parent object draft.
      let dr = this.$store.state.drafts.get(this.baseData.draft.parentFp)
      if (!dr.has(this.baseData.draft.contentType)) {
        return
      }
      dr = dr.get(this.baseData.draft.contentType)
      for (let field of this.baseData.fields) {
        if (!dr.has(field.id)) {
          continue
        }
        field.content = dr.get(field.id)
      }
    },
    // Mostly the same as above, just in reverse - we're deleting keys we no longer need.
    clearDrafts(this: any) {
      // If basedata or drafts inside it doesn't exist, return
      if (
        typeof this.baseData === 'undefined' ||
        typeof this.baseData.draft == 'undefined'
      ) {
        return
      }
      // If drafts have no draft saved for this parentFp, return
      if (!this.$store.state.drafts.has(this.baseData.draft.parentFp)) {
        return
      }
      // We have a parent object draft.
      let dr = this.$store.state.drafts.get(this.baseData.draft.parentFp)
      if (!dr.has(this.baseData.draft.contentType)) {
        return
      }
      dr = dr.get(this.baseData.draft.contentType)
      for (let field of this.baseData.fields) {
        if (!dr.has(field.id)) {
          continue
        }
        dr.delete(field.id)
      }
    },
  },
}
</script>

<!--
This is the full spec of a composer item. Keep this updates as you add new fields.

postComposer: {
  fields: [{
    id: "postBody",
    emptyWarningDisabled: true,
    visibleName: "",
    description: "",
    placeholder: "",
    maxCharCount: 20480,
    minCharCount: 2,
    heightRows: 5,
    previewDisabled: false,
    content: "",
    optional: false,
    spaceDisabled: false,
    preInfoEnabled:false,
    preInfoText: @,
  }],
  draft: {
    parentFp: this.post.fingerprint,
    contentType: "post"
  },
  commitAction: function() {},
  commitActionName: "SUBMIT",
  cancelAction: function() {},
  cancelActionName: "",
  autofocus:false,
  preventClearAfterSuccessfulCommit:false,
  cancelRetainsDraft:false,
},

 -->

<style lang="scss" scoped>
@import '../scss/globals';
.composer {
  flex: 1; // display: flex;
  // flex-direction: column;
  max-width: 100%;
}

.block-meta-container {
  margin-bottom: 10px;
  .visible-name {
    font-size: 110%;
    color: $a-grey-800;
    font-family: 'SSP Bold';
  }

  .description {
    font-family: 'SSP Regular';
  }
}

.composer-block {
  // flex: 1;
  padding: 10px 10px;
  padding-top: 5px;
  width: fit-content;
  background-color: $dark-base;
  border-radius: 5px;
  border: 1px solid $a-grey-300;
  width: 100%;
  margin-bottom: 15px;
  display: flex;
  flex-direction: column;

  &:last-of-type {
    margin-bottom: 0;
  }

  &.single-row {
    padding-top: 10px;
    .composer-top-rack {
      display: none;
    }
    .single-row-text-container {
      display: flex;
      .info-container {
        display: flex;
        flex: 1;
        .info-text {
          margin: auto;
          color: $a-grey-400;
          font-family: 'SCP Regular';
          cursor: default;
          padding-right: 4px;
          &.toolong {
            @extend .toolong;
          }
        }
      }
      .pre-info-container {
        display: flex;
        .text {
          margin: auto;
          position: relative;
          left: 8px;
          bottom: 1px;
          font-family: 'SSP Black';
          color: $a-grey-600;
        }
      }
    }
  }

  &:focus-within {
    // cool! css ftw
    outline: none;
    border: 1px solid $a-grey-800;
  }

  &.preview-is-visible {
    background-color: #1f313c; // $dark-base * 1.15;
  }
  .markdowned {
    font-family: 'SSP Regular';
    font-size: 110%; // margin-top: 3px;
  }
  .composer-text-entry {
    background-color: $a-transparent;
    color: $a-grey-800; // font-size: 100%;
    font-size: 17.6px;
    font-family: 'SSP Regular';
    border: none;
    padding: 5px 0;
    max-width: 100%;
    width: 100%;
    caret-color: $a-cerulean;
    resize: none;
    min-height: 0;
    line-height: 1.5;
    padding-right: 10px;
    padding-left: 10px;
    overflow: hidden;
    &:focus {
      outline: none;
    }
    &::placeholder {
      color: $a-grey-400;
    }
  }
  .composer-top-rack {
    display: flex;
    * {
      user-select: none;
    }
    .entry-actions {
      padding-right: 6px;
      display: flex;
    }
    .info {
      color: $a-grey-400;
      &.toolong {
        @extend .toolong;
      }
      font-family: 'SCP Regular';
      cursor: default;
      padding-right: 4px;
    }
  }
  .composer-bottom-rack {
    display: flex;
    cursor: default;
    * {
      user-select: none;
    }
    .actions {
      font-family: 'SSP Bold';
    }
  }

  &.search {
    padding-top: 15px;
    display: flex;
    flex-direction: row;
    .single-row-text-container {
      flex: 1;
      .composer-text-entry {
        font-family: 'SCP Regular';
        font-size: 140%;
      }
    }
    .search-actions {
      margin: 5px 10px;
    }
  }
}

.actions-rack {
  display: flex;
  .actions {
  }
  .errors-container {
    padding-left: 0;
    padding-top: 0;
    flex: 1;
    flex-grow: 10;
    .error {
      white-space: initial;
      height: auto;
      margin-left: 0;
      margin-right: 6px;
    }
  }
}

a.button {
  @extend %link-hover-ghost-extenders-disable;
}

.action-tag {
  display: flex;
  cursor: default;
  font-family: 'SSP Bold';
  .preview-tag {
    cursor: pointer;
    font-size: 80%;
    color: $a-grey-400;
    padding-right: 6px;
    margin: auto;
    margin-top: 3px;
    letter-spacing: 1px;
    a:hover {
      color: $a-grey-600;
    }
  }
}

.errors-container {
  padding-left: 10px;
  padding-top: 2px;
  display: flex;
  .error {
    margin: auto;
    border-radius: 2px;
    height: 22px;
    border: 1px solid $a-orange;
    white-space: nowrap;
    color: $a-orange;
    display: flex;
    padding: 1px 3px;
    margin-right: 6px;
    .error-text {
      line-height: 110%;
      margin: auto;
      font-family: 'SSP Regular'; // font-size: 90%;
    }
  }
}

.toolong {
  color: $a-orange;
  font-family: 'SCP Semibold';
}

.flex-spacer {
  flex: 1;
}
</style>
<style lang="scss">
@import '../scss/globals';
.preview-button-in-composer {
  margin: auto; // border: 1px solid $a-grey-600;
  border-radius: 50%;
  width: 18px;
  height: 18px;
  padding: 0 0 1px 2px;
  cursor: pointer;
  margin-top: 4px;
  &:hover {
    background-color: $a-grey-600;
    svg {
      fill: $a-grey-100;
    }
  }
  &:active {
    background-color: $a-grey-800;
    svg {
      fill: $a-grey-100;
    }
  }
  &.enabled {
    background-color: $a-turquoise;
    svg {
      fill: $a-grey-100;
    }
    &:hover {
      background-color: #72d7d7; // $a-turquoise * 1.2;
    }
    &:active {
      background-color: #85fbfb; // $a-turquoise * 1.4;
    }
  }
  svg {
    fill: $a-grey-400;
    position: relative;
    left: 0px;
    top: -2px;
    width: 14px;
    height: 14px;
  }
}

.composer-markdown-block {
  padding: 5px 0;
  overflow: hidden;
  padding-right: 10px;
  padding-left: 10px;
  p:last-of-type {
    margin-bottom: 0;
  }
}

.composer .description {
  b {
    font-family: 'SSP Bold';
  }
  i {
    font-family: 'SSP Regular Italic';
  }
}
</style>
