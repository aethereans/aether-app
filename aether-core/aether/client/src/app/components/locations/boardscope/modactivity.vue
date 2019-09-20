<template>
  <div class="board-sublocation">
    <div class="board-modactivity">
      <template v-if="!$store.state.currentBoardsModActionsArrived">
        <div class="spinner-container">
          <a-spinner :delay="300"></a-spinner>
        </div>
      </template>
      <template v-else>
        <!-- ModActions arrived -->
        <template v-if="$store.state.currentBoardsModActions.length === 0">
          <!-- No content available -->
          <a-no-content
            no-content-text="There are no mod actions in recorded history."
            quoteDisabled="true"
          ></a-no-content>
        </template>
        <template v-else>
          <!-- Some modactions available -->
          <template v-for="modaction in $store.state.currentBoardsModActions">
            <!-- It's a thread -->
            <template v-if="modaction.threadpayload.fingerprint.length > 0">
              <a-thread-header-entity
                :thread="modaction.threadpayload"
                :isinmodactionsview="true"
                :useBodyPreview="true"
              ></a-thread-header-entity>
              <div class="divider"></div>
            </template>
            <!-- Or a post -->
            <template v-if="modaction.postpayload.fingerprint.length > 0">
              <a-post
                :isinmodactionsview="true"
                :post="modaction.postpayload"
                :useBodyPreview="true"
              ></a-post>
              <div class="divider"></div>
            </template>
          </template>
          <a-fin-puck></a-fin-puck>
        </template>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'modactivity',
  data() {
    return {}
  },
  methods: {
    actionTaken(this: any, ccs: any) {
      if (ccs.selfmodblocked || ccs.modblocked) {
        return true
      }
      return false
    },
  },
}
</script>

<style lang="scss" scoped>
.board-sublocation .board-modactivity {
  .divider {
    width: 100%;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    &:nth-last-of-type(2) {
      // ^ To accommodate the fin puck
      border-bottom: none;
    }
  }
  .spinner-container {
    display: flex;
    .spinner {
      margin: auto;
      padding-top: 50px;
    }
  }
}
</style>
