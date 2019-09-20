<template>
  <div class="inflight-info">
    <a-progress-bar
      :percent="inflightStatus.completionpercent"
    ></a-progress-bar>
    <div class="text-info">
      <div class="current-status">
        <div class="current-status-text">
          {{ inflightStatus.statustext }}
          <template
            v-if="
              inflightStatus.completionpercent === 100 &&
                inflightStatus.eventtype != 'UPDATE'
            "
          >
            You can
            <a class="refresher-link" @click.prevent="refresherFunc">refresh</a>
            to interact with it.
          </template>
        </div>
        <div class="flex-spacer"></div>
        <div class="more-info">
          <a-info-marker
            header="Your submission is currently being minted."
            text="<p>Every new change you make requires a proof of work to be minted to prove the other nodes that you're a real user. This can take a couple seconds to a couple minutes depending on your preferred proof-of-work strength. </p><p>If these are taking too long, you can reduce the strength from settings. (Be mindful that dropping it too low can cause them to reject your messages.)</p>"
          ></a-info-marker>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-inflight-info',
  props: ['status', 'refresherFunc'],
  data() {
    return {}
  },
  computed: {
    inflightStatus(this: any) {
      return this.status
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.inflight-info {
  // width: 250px; // padding: 10px;
  // padding-top: 15px;
  background-color: rgba(
    255,
    255,
    255,
    0.03
  ); // border: 1px solid rgba(0, 0, 0, 0.2);
  border-radius: 2px;
  margin-bottom: 5px;
  padding: 10px;
  padding-bottom: 5px;
  .text-info {
    padding-top: 4px;
    color: $a-grey-600;
    font-family: 'SSP Regular';
  }
}

.current-status {
  display: flex;
  .current-status-text {
  }
  .flex-spacer {
    flex: 1;
  }
  .more-info {
    padding-top: 1px;
    padding-left: 5px;
  }
}

.refresher-link {
  font-family: 'SSP Semibold';
}
</style>
