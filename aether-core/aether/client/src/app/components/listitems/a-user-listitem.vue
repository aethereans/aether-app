<template>
  <router-link class="user-listitem" :to="linkToUser">
    <div class="flex-carrier image">
      <div class="image-container">
        <a-hashimage :hash="userFingerprint" height="72px"></a-hashimage>
      </div>
    </div>
    <div class="flex-carrier main">
      <div class="main-data-container">
        <div class="user-name">
          <a-username :owner="this.user"></a-username>
        </div>
        <div class="user-info">
          {{ user.info }}
        </div>
      </div>
    </div>
  </router-link>
</template>

<script lang="ts">
export default {
  name: 'a-user-listitem',
  props: ['user'],
  data() {
    return {}
  },
  methods: {},
  computed: {
    userFingerprint(this: any) {
      return this.user.fingerprint
    },
    userName(this: any) {
      return this.user.noncanonicalname
    },
    linkToUser(this: any) {
      return '/user/' + this.userFingerprint
    },
  },
}
</script>

<style lang="scss" scoped>
@import '../../scss/globals';
a:hover {
  color: inherit;
}

.user-listitem {
  display: block;
  @extend %link-hover-ghost-extenders-disable;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  padding: 15px 5px;
  margin: 0 20px; // cursor: default;
  color: $a-grey-800;
  word-break: break-word;
  &:hover {
    color: $a-grey-800;
    background-color: rgba(255, 255, 255, 0.05);
  }
  a {
    @extend %link-hover-ghost-extenders-disable;
    &:hover {
      color: inherit;
    }
  }
}

.flex-carrier {
  display: flex;
  &.signals {
  }
  &.main {
    flex: 1;
  }
}

.main-data-container {
  flex: 1;
  padding-right: 15px;
  .user-name {
    font-size: 120%;
    cursor: pointer;
    padding-bottom: 5px;
  }
  .user-info {
    font-family: 'SSP Regular';
    color: $a-grey-500;
    font-size: 95%;
    @include lineClamp(3);
  }
}

.image-container {
  margin: auto;
  padding-right: 15px;
}
</style>
