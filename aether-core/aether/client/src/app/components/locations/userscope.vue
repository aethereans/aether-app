<template>
  <div class="location user" v-if="$store.state.currentUserLoadComplete">
    <template v-if="entityNotFound">
      <a-notfound></a-notfound>
    </template>
    <template v-else>
      <div class="user-side">
        <a-avatar-block
          :user="$store.state.currentUserEntity"
          :nofingerprint="true"
        ></a-avatar-block>
        <div class="user-info-flex-carrier">
          <div class="user-info-flex-child">
            <div class="user-info">
              <a-markdown
                :content="$store.state.currentUserEntity.info"
              ></a-markdown>
            </div>
          </div>
        </div>
      </div>

      <div class="user-main">
        <a-tabs :tabslist="tabslist"></a-tabs>
        <router-view></router-view>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
var globalMethods = require('../../services/globals/methods')
export default {
  name: 'userscope',
  data() {
    return {}
  },
  methods: {
    scrollSubviewToTop(this: any) {
      // This is necessary because the subview is somehow not responding to window.scrollTo. So we have to specifically target it. (I tried setting body/html to 100%, etc. Seems it's specific to this case.)
      let el = document.getElementsByClassName('main-block')[0]
      el.scrollTo(0, 0)
    },
  },
  computed: {
    // currentUser(this: any) {
    //   return this.$store.state.currentUserEntity
    // },
    entityNotFound(this: any) {
      return this.$store.state.currentUserEntity.fingerprint.length === 0
    },
    tabslist(this: any) {
      let selfUserTabsList = [
        {
          name: 'BIO',
          link: '/user/' + this.$store.state.currentUserEntity.fingerprint,
        },
        {
          name: 'NOTIFICATIONS',
          link:
            '/user/' +
            this.$store.state.currentUserEntity.fingerprint +
            '/notifications',
        },
        {
          name: 'POSTS',
          link:
            '/user/' +
            this.$store.state.currentUserEntity.fingerprint +
            '/posts',
        },
        {
          name: 'THREADS',
          link:
            '/user/' +
            this.$store.state.currentUserEntity.fingerprint +
            '/threads',
        },
        {
          name: 'COMMUNITIES CREATED',
          link:
            '/user/' +
            this.$store.state.currentUserEntity.fingerprint +
            '/boards',
        },
      ]
      let nonSelfUserTabsList = [
        {
          name: 'BIO',
          link: '/user/' + this.$store.state.currentUserEntity.fingerprint,
        },
        {
          name: 'POSTS',
          link:
            '/user/' +
            this.$store.state.currentUserEntity.fingerprint +
            '/posts',
        },
        {
          name: 'THREADS',
          link:
            '/user/' +
            this.$store.state.currentUserEntity.fingerprint +
            '/threads',
        },
        {
          name: 'COMMUNITIES CREATED',
          link:
            '/user/' +
            this.$store.state.currentUserEntity.fingerprint +
            '/boards',
        },
      ]
      if (globalMethods.IsUndefined(this.$store.state.localUser)) {
        return nonSelfUserTabsList
      }
      if (
        this.$store.state.currentUserEntity.fingerprint !==
        this.$store.state.localUser.fingerprint
      ) {
        return nonSelfUserTabsList
      }
      return selfUserTabsList
    },
  },
  mounted(this: any) {
    this.scrollSubviewToTop()
  },
  updated(this: any) {
    this.scrollSubviewToTop()
  },
}
</script>

<style lang="scss" scoped>
@import '../../scss/globals';
$user-sidebar-width: 320px;
$user-sidebar-base: #25272a; // $mid-base * 0.9;
.location {
  display: flex;
  flex-direction: row;
  padding-bottom: 0;
  flex: none;
  min-height: 100%;
}

.user-side {
  width: $user-sidebar-width;
  background-color: $user-sidebar-base; // background-color: $mid-base;
  position: fixed;
  z-index: 1;
  height: 100%;
  padding-top: 10px;
  display: flex;
  flex-direction: column;
  /*
    Why fixed instead of more modern flexbox? Because if we do a flex here, it creates a second scroll bar for that specific div instead of using the main one.
     */
  .user-info-flex-carrier {
    flex: 1;
    display: flex;
    .user-info-flex-child {
      margin: auto;
      margin-top: 0;
      max-width: 100%;
    }
  }
  .user-info {
    padding: 0 15px;
    font-family: 'SSP Regular';
    @include lineClamp(5);
    text-align: center;
    position: relative;
    .markdowned {
      margin: auto;
    }
    &:after {
      content: '';
      position: absolute;
      bottom: 0;
      right: 0;
      width: 70%;
      height: 1.3em;
      background: linear-gradient(
        to right,
        $a-transparent,
        $user-sidebar-base 75%
      );
    }
  }
}

.user-main {
  flex: 1;
  margin-left: $user-sidebar-width; // padding-left: 20px;
  border-left: 1px solid rgba(0, 0, 0, 0.25);
  border-radius: 5px 0 0 5px;
  overflow: hidden;
  padding-bottom: 50px;
}
</style>
