<template>
  <div class="view-header">
    <div class="center-container">
      <div class="header-container side-layout">
        <template>
          <div class="header-text">Search</div>
          <div class="header-subtext">
            Search for communities, content or people
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-search-header',
  data() {
    return {}
  },
  computed: {},
  mounted() {},
  updated() {},
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.view-header {
  height: 175px;
  display: flex;
  border-bottom: 1px solid rgba(0, 0, 0, 0.15);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.25), 1px 0 0 0 rgba(0, 0, 0, 0.2) inset;
  background-color: #1a1f23; //$a-grey-200 * 0.5;
  border-left: 5px solid $a-turquoise;
  line-height: 200%;

  .center-container {
    flex: 1;
    display: flex;
    flex-direction: column;

    .header-container {
      margin: auto;
      text-align: center;
      margin-bottom: 32px;

      &.side-layout {
        .header-text,
        .header-subtext {
          text-align: left;
          padding: 0;
          line-height: 140%;
          padding-top: 15px;
        }

        .header-subtext {
          font-family: 'SSP Regular';
          padding-top: 0;
          &.multiline {
            padding-top: 15px;
          }
        }
        margin-left: 50px;
      }

      .header-text {
        font-size: 200%;
      }

      .header-subtext {
        padding: 0px 50px;
      }
    }
  }
}
</style>
