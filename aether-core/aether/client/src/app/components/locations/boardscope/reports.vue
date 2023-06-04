<template>
  <div class="board-sublocation">
    <div class="board-reports">
      <template v-if="!$store.state.currentBoardsReportsArrived">
        <div class="spinner-container">
          <a-spinner :delay="300"></a-spinner>
        </div>
      </template>
      <template v-else>
        <!-- Reports arrived -->
        <template v-if="reportsPaneEmpty">
          <!-- No content available -->
          <a-no-content
            no-content-text="Report inbox zero. Nice work! <br><br>There are no un-actioned reports in recorded history. <br><br>(Reports are only compiled when mod mode is enabled. If you just enabled it, you'll only see the reports from this point onwards)."
            quoteDisabled="true"
          ></a-no-content>
        </template>
        <template v-else>
          <!-- Some reports available -->
          <template v-for="report in $store.state.currentBoardsReports">
            <!-- It's a thread -->
            {{ report }}
            <template
              v-if="
                report.threadpayload.fingerprint.length > 0 &&
                !actionTaken(report.threadpayload.compiledcontentsignals)
              "
            >
              <a-thread-header-entity
                :thread="report.threadpayload"
                :isinreportsview="true"
              ></a-thread-header-entity>
              <div class="divider"></div>
            </template>
            <!-- Or a post -->
            <template
              v-if="
                report.postpayload.fingerprint.length > 0 &&
                !actionTaken(report.postpayload.compiledcontentsignals)
              "
            >
              <a-post
                :isinreportsview="true"
                :post="report.postpayload"
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
  name: 'reports',
  data() {
    return {}
  },
  methods: {
    reportsPaneEmpty(this: any) {
      for (let report of this.$store.state.currentBoardsReports) {
        if (
          report.threadpayload.fingerprint.length > 0 &&
          !this.actionTaken(report.threadpayload.compiledcontentsignals)
        ) {
          return false
        }
        if (
          report.postpayload.fingerprint.length > 0 &&
          !this.actionTaken(report.postpayload.compiledcontentsignals)
        ) {
          return false
        }
      }
      return true
    },
    actionTaken(this: any, ccs: any) {
      if (ccs.selfmodapproved || ccs.selfmodblocked) {
        return true
      }
      return false
    },
  },
}
</script>

<style lang="scss" scoped>
.board-sublocation .board-reports {
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
