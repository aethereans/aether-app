// Services > Metrics
// This service, if enabled, collects non-identifying user metrics for network performance, product improvements and bug fixes.

// This library can be imported by both the main and the renderer.

export {}

declare var mixpanel: any

// This is a little different than usual because this can take an argument in its instantiation.

var moduleCache: any

const Store = require('electron-store')
const store = new Store()
const isDev = require('electron-is-dev')

module.exports = function (
  instantiatedByMainMain: boolean,
  mainMainMetricsDisabled: boolean
) {
  // If we have an already instantiated module, return that.
  if (moduleCache) {
    return moduleCache
  }
  // if (instantiatedByMainMain) {
  //   console.log('Metrics initialised by MainMain')
  //   console.log('MetricsDisabled: ', metricsDisabled())
  // } else {
  //   console.log('Metrics initialised by Renderer')
  //   console.log('MetricsDisabled: ', metricsDisabled())
  // }
  let nodeMixpanel: any = {}
  if (instantiatedByMainMain) {
    var mp = require('mixpanel')
    nodeMixpanel = mp.init('b48754d816a75407938965c24debbe46', {
      protocol: 'https',
    })
  }
  let mixpanelInstance: any = {}
  if (instantiatedByMainMain) {
    mixpanelInstance = nodeMixpanel
  } else {
    mixpanelInstance = mixpanel
  }
  // We don't, this is our first instantiation.
  var module: any = {}
  /*
    If this is instantiated by the MainMain, not the Renderer, this will be flipped true. This allows this library to be compatible with both.
  */
  let appVersionAndBuild = ''
  if (instantiatedByMainMain) {
    appVersionAndBuild = require('electron').app.getVersion().split('+')
  } else {
    appVersionAndBuild = require('electron').remote.app.getVersion().split('+')
  }
  let appVersion = appVersionAndBuild[0]
  let appBuild = appVersionAndBuild[1]

  /*----------  Public methods  ----------*/
  module.SendRaw = function (metric: string, payload: any) {
    // console.log('Metrics received a send request.')
    if (metricsDisabled()) {
      return
      // If metrics are disabled, nothing goes out.
    }
    let fields = getFields(payload)
    insertAmbientInfo(fields)
    // console.log(metric, fields)
    if (!instantiatedByMainMain) {
      let did = store.get('did')
      if (!did) {
        // No known did beforehand
        store.set('did', 'MIXPANEL_ELECTRON_CLIENT_DID_NOT_LOAD')
        if (typeof mixpanelInstance.get_distinct_id !== 'undefined') {
          store.set('did', mixpanelInstance.get_distinct_id())
        }
      } else {
        // We have a prior did
        if (typeof mixpanelInstance.cookie !== 'undefined') {
          mixpanelInstance.cookie.props.distinct_id = did
          mixpanelInstance.cookie.props.$device_id = did
        }
      }
    }
    if (!isDev) {
      mixpanelInstance.track(metric, fields)
      mixpanelInstance.track('App interacted', fields)
    }
  }

  module.SendContentEvent = function (
    contentType: string,
    eventType: string,
    customFields: any
  ) {
    let fields = getFields(customFields)
    fields['A-Type'] = contentType
    fields['A-EventType'] = eventType
    module.SendRaw('ContentEvent', fields)
  }

  module.SendSignalEvent = function (
    signalType: string,
    eventType: string,
    customFields: any
  ) {
    let fields = getFields(customFields)
    fields['A-Type'] = signalType
    fields['A-EventType'] = eventType
    module.SendRaw('SignalEvent', fields)
  }

  /*----------  Private methods  ----------*/
  function insertAmbientInfo(payload: any): any {
    payload['A-Ambient-SourceType'] = 'App'
    payload['A-Ambient-AppVersion'] = appVersion
    payload['A-Ambient-AppBuild'] = appBuild
    if (instantiatedByMainMain) {
      payload['A-Ambient-Page'] = 'Client MainMain'
      payload['A-Ambient-PageName'] = 'Client MainMain'
      return
    }
    const vuexStore = require('../../store/index').default
    // payload['A-Ambient-Page'] = vuexStore.state.route.path
    payload['A-Ambient-Page'] = 'Client Renderer'
    payload['A-Ambient-PageName'] = vuexStore.state.route.name
  }

  function getFields(customFields: any): any {
    let fields = {}
    if (typeof customFields !== 'undefined') {
      fields = customFields
    }
    return fields
  }

  function metricsDisabled(): boolean {
    if (instantiatedByMainMain) {
      return mainMainMetricsDisabled
    }
    const vuexStore = require('../../store/index').default
    return vuexStore.state.metricsDisabled
  }

  // Save it to the cache, so the next time it is instantiated, we use the prior.
  moduleCache = module
  return module
}
