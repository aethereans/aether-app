<template>
  <div
    @click="cryptoFundLinkClicked"
    class="button is-success crypto-fund-button"
  >
    FUND WITH CRYPTO
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-crypto-fund-button',
  data() {
    return {}
  },
  methods: {
    cryptoFundLinkClicked(this: any) {
      const metrics = require('../services/metrics/metrics')()
      metrics.SendRaw('Cryptofund link clicked')
      window.location.href =
        'https://commerce.coinbase.com/checkout/ae9092da-cb97-4e3f-bc98-b58d5c3bc7c8'
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.crypto-fund-button {
  font-family: 'SCP Bold';
  height: 41px;
  width: 175px;
  font-size: 14px;
  letter-spacing: 1px;
}
</style>
