<template>
  <router-link
    class="profile-header"
    tag="div"
    :to="link"
    :class="{ clickable: clickable }"
  >
    <div class="profile-avatar">
      <div class="profile-img">
        <a-hashimage
          :hash="user.fingerprint"
          isUser="true"
          :height="hashImageHeight + 'px'"
        ></a-hashimage>
      </div>
    </div>
    <div class="profile-name">
      <div class="profile-name-text">
        <a-username :owner="this.user" :styling="'avatar'"></a-username>
        <!-- {{getUserName()}} -->
      </div>
    </div>
    <div class="profile-fingerprint" v-if="!nofingerprint">
      <div class="profile-fingerprint-text">
        <a-fingerprint
          :fp="user.fingerprint"
          :fontSize="11.5"
          :link="externalUrl"
        ></a-fingerprint>
      </div>
    </div>
  </router-link>
</template>

<script lang="ts">
var globalMethods = require('../services/globals/methods')
export default {
  name: 'a-avatar-block',
  props: ['user', 'clickable', 'nofingerprint', 'imageheight'],
  data() {
    return {}
  },
  methods: {
    getUserName(this: any): string {
      // console.log(this.user)
      return globalMethods.GetUserName(this.user)
    },
  },
  computed: {
    // userName(this: any) {
    //   return this.user.NonCanonicalName
    // },
    link(this: any): string {
      if (this.clickable) {
        return '/user/' + this.user.fingerprint
      }
      return ''
    },
    hashImageHeight(this: any) {
      if (typeof this.imageheight === 'undefined') {
        return 128
      }
      return this.imageheight
    },
    externalUrl(this: any) {
      return 'aether:/' + this.$store.state.route.path
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../scss/globals';
.profile-header {
  width: 320px;
  display: flex;
  flex-direction: column;
  padding: 20px 15px 12px 15px;
  &.clickable {
    cursor: pointer;
  }
  .profile-avatar {
    display: flex;

    .profile-img {
      // width: 128px;
      // height: 128px;
      margin: auto;
      .hashimage {
        margin: auto;
      }
    }
  }
  .profile-fingerprint {
    display: flex;

    .profile-fingerprint-text {
      margin: auto;
      margin-top: 15px;
    }
  }
  .profile-name {
    display: flex;
    margin-top: 20px;
    .profile-name-text {
      font-size: 125%;
      margin: auto;
    }
  }
}
</style>
