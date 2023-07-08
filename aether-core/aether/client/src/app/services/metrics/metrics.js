"use strict";
// Services > Metrics
// This service, if enabled, collects non-identifying user metrics for network performance, product improvements and bug fixes.
Object.defineProperty(exports, "__esModule", { value: true });
// This is a little different than usual because this can take an argument in its instantiation.
var moduleCache;
var Store = require('electron-store');
var store = new Store();
var isDev = require('electron-is-dev');
module.exports = function (instantiatedByMainMain, mainMainMetricsDisabled) {
    // If we have an already instantiated module, return that.
    if (moduleCache) {
        return moduleCache;
    }
    // if (instantiatedByMainMain) {
    //   console.log('Metrics initialised by MainMain')
    //   console.log('MetricsDisabled: ', metricsDisabled())
    // } else {
    //   console.log('Metrics initialised by Renderer')
    //   console.log('MetricsDisabled: ', metricsDisabled())
    // }
    var nodeMixpanel = {};
    if (instantiatedByMainMain) {
        var mp = require('mixpanel');
        nodeMixpanel = mp.init('b48754d816a75407938965c24debbe46', {
            protocol: 'https',
        });
    }
    var mixpanelInstance = {};
    if (instantiatedByMainMain) {
        mixpanelInstance = nodeMixpanel;
    }
    else {
        mixpanelInstance = mixpanel;
    }
    // We don't, this is our first instantiation.
    var module = {};
    /*
      If this is instantiated by the MainMain, not the Renderer, this will be flipped true. This allows this library to be compatible with both.
    */
    var appVersionAndBuild = '';
    if (instantiatedByMainMain) {
        appVersionAndBuild = require('electron').app.getVersion().split('+');
    }
    else {
        appVersionAndBuild = require('electron').remote.app.getVersion().split('+');
    }
    var appVersion = appVersionAndBuild[0];
    var appBuild = appVersionAndBuild[1];
    /*----------  Public methods  ----------*/
    module.SendRaw = function (metric, payload) {
        // console.log('Metrics received a send request.')
        if (metricsDisabled()) {
            return;
            // If metrics are disabled, nothing goes out.
        }
        var fields = getFields(payload);
        insertAmbientInfo(fields);
        // console.log(metric, fields)
        if (!instantiatedByMainMain) {
            var did = store.get('did');
            if (!did) {
                // No known did beforehand
                store.set('did', 'MIXPANEL_ELECTRON_CLIENT_DID_NOT_LOAD');
                if (typeof mixpanelInstance.get_distinct_id !== 'undefined') {
                    store.set('did', mixpanelInstance.get_distinct_id());
                }
            }
            else {
                // We have a prior did
                if (typeof mixpanelInstance.cookie !== 'undefined') {
                    mixpanelInstance.cookie.props.distinct_id = did;
                    mixpanelInstance.cookie.props.$device_id = did;
                }
            }
        }
        if (!isDev) {
            mixpanelInstance.track(metric, fields);
            mixpanelInstance.track('App interacted', fields);
        }
    };
    module.SendContentEvent = function (contentType, eventType, customFields) {
        var fields = getFields(customFields);
        fields['A-Type'] = contentType;
        fields['A-EventType'] = eventType;
        module.SendRaw('ContentEvent', fields);
    };
    module.SendSignalEvent = function (signalType, eventType, customFields) {
        var fields = getFields(customFields);
        fields['A-Type'] = signalType;
        fields['A-EventType'] = eventType;
        module.SendRaw('SignalEvent', fields);
    };
    /*----------  Private methods  ----------*/
    function insertAmbientInfo(payload) {
        payload['A-Ambient-SourceType'] = 'App';
        payload['A-Ambient-AppVersion'] = appVersion;
        payload['A-Ambient-AppBuild'] = appBuild;
        if (instantiatedByMainMain) {
            payload['A-Ambient-Page'] = 'Client MainMain';
            payload['A-Ambient-PageName'] = 'Client MainMain';
            return;
        }
        var vuexStore = require('../../store/index').default;
        // payload['A-Ambient-Page'] = vuexStore.state.route.path
        payload['A-Ambient-Page'] = 'Client Renderer';
        payload['A-Ambient-PageName'] = vuexStore.state.route.name;
    }
    function getFields(customFields) {
        var fields = {};
        if (typeof customFields !== 'undefined') {
            fields = customFields;
        }
        return fields;
    }
    function metricsDisabled() {
        if (instantiatedByMainMain) {
            return mainMainMetricsDisabled;
        }
        var vuexStore = require('../../store/index').default;
        return vuexStore.state.metricsDisabled;
    }
    // Save it to the cache, so the next time it is instantiated, we use the prior.
    moduleCache = module;
    return module;
};
//# sourceMappingURL=metrics.js.map