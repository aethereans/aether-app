<template>
  <div class="location status">
    <div class="status-cards-container">
      <div class="card-block">
        <div class="status-card">
          <div class="card-header">
            <div class="card-header-text">Frontend</div>
            <div class="info-marker-container" v-show="loadComplete">
              <a-info-marker header="Front end is the part that makes this app meaningful to humans." text="<p>It generates a graph of entities and relationships from raw network data, and converts the content you create into raw blocks that other users can get from you, read, and interpret. </p><p>Since these users are also nodes, they will forward the content that you generated to the nodes that connect to them as well.</p>"></a-info-marker>
            </div>
            <div class="flex-spacer"></div>
            <a-guidelight :notifyState="dotStates.frontendDotState"></a-guidelight>
          </div>
          <div class="card-body">
            <div class="sub-block">
              <div class="sub-header">
                <div class="sub-header-text">Content
                </div>
                <div class="flex-spacer"></div>
                <a-guidelight :notifyState="dotStates.refresherDotState"></a-guidelight>
              </div>
              <div class="sub-body">
                <div class="row">
                  <div class="row-text">Last update</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderTimestamp(fas.lastrefreshtimestamp)}}</div>
                </div>
              </div>
            </div>

            <div class="sub-block">
              <div class="sub-header">
                <div class="sub-header-text">Graph compiler
                </div>

                <div class="info-marker-container" v-show="loadComplete">
                  <a-info-marker header="Graph compiler processes raw data into human readable relationships." text="<p>This part of Aether is tasked with converting raw Mim entities to what is visible as communities, threads, posts, signals such as upvotes, downvotes, and so on. </p><p>It receives the incremental updates that the backend has gathered from the network, and then inserts it into its existing network of meaningful relationships. </p><p>This component runs in intervals.</p>"></a-info-marker>
                </div>
                <div class="flex-spacer"></div>
                <a-guidelight :notifyState="dotStates.refresherDotState"></a-guidelight>
              </div>
              <div class="sub-body">
                <div class="row">
                  <div class="row-text">Status</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{fas.refresherstatus}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Last compile</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderTimestamp(fas.lastrefreshtimestamp)}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Last compile duration</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderDuration(fas.lastrefreshdurationseconds)}}</div>
                </div>
              </div>
            </div>

            <div class="sub-block">
              <div class="sub-header">
                <div class="sub-header-text">Entities in progress
                </div>

                <div class="info-marker-container" v-show="loadComplete">
                  <a-info-marker header="In-progress entities are those that you create or update while you're using the app." text="
                    <p>Here you can see the progress of content and signals you've recently created.</p>

                    <p>The creation and updates for these entities are not instant, because your device needs to expend a certain amount of effort to convince other nodes that it is acting in good faith. This is called a proof-of-work. </p>

                    <p>When an entity is completely minted and made available in the network, it will be removed from the inflights list.</p>

                    <p>All interactions with the network require varying strengths of proof-of-work. Your device will automatically mint the appropriate amount for you.</p>

                    <p> If this process seems to be taking a long time, you can reduce the proof-of-work requirement. Be mindful that dropping it too low can cause other nodes to reject your messages.<p>
                      ">
                  </a-info-marker>
                </div>
                <div class="flex-spacer"></div>
                <a-guidelight :notifyState="dotStates.inflightsDotState"></a-guidelight>
              </div>
              <table class="inflights-table">
                <thead>
                  <tr>
                    <th>Type</th>
                    <th>Event</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  <template v-if="mergedInflights.length > 0">
                    <template v-for="inflight in mergedInflights">
                      <tr>
                        <td>{{inflight.EntityType}}</td>
                        <td>{{inflight.ActionType}}</td>
                        <td>
                          <div class="progress-bar-container">
                            <a-progress-bar :percent="inflight.ProgressPercentage"></a-progress-bar>
                          </div>
                        </td>
                      </tr>
                    </template>
                  </template>
                  <template v-else>
                    <tr>
                      <td class="nothing-in-progress">Nothing is currently in progress.</td>
                      <td></td>
                      <td></td>
                    </tr>
                  </template>
                </tbody>
              </table>
            </div>

          </div>
        </div>
      </div>
      <div class="card-block">
        <div class="status-card">
          <div class="card-header">
            <div class="card-header-text">Backend</div>
            <div class="info-marker-container" v-show="loadComplete">
              <a-info-marker header="The back end is the node instance running on your machine that talks to other nodes, and to your front end." text="<p>It handles the distribution, validation, maintenance, and propagation of raw Mim content and signals.</p><p>There is a backend running on your device alongside a front end. This is the default way the client works, but it is not the only one. </p><p> For example, some users prefer to run their backends on a cloud service, and have the Aether app connect that backend. The benefit of doing this is that it makes your node always online, and always up to date, even when your computer is shut down. </p><p>It also comes with an additional benefit of being able to service multiple frontends from the same backend, so you can run a node for your family, or as a public service. Back ends are zero-knowledge, and they know nothing about the users present in front ends connected to them. </p>"></a-info-marker>
            </div>
            <div class="flex-spacer"></div>
            <a-guidelight :notifyState="dotStates.backendDotState"></a-guidelight>
          </div>
          <div class="card-body">

            <div class="sub-block">
              <div class="sub-header">
                <div class="sub-header-text">Network
                </div>
                <div class="info-marker-container" v-show="loadComplete">
                  <a-info-marker header="The content in this app comes from a network of nodes like you." text="<p>Your computer connects to these nodes in certain patterns to maximise its efficiency.</p><p>In their ideal state, all nodes in the network retain all of the information on the network, bracketed by time.</p><p> Network connections are pull-only. This means your node cannot push data out into another node, that node needs to ask first. As a result, it's important for other nodes to be able to connect to you, because if they cannot, content you create will not leave your device and it will not become available on the network.</p><p>Aether will attempt to reverse open connections in the case that no one connects to your computer for a while. In this case, your device will ask other people to connect to it, and hand over a live TCP socket to do so.</p><p>Outbound requests are triggered in intervals.</p>"></a-info-marker>
                </div>
                <div class="flex-spacer"></div>
                <a-guidelight :notifyState="dotStates.networkDotState"></a-guidelight>
              </div>
              <div class="sub-body">
                <div class="row">
                  <div class="row-text">Inbounds conns. in last 15m</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{bas.inboundscount15}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Last inbound</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderTimestamp(bas.lastinboundconntimestamp)}}</div>
                </div>
                <div class="spacer-row"></div>
                <div class="row">
                  <div class="row-text">Outbound conns. in last 15m</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{bas.outboundscount15}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Last outbound</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderTimestamp(bas.lastoutboundconntimestamp)}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Last outbound duration</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderDuration(bas.lastoutbounddurationseconds)}}</div>
                </div>
                <div class="spacer-row"></div>
                <!-- <div class="row">
                  <div class="row-text">Your node address</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{bas.localnodeexternalip}}</div>
                </div> -->
                <!-- Removed because it's not used anywhere, and since it's sourced from the router, it's often wrong. The user is better served finding own IP address through some other service, this can be misleading. -->
                <div class="row">
                  <div class="row-text">Your node port</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{bas.localnodeexternalport}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Port mapping status</div>
                  <div class="info-marker-container">
                    <a-info-marker header="Port mapping status indicates whether your router accepted a port open request from Aether." text="
                      <p>If you're seeing <em>mapping failed</em>, you might either have no router (e.g. tethering from your phone to connect to the internet), or you have one, and it refuses to map. Public routers (e.g. at Starbucks, at work) will usually refuse to map.</p>

                      <p>Even if the mapping fails, Aether can still work via a process called <em>reverse open</em>s. However, it will be slower to deliver the posts you create to other nodes. </p>

                      <p>If you're seeing this message at home, it might be worth trying to manually open the port from your router. If you do that, this entry will still say <em>failed</em>, but you should see an increase in your inbound connections. If you don't know how to do that, don't worry about it, it's not a big deal. </p>
                        ">
                    </a-info-marker>
                  </div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{bas.upnpstatus}} </div>
                </div>
              </div>
            </div>

            <div class="sub-block">
              <div class="sub-header">
                <div class="sub-header-text">Database
                </div>

                <div class="info-marker-container" v-show="loadComplete">
                  <a-info-marker header="Database is where the uncompiled Mim entities are stored for your frontend to consume." text="<p>Aether is an ephemeral network. By default it will delete data entities that haven't been referenced in 6 months.</p><p>That means unused and old data will be gone first, but all data on your database will eventually disappear to open space for new content. </p><p> You can make your database retain data indefinitely. This is useful if you want to keep a personal archive of all network. However, that means it will grow to an arbitrary size, and it won't be useful to others since other nodes will treat any data older than their own lifespan setting setting as gone. </p><p>The database also has a maximum size that can be set. If this size is reached (by default, 10Gb) before the livespan is reached, the node will start to delete from history starting from the oldest, even if it is otherwise within the 6-month lifespan. </p>"></a-info-marker>
                </div>
                <div class="flex-spacer"></div>
                <a-guidelight :notifyState="dotStates.dbDotState"></a-guidelight>
              </div>
              <div class="sub-body">
                <div class="spacer-row"></div>
                <div class="row progress-bar-row">
                  <div class="progress-bar-meta">
                    <div class="row-text">Disk used</div>
                    <div class="flex-spacer"></div>
                    <div class="row-text">Total allowed</div>
                  </div>
                  <a-progress-bar class="prog-bar" :percent="bas.dbsizemb" :max="bas.maxdbsizemb"></a-progress-bar>
                  <div class="progress-bar-meta">
                    <div class="left-space">{{bas.dbsizemb}} Mb</div>
                    <div class="flex-spacer"></div>
                    <div class="right-space">{{bas.maxdbsizemb}} Mb</div>
                  </div>
                </div>
                <div class="spacer-row"></div>
                <div class="row">
                  <div class="row-text">Status</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{bas.databasestatus}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Last insert</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderTimestamp(bas.lastdbinserttimestamp)}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Last insert duration</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderDuration(bas.lastinsertdurationseconds)}}</div>
                </div>

              </div>
            </div>

            <div class="sub-block">
              <div class="sub-header">
                <div class="sub-header-text">Caching
                </div>

                <div class="info-marker-container" v-show="loadComplete">
                  <a-info-marker header="Caches are pre-baked responses your node generates to efficiently serve other nodes with content." text="<p>Your node generates static caches every 6 hours by default. When a remote node tries to sync, the first thing it does is to consume your caches starting from its last sync timestamp. Since caches are pre-generated, serving these caches are effectively free (that is, only consumes bandwidth). </p><p> When a remote node reads the end of the most recent cache, it sends the end timestamp as the start timestamp of the live request. As a result, your node will only have to actively query and serve the smaller set of changes that has occurred since the end of the last cache.</p><p>If the live part of node's query ends up big enough that it takes more than one cache page, your node will save that response into an ad-hoc cache and serve that. <p> Caches are critical to correct and efficient function of your node. A node with a non-functional caching system will struggle to get the locally created content out, because other nodes will mark it as poorly behaved and prefer other nodes over it.</p> <p>This component runs in intervals.</p>"></a-info-marker>
                </div>
                <div class="flex-spacer"></div>
                <a-guidelight :notifyState="dotStates.cachingDotState"></a-guidelight>
              </div>
              <div class="sub-body">
                <div class="row">
                  <div class="row-text">Status</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{bas.cachingstatus}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Last cache generation</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderTimestamp(bas.lastcachegenerationtimestamp)}}</div>
                </div>
                <div class="row">
                  <div class="row-text">Last cache generation duration</div>
                  <div class="flex-spacer"></div>
                  <div class="row-data">{{renderDuration(bas.lastcachegenerationdurationseconds)}}</div>
                </div>
              </div>
            </div>

          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
  var globalMethods = require('../../services/globals/methods')
  var Vuex = require('../../../../node_modules/vuex').default
  var fe = require('../../services/feapiconsumer/feapiconsumer')
  interface Inflight {
    EntityType: string;
    ProgressPercentage: number;
    StateText: string;
    RequestedTimestamp: number;
    LastActionTimestamp: number;
    ActionType: string;
  }
  export default {
    name: 'status',
    data() {
      return {
        loadComplete: false,
      }
    },
    beforeMount(this: any) {
      fe.RequestAmbientStatus(function(resp: any) {
        console.log(resp)
      })
    },
    mounted(this: any) {
      let vm = this
      setTimeout(function() { vm.loadComplete = true }, 0)
      /*
        ^ I'm not super sure why this is needed, but it appears that if you include this without this gate, it actually makes the info markers shift slightly right after the page paint, which is ugly. Putting it behind a setTimeout makes it appear at their final place. I also tried Vue.nextTick as a more canonical replacement, but that did not help.

        Likely this is something related to which component gets priority in paint.
      */
    },
    computed: {
      ...Vuex.mapState(['ambientStatus', 'dotStates']),
      // /*----------  Frontend graph compiler  ----------*/
      // frontendLastCompile(this: any): string {
      //   if (globalMethods.IsUndefined(this.ambientStatus.frontendambientstatus.refreshstatus.lastrefreshtimestamp)) {
      //     return 'Unknown'
      //   }
      //   return globalMethods.TimeSince(this.ambientStatus.frontendambientstatus.refreshstatus.lastrefreshtimestamp)
      // },
      // frontendLastCompileDuration(this: any): string {
      //   if (globalMethods.IsUndefined(this.ambientStatus.frontendambientstatus.refreshstatus.lastrefreshdurationseconds)) {
      //     return 'Unknown'
      //   }
      //   return this.ambientStatus.frontendambientstatus.refreshstatus.lastrefreshdurationseconds + 's'
      // },
      /*----------  Inflights  ----------*/
      mergedInflights(this: any): Inflight[] {
        // This creates one single list from the inflights sorted by recency.
        let ifl: Inflight[] = []
        for (let val of this.ambientStatus.inflights.boardsList) {
          ifl.push({
            EntityType: 'Community',
            ProgressPercentage: val.status.completionpercent,
            StateText: val.status.statustext,
            RequestedTimestamp: val.status.requestedtimestamp,
            LastActionTimestamp: val.status.lastactiontimestamp,
            ActionType: val.status.eventtype,
          })
        }
        for (let val of this.ambientStatus.inflights.threadsList) {
          ifl.push({
            EntityType: 'Thread',
            ProgressPercentage: val.status.completionpercent,
            StateText: val.status.statustext,
            RequestedTimestamp: val.status.requestedtimestamp,
            LastActionTimestamp: val.status.lastactiontimestamp,
            ActionType: val.status.eventtype,
          })
        }
        for (let val of this.ambientStatus.inflights.postsList) {
          ifl.push({
            EntityType: 'Post',
            ProgressPercentage: val.status.completionpercent,
            StateText: val.status.statustext,
            RequestedTimestamp: val.status.requestedtimestamp,
            LastActionTimestamp: val.status.lastactiontimestamp,
            ActionType: val.status.eventtype,
          })
        }
        for (let val of this.ambientStatus.inflights.votesList) {
          ifl.push({
            EntityType: 'Vote',
            ProgressPercentage: val.status.completionpercent,
            StateText: val.status.statustext,
            RequestedTimestamp: val.status.requestedtimestamp,
            LastActionTimestamp: val.status.lastactiontimestamp,
            ActionType: val.status.eventtype,
          })
        }
        for (let val of this.ambientStatus.inflights.keysList) {
          ifl.push({
            EntityType: 'User',
            ProgressPercentage: val.status.completionpercent,
            StateText: val.status.statustext,
            RequestedTimestamp: val.status.requestedtimestamp,
            LastActionTimestamp: val.status.lastactiontimestamp,
            ActionType: val.status.eventtype,
          })
        }
        for (let val of this.ambientStatus.inflights.truststatesList) {
          ifl.push({
            EntityType: 'Truststate',
            ProgressPercentage: val.status.completionpercent,
            StateText: val.status.statustext,
            RequestedTimestamp: val.status.requestedtimestamp,
            LastActionTimestamp: val.status.lastactiontimestamp,
            ActionType: val.status.eventtype,
          })
        }
        ifl = this.sortByKey(ifl, 'RequestedTimestamp')
        return ifl
      },
      /*----------  Backend Ambient Status  ----------*/
      bas(this: any) {
        let bas = Object.assign({}, this.ambientStatus.backendambientstatus)
        // let bas = this.ambientStatus.backendambientstatus
        Object.keys(bas).forEach(function(key) {
          if (bas[key] === 0 || bas[key] === "" || globalMethods.IsUndefined(bas[key])) {
            // The ones below are excluded from labeled 'unknown' and will show zeroes.
            if (key === "dbsizemb" || key === "maxdbsizemb" || key === "lastinsertdurationseconds" || key === "lastcachegenerationdurationseconds" || key === "inboundscount15" || key === "outboundscount15") {
              return
            }
            bas[key] = 'Unknown'
          }
        })
        return bas
      },
      /*----------  Frontend Ambient status  ----------*/
      fas(this: any) {
        let fas = Object.assign({}, this.ambientStatus.frontendambientstatus)
        // let fas = this.ambientStatus.frontendambientstatus
        Object.keys(fas).forEach(function(key) {
          if (fas[key] === 0 || fas[key] === "" || globalMethods.IsUndefined(fas[key])) {
            // The ones below are excluded from labeled 'unknown' and will show zeroes.
            if (key === "lastrefreshdurationseconds") {
              return
            }
            fas[key] = 'Unknown'
          }
        })
        return fas
      },
    },
    methods: {
      sortByKey(this: any, array: any, key: string) {
        return array.sort(function(a: any, b: any) {
          var x = a[key]
          var y = b[key]
          return ((x > y) ? -1 : ((x < y) ? 1 : 0))
        });
      },
      renderTimestamp(this: any, ts: any) {
        if (ts !== 'Unknown') {
          return globalMethods.TimeSince(ts)
        }
        return ts
      },
      renderDuration(this: any, sec: any) {
        if (sec !== 'Unknown') {
          return sec + 's'
        }
        return sec
      },
      // renderDotStatees(this: any) {
      //   this.renderRefresherDotState()
      //   this.renderInflightsDotState()
      //   this.renderNetworkDotState()
      //   this.renderDbDotState()
      //   this.renderCachingDotState()
      //   this.renderFrontendDotState()
      //   this.renderBackendDotState()
      // },
      // renderFrontendDotState(this: any) {
      //   if (this.dotStatees)
      // }
    }
  }
