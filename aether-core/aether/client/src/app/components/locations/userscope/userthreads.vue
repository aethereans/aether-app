<template>
  <div class="user-sublocation">
    <div class="user-threads">
      <template v-if="loadingComplete">
        <template v-if="isSelf">
          <template v-for="thr in inflightCreates.slice().reverse()">
            <a-thread-header-entity
              :thread="thr.entity"
              :key="thr.entity.Fingerprint"
              :inflightStatus="thr.status"
              :uncompiled="true"
            ></a-thread-header-entity>
            <!-- Well, it is not uncompiled, it's inflight - but we're using uncompiled here as a way to denote that this is in this specific list, not that it is necessarily uncompiled, it's already inflight, which is a superset of uncompiled. -->
            <div class="divider"></div>
          </template>
        </template>
        <div v-for="thr in threadsList" :key="thr.Fingerprint">
          <a-thread-header-entity
            :thread="thr"
            :uncompiled="true"
          ></a-thread-header-entity>
          <div class="divider"></div>
        </div>
        <div class="load-more-carrier" v-show="loadMoreVisible">
          <a
            class="button is-warning is-outlined load-more-button"
            @click="loadMore"
          >
            LOAD MORE
          </a>
        </div>
        <a-no-content
          no-content-text="No threads created in retained history."
          quoteDisabled="true"
          v-if="threadsList.length === 0 && inflightCreates.length === 0"
        >
        </a-no-content>
        <a-fin-puck v-show="!loadMoreVisible"></a-fin-puck>
      </template>
      <template v-else>
        <div class="spinner-container">
          <a-spinner></a-spinner>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
var fe = require('../../../services/feapiconsumer/feapiconsumer')
var Vuex = require('../../../../../node_modules/vuex').default
var globalMethods = require('../../../services/globals/methods')
export default {
  name: 'userthreads',
  data() {
    return {
      threadsList: [],
      currentUserFp: '',
      loadingComplete: false,
      limit: 25,
      offset: 0,
      lastLoadSize: 0,
    }
  },
  computed: {
    ...Vuex.mapState(['currentUserEntity']),
    inflightCreates(this: any) {
      let inflightCreates = []
      for (let val of this.$store.state.ambientStatus.inflights.threadsList) {
        if (val.status.eventtype !== 'CREATE') {
          continue
        }
        inflightCreates.push(val)
      }
      return inflightCreates
    },
    isSelf(this: any) {
      if (
        globalMethods.IsUndefined(this.$store.state.currentUserEntity) ||
        globalMethods.IsUndefined(this.$store.state.localUser)
      ) {
        return false
      }
      if (
        this.$store.state.currentUserEntity.fingerprint !==
        this.$store.state.localUser.fingerprint
      ) {
        return false
      }
      return true
    },
    loadMoreVisible(this: any) {
      if (this.lastLoadSize < this.limit) {
        return false
      }
      return true
    },
  },
  // watch: {
  //   currentUserEntity(this: any, val: any) {
  //     if (this.currentUserFp === val.fingerprint) {
  //       return
  //     }
  //     if (typeof val === 'undefined') {
  //       return
  //     }
  //     this.fetchData(val.fingerprint)
  //   }
  // },
  methods: {
    fetchData(this: any, targetuserfp: string) {
      let vm = this
      fe.GetUncompiledEntityByKey(
        'Thread',
        targetuserfp,
        '',
        '',
        this.limit,
        this.offset,
        function(resp: any) {
          console.log(resp)
          vm.threadsList = resp.threadsList
          vm.currentUserFp = targetuserfp
          vm.loadingComplete = true
          vm.lastLoadSize = resp.postsList.length
        }
      )
    },
    loadMore(this: any) {
      this.offset = this.offset + this.limit
      this.fetchData(this.currentUserFp)
    },
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
  },
}
</script>

<style lang="scss" scoped>
.user-sublocation {
  .user-threads {
    .thread-entity.no-user-present {
      border-bottom: none;
      padding-bottom: 15px;
    }
  }
}

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
