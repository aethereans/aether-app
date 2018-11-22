<template>
  <div class="user-sublocation">
    <div class="user-posts">
      <template v-if="!loadingComplete">
        <div class="spinner-container">
          <a-spinner></a-spinner>
        </div>
      </template>
      <template v-else>
        <template v-if="isSelf">
          <template v-for="iflChild in inflightCreates.slice().reverse()">
            <a-post :post="iflChild.entity" :inflightStatus="iflChild.status"></a-post>
            <div class="divider"></div>
          </template>
        </template>
        <div v-for="post in postsList" :key="post.Fingerprint">
          <a-post :post="post" :uncompiled="true"></a-post>
          <div class="divider"></div>
        </div>
        <div class="load-more-carrier" v-show="loadMoreVisible">
          <a class="button is-warning is-outlined load-more-button" @click="loadMore">
        LOAD MORE
              </a>
        </div>
        <a-no-content no-content-text="No posts created in retained history." quoteDisabled="true" v-if="postsList.length === 0 && inflightCreates.length === 0">
        </a-no-content>
        <a-fin-puck v-show="!loadMoreVisible"></a-fin-puck>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
  var Vuex = require('../../../../../node_modules/vuex').default
  var fe = require('../../../services/feapiconsumer/feapiconsumer')
  var globalMethods = require('../../../services/globals/methods')
  export default {
    name: 'userposts',
    data() {
      return {
        postsList: [],
        currentUserFp: '',
        loadingComplete: false,
        limit: 25,
        offset: 0,
        lastLoadSize: 0,
      }
    },
    computed: {
      ...Vuex.mapState(['currentUserEntity']),
      /*----------  Inflight computeds  ----------*/
      inflightCreates(this: any) {
        let inflightCreates = []
        for (let val of this.$store.state.ambientStatus.inflights.postsList) {
          if (val.status.eventtype !== 'CREATE') {
            continue
          }
          inflightCreates.push(val)
        }
        return inflightCreates
      },
      isSelf(this: any) {
        if (globalMethods.IsUndefined(this.$store.state.currentUserEntity) || globalMethods.IsUndefined(this.$store.state.localUser)) {
          return false
        }
        if (this.$store.state.currentUserEntity.fingerprint !== this.$store.state.localUser.fingerprint) {
          return false
        }
        return true
      },
      loadMoreVisible(this: any) {
        if (this.lastLoadSize < this.limit) {
          return false
        }
        return true
      }
    },
    methods: {
      fetchData(this: any, targetuserfp: string) {
        let vm = this
        fe.GetUncompiledEntityByKey('Post', targetuserfp, this.limit, this.offset, function(resp: any) {
          console.log(resp)
          vm.postsList.push(...resp.postsList)
          vm.currentUserFp = targetuserfp
          vm.loadingComplete = true
          vm.lastLoadSize = resp.postsList.length
          console.log('post returned a response')
        })
      },
      loadMore(this: any) {
        this.offset = this.offset + this.limit
        this.fetchData(this.currentUserFp)
      }
    },
    beforeMount(this: any) {
      if (typeof this.currentUserEntity === 'undefined') {
        return
      }
      this.fetchData(this.currentUserEntity.fingerprint)
    },
    mounted(this: any) {},
    updated(this: any) {
      if (typeof this.currentUserEntity === 'undefined') {
        return
      }
      if (this.currentUserEntity.fingerprint === this.currentUserFp) {
        return
      }
      this.fetchData(this.currentUserEntity.fingerprint)
    }
  }
</script>

<style lang="scss">
  .user-sublocation .user-posts .post {
    .markdowned p:last-child {
      margin-bottom: 0;
    }
  }
</style>

<style lang="scss" scoped>
  @import "../../../scss/bulmastyles";
  .divider {
    width: 100%;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .spinner-container {
    display: flex;
    .spinner {
      margin: auto;
      padding-top: 50px;
    }
  }

  .load-more-carrier {
    display: flex;
    padding: 20px 0 0 0;
    .load-more-button {
      margin: auto;
    }
  }
</style>