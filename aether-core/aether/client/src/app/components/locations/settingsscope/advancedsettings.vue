<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content"></a-markdown>
    <a-markdown :content="addNodeContent"></a-markdown>
    <div class="add-address-container">
      <h3>Node details</h3>
      <ul>
        <li>
          Your node needs to be online and accessible <i>now</i>. Your node will
          immediately attempt to establish a connection.
        </li>

        <li>
          If the connection fails for any reason, the node will not be added to
          the existing nodes database.
        </li>
      </ul>
      <template v-if="!addressSendInProgress">
        <!-- Not in progress: not started or done -->
        <template v-if="!addressSendResultArrived">
          <!-- Not started -->
          <div class="address-composer-container">
            <a-composer
              class="address-composer"
              :spec="addressComposerSpec"
            ></a-composer>
          </div>
        </template>
        <template v-else>
          <!-- Done -->
          <template v-if="addressSendSuccessful">
            <!-- Successful -->
            <div class="result-box">
              <p class="bold">Node sync successful</p>
              <p>
                Your node is successfully synced with, and added to the backend
                as an online node.
              </p>
            </div>
          </template>
          <template v-else>
            <!-- Failure -->
            <div class="result-box">
              <p class="bold">Node sync failed.</p>
              <p>
                Please check your values, make sure the remote is online and
                accessible, and return to this screen to try again.
              </p>
              <p class="bold">Error received</p>
              <p>
                <i>{{ addressSendResultErrorMessage }}</i>
              </p>
            </div>
          </template>
        </template>
      </template>
      <template v-else>
        <!-- In progress: spinner -->
        <div class="result-box">
          <p class="bold">Attempting to sync with the remote node...</p>
          <p>
            Depending on a) how busy either node is, b) graph delta between
            nodes, this process can take from a couple seconds to 10 minutes.
          </p>
          <p>
            If you'd like to see the result of this process (as well as the
            error message, if any), do not leave this page until this process is
            complete.
          </p>
          <div class="spinner-container">
            <div class="spinner-carrier">
              <a-spinner :hidetext="true"></a-spinner>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
var fe = require('../../../services/feapiconsumer/feapiconsumer')
var mimobjs = require('../../../../../../protos/mimapi/mimapi_pb.js')
let globalMethods = require('../../../services/globals/methods')
export default {
  name: 'advancedsettings',
  data(this: any) {
    return {
      headline: headline,
      intro: intro,
      content: content,
      addNodeContent: addNodeContent,
      addressSendResultArrived: false,
      addressSendInProgress: false,
      addressSendSuccessful: false,
      addressSendResultErrorMessage: '',
      addressComposerSpec: {
        fields: [
          {
            id: 'addressLocation',
            emptyWarningDisabled: false,
            visibleName: 'Location',
            description:
              'Location (URL or IPv4) of your node. <br>If URL, omit <b>http://</b> or <b>https://</b>. Aether is TLS-only, <b>https</b> will be assumed.',
            placeholder: '172.31.88.174',
            maxCharCount: 1024,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: false,
          },
          {
            id: 'addressSub',
            visibleName: 'URL Sublocation',
            description:
              "Optional. If you're using an IP address (almost always), leave blank.  <br>Example: If your node is at <b>https://www.example.com/myhome/mynode</b>, <b>www.example.com</b> is the location, and <b>myhome/mynode</b> is the sublocation.",
            placeholder: 'myhome/mynode',
            maxCharCount: 1024,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: true,
          },
          {
            id: 'addressPort',
            emptyWarningDisabled: false,
            visibleName: 'Port',
            description: '',
            placeholder: '49999',
            maxCharCount: 5,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: false,
          },
        ],
        commitAction: this.sendAddressEntity,
        commitActionName: 'SUBMIT',
        cancelAction: function () {},
        cancelActionName: '',
        autofocus: false,
      },
    }
  },
  methods: {
    sendAddressEntity(this: any, fields: any) {
      this.addressSendInProgress = true
      let addressLocation = ''
      let addressSublocation = ''
      let addressPort = ''
      for (let val of fields) {
        if (val.id === 'addressLocation') {
          addressLocation = val.content
          continue
        }
        if (val.id === 'addressSublocation') {
          addressSublocation = val.content
          continue
        }
        if (val.id === 'addressPort') {
          addressPort = val.content
          continue
        }
      }
      let addr = new mimobjs.Address()
      addr.setLocation(addressLocation)
      addr.setSublocation(addressSublocation)
      addr.setPort(addressPort)
      let vm = this
      fe.SendAddress(addr, function (resp: any) {
        if (!globalMethods.IsUndefined(resp.reportErrorToServer)) {
          // This is an error
          vm.addressSendInProgress = false
          vm.addressSendResultArrived = true
          return
        }
        vm.addressSendInProgress = false
        vm.addressSendResultArrived = true
        if (resp.statuscode === 200) {
          vm.addressSendSuccessful = true
        } else {
          vm.addressSendResultErrorMessage = resp.errormessage
        }
      })
    },
  },
}
// These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
var headline = '# Advanced'
var intro = `**This describes the settings in more detail and provides instructions on how to change them.**

* **These descriptions are intended for power-users.**
All settings come with sane defaults. If none of this makes sense to you, you can safely ignore them.
  `