</script>

<style lang="scss" scoped>
  @import "../../scss/globals";

  .flex-spacer {
    flex: 1;
  }

  .location.status {
    padding: 20px;
    padding-bottom: 50px;
    color: $a-grey-600;
    min-height: 100%;
    .status-cards-container {
      display: flex;
      flex-direction: row;
      // width: 1000px; // Same as settings
      width: 775px;
      margin: 0 auto;
    }
    .card-block {
      flex: 1;
      &:first-of-type {
        margin-right: 10px;
      }
      &:last-of-type {
        margin-left: 10px;
        margin-right: 0;
      }
    }
    .status-card {
      // margin: 15px;
      // margin-top: 10px;
      padding: 15px 20px;
      background-color: rgba(0, 0, 0, 0.25);
      border-radius: 3px; // width: 450px;
      .card-header {
        font-size: 160%;
        border-bottom: 3px solid rgba(255, 255, 255, 0.25);
        padding-bottom: 10px;
        display: flex;
        .info-marker-container {
          display: flex;
          div {
            margin: auto;
            margin-left: 6px;
            margin-bottom: 9px;
          }
        }
      }
      .card-header-text {}
      .card-body {}
    }
  }

  .sub-block {
    font-family: "SSP Regular";
    padding-top: 8px;
    padding-bottom: 15px;
    border-bottom: 3px solid rgba(255, 255, 255, 0.25);
    &:last-of-type {
      border-bottom: none;
      padding-bottom: 10px;
    }
    .sub-header {
      font-size: 125%;
      letter-spacing: 0.2px;
      padding-bottom: 5px;
      display: flex;

      .sub-header-text {
        // padding-right: 1px;
      }
      .info-marker-container {
        width: 12px;
        margin-left: 5px;
        display: flex;
        div {
          margin: auto;
          margin-bottom: 7px;
        }
      }
    }
  }

  .sub-body {}

  .row {
    display: flex;
    .info-marker-container {
      margin-left:4px;
    }
    .row-text {
    }
    .row-data {
      text-align: right;
      max-width: 200px;
    }
    &.progress-bar-row {
      // padding-top: 4px;
      // padding-bottom: 3px;
      display: block;
      .prog-bar {
        padding-top: 3px;
      }
      .progress-bar-meta {
        display: flex;
      }
    }
  }

  .inflights-table {
    width: 100%;
    font-family: "SCP Regular";
    font-size: 90%;
    border-collapse: collapse;
    tr {
      // border-bottom: 1px solid rgba(255, 255, 255, 0.25);
      td {
        padding: 2px 5px;
        &.nothing-in-progress {
          padding: 5px 7px;
        }
      }
      th {
        padding: 2px 5px;
        text-align: left;
        &:last-of-type {
          text-align: right; // padding-right: 0;
          width: 100px;
        }
      }
    }
    tbody tr {
      &:nth-child(odd) {
        background-color: rgba(255, 255, 255, 0.05);
      }
    }
  }

  .spacer-row {
    height: 10px;
  }
</style>