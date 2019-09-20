<template>
  <div class="no-content">
    <div class="no-content-carrier">
      <template v-if="noContentText">
        <div class="no-content-container">
          <div v-html="noContentText" class="no-content-text"></div>
        </div>
      </template>
      <template v-if="!quoteDisabled">
        <div class="quote-container">
          <div class="quote-text" v-html="chosenQuote"></div>
          — <span class="quote-author"> {{ chosenQuoteAuthor }}</span>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
var GetQuote = require('../services/quotepicker/quotepicker')
export default {
  name: 'a-no-content',
  props: ['noContentText', 'quoteDisabled'],
  data() {
    return {
      chosenQuote: '',
      chosenQuoteAuthor: '',
    }
  },
  methods: {},
  beforeMount(this: any) {
    if (!this.quoteDisabled) {
      let q = GetQuote()
      this.chosenQuote = q.Quote
      this.chosenQuoteAuthor = q.Author
    }
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.no-content {
  font-family: 'SCP Regular';
  display: flex;
  border-radius: 3px;
  margin: 0 20px;
  margin-top: 20px; // background-color: rgba(255, 255, 255, 0.05);
  background-color: rgba(0, 0, 0, 0.25);
  .no-content-carrier {
    // font-family: "SSP Regular";
    padding: 15px 20px; // margin: auto;
    // width: 50%;
    max-width: 750px;
  }
  .no-content-container {
    padding-bottom: 25px;
    font-family: 'SCP Semibold';
    color: $a-grey-800;
  }
  .quote-container {
    color: $a-grey-400;
    .quote-text {
      // font-family: "SSP Regular Italic";
    }
    .quote-author {
      // font-family: "SSP Semibold Italic";
    }
  }
}
</style>
