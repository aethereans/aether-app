<template>
  <div class="notfound-container">
    <div class="not-found">
      <div class="headline">
        <icon name="unlink" scale="4"></icon>
      </div>
      <div class="subheader">
        Entity not found
      </div>
      <div class="description">
        If you opened the app via clicking a link on the Internet, welcome back!
        <br /><b
          >Aether works best when it's kept in the tray or menu bar, rather than
          completely shut down.</b
        >
      </div>
      <div class="description">
        If that's the case, your app didn't have the chance to get fresh
        content. Keep it open and try again in a few minutes.
      </div>
      <div class="description">
        If that's not the case, and you see this often,
        <a href="https://meta.getaether.net/c/support">please reach out</a>,
        it's worth a look.
      </div>
      <div class="goback">
        <a class="button is-warning is-outlined" @click="goToPopular"
          >GO TO POPULAR</a
        >
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'a-notfound',
  data() {
    return {}
  },
  methods: {
    goToPopular(this: any) {
      this.$router.push('/popular')
    },
    send404(this: any) {
      const metrics = require('../services/metrics/metrics')()
      metrics.SendRaw('Target FP not found', {
        'A-NotFoundObjectFp': this.$store.state.route.path,
      })
    },
  },
  mounted(this: any) {
    this.send404()
  },
  updated(this: any) {
    this.send404()
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.notfound-container {
  width: 100%;
  display: flex;
  // padding-top: 100px;

  .not-found {
    margin: auto; // padding-bottom: 7%; // To compensate for the line-height of the top thing.
    font-family: 'SCP Semibold';
    color: $a-grey-600;
    text-align: center;
    .headline {
      font-size: 700%;
      line-height: 100%;
      font-family: 'SCP ExtraLight';
      padding-bottom: 15px;
    }
    .subheader {
      font-size: 125%;
      font-family: 'SCP Bold';
    }
    .description {
      font-family: 'SSP Regular';
      margin: 0 auto;
      margin-top: 20px;
      width: 70%;
      font-size: 110%;
      b {
        font-family: 'SSP Black';
      }
      a {
        font-family: 'SSP Bold';
      }
    }
    .goback {
      padding-top: 40px;
      font-family: 'SSP Bold';
    }
  }
}
</style>