var content = `


### Changing preferences manually

1. Shut down the app completely.
2. Save your old config file as a copy somewhere else.
3. Go to the appropriate config file, and make your changes.
4. Restart. If the app fails to boot, bring back the old config file.


* **Descriptions of more consequential settings are below.**
For the rest, the descriptions can be found [here](https://github.com/Frigyes06/aether-app/blob/master/aether-core/aether/services/configstore/permanent.go).

* **Network settings, if misconfigured, can get your machine and user key permanently banned by other nodes.** The settings that modify local behaviour are generally safe to fiddle with. Ones that relate to network behaviour are not.

### Defaults

| Setting | Description | Value |
|--- | --- |--- |
| Maximum disk space for database | The maximum disk space the backend database is allowed to take. Whenever the app reaches this threshold, it starts to delete from history to remain under the threshold. <br><br> Mind that this is not the entire disk space used by this app. The other things using the cache are the pre-baked HTTP caches that are used to ease outbound serves, and the frontend key-value store that holds precompiled graph of human-readable objects, such as boards, threads, users. <br><br> Both of these take an additional 10-15% of the database size, so at the maximum DB size of 10 Gb, the total disk use would be something around 13 Gb. | 10 Gb
| Local memory | For how long the local node remembers data for, in absence of disk space pressure. <br><br> This means the data will be deleted at the 6-month mark even if the maximum disk space is not reached. If it is, the local memory will be less. In other words, you can have an 6-month local memory but give the maximum disk space a value of 1 Gb, you will likely have much less than 6 months worth of content.  <br><br> This is somewhat akin to a reference-counting garbage collector, in that it will traverse the content graph, and if an object is older than 6 months but there exists a vertex that links to a newer-than-6-months graph node, it will not be deleted.  <br><br> An example of this is board graph nodes, which will not be deleted as long as someone has posted a thread in them in the last 6 months. | 6 months |
| Neighbourhood size | This is the size of the local node's neighbourhood that it will try to keep in sync with. <br><br> This is a push/pop stack, and at the end of a cycle, the oldest neighbour will be evicted, and a new one that wasn't synced before will be added to it. This means the local node will be connecting to known nodes for the 90% of the time, and new ones 10% of the time. | 10
| Tick duration | This is the base unit of time. The local node will attempt to establish a connection to a standard live node every tick by popping a candidate from the neighbourhood.  <br><br> Every 10 ticks the neighbourhood is cycled by one (see 'Neighbourhood size' above), every 60 ticks the static nodes will be hit, and every 360 ticks, bootstrap nodes. | 60 seconds
| Reverse open | Reverse opens are a way to 'request' a connection from a remote node by connecting to that remote and passing over a raw Mim request asking for a reverse-connect using the same TCP socket.  <br><br> This is useful for nodes that are behind firewalls and uncooperating NATs. Without this, no other nodes would be able to connect to them directly in another way because of the firewall, rendering the content they create unable to reach the network.  <br><br> In case of erratic behaviour, this is a good first thing to consider disabling (as long as your network is configured right, and UPNP can port-map your router). | Enabled
| Maximum address table size | How many other nodes' addresses will be kept in the database. <br><br> Whenever this threshold is crossed, the addresses with the oldest last successful connection timestamp will be purged from memory. | 1000
| Maximum simultaneous inbound connections | The number of remotes that can be syncing with the local node at the same time. <br><br> Mind that the inbound and outbounds are different types of syncs, because syncs are one-way pulls. A node syncing with you doesn't mean that you get the changes on that node, it just means the remote node gets the changes in yours. This improves security, since no one can 'push' data into your machine. A result of this is that it is imperative that other nodes are able to connect to you, because if they do not, the content you create will never be able to leave your node, and reach the network. <br><br> Therefore setting this value to 0 might at first seem like a good hack to reduce bandwidth use at the expense of others, but it will also render you effectively invisible to everyone. The larger this value is, the faster your content will reach other users of the network, up to the point that your uplink bandwidth, CPU or disk is saturated to such a degree that the remotes are abandoning syncs with you because it takes so long. <br><br> If the app is taxing your computer too much, this is a good value to try reducing one by one. The default value is chosen as a balance between network connectivity and system resource use, and assuming you have a CPU made in the last decade, it should not be taxing your CPU very much, if at all. | 5
| Maximum simultaneous outbound connections | The number of remotes that the local node be syncing with at the same time. | 1
  `
var addNodeContent = `
## Manual node insert

* You can add a new node to your backend here.
* If you are not having issues or debugging something, you do not need this.

### The reasons you might want to do this

* **You are using Aether in a local network isolated from the Internet.** In this case, you need to add at least one node to another, so that they can discover each other for the first time.

  *Heads up - if you want to do this, you want to create a **Realm**. Realms are encrypted, private instances of Aether network that don't touch the main universe, and they maintain their own node lists. Realms are an upcoming feature.*

  *If you don't use a realm, if any computer on your local network ever touches the main network, it will dump all your content into the mainnet, as it considers it a repair of a temporary disconnect from the main network.*

* **You are testing connectivity.** Since Aether nodes' caches are accessible with a browser, it's usually much easier to just open your browser and type the node IP:Port to see if it connects. That said, this can help you test all-around connectivity and responsiveness.

* **Your node is getting content slowly.**

  - It is likely that you have a different problem than most nodes in your neighbourhood being extremely slow.

  - First, try a different internet connection.

  - In some rare cases, doing this might help.


### The reasons you do ***not*** want to do this

* **There's more content available in a specific node.**
This is usually not true. All Aether nodes keep a time-limited copy of all of the network. Connecting to any node is indistinguishable from connecting to another, assuming both nodes are tracking the \`\`\`NETWORK_HEAD\`\`\`.

* **The content you create is being broadcast slowly.**

  - This is usually because you have a hostile router (work network, coffee shop?) that won't let you map your external port and allow other nodes to connect to you.

  - In this case, your node will request other nodes to connect to you by opening an outbound connection to them, and hand over the ownership of connection. It is still within other nodes' decision to connect to you, and they can decline it.

  - This process is more involved than being able to accept inbound requests normally, and it will be slower to push your content out.

  - Adding a node won't help. To fix this, find a better Internet connection that is more friendly to you.

  - Your issue isn't that you don't have enough nodes, but that they cannot connect to you directly.

* **You just really want to connect to that specific node and have it in your neighbourhood.**

  - Neighbourhood is maintained programmatically on a inject/eject basis. Having any specific node in your neighbourhood does not come with any benefits. All nodes are the same, and behave the same.

  - In Aether, *the content graph is independent of network topology*.

`
</script>

<style lang="scss" scoped>
@import '../../../scss/globals';
.settings-sublocation {
  color: $a-grey-600;
  .markdowned {
    &:first-of-type {
      margin-bottom: 0;
    }
    margin-bottom: 40px;
  }
}

h1,
h2,
h3,
h4,
h5,
h6 {
  font-family: 'SSP Bold';
}

hr {
  background-color: rgba(255, 255, 255, 0.25);
  height: 3px;
  border: none;
}

.result-box {
  font-family: 'SCP Regular';
  margin: 20px;
  padding: 20px;
  background-color: rgba(0, 0, 0, 0.25);
  border-radius: 3px;
  font-size: 16px;
  .bold {
    font-family: 'SCP Bold';
  }
}

.spinner-container {
  display: flex;
  .spinner-carrier {
    margin: auto;
  }
}

.add-address-container {
  // font-family: "SSP Bold"
}
</style>

<style lang="scss">
.address-composer {
  .description b {
    font-family: 'SCP Bold';
    font-size: 93%;
  }
  .actions {
    font-family: 'SSP Bold';
  }
}
</style>
