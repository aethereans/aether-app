"use strict";
// Store > Status
Object.defineProperty(exports, "__esModule", { value: true });
var globalMethods = require('../services/globals/methods');
var dotStateOrderSection = {
    status_section_unknown: 0,
    status_section_ok: 1,
    status_section_warn: 2,
    status_section_fail: 3,
};
var dotStateOrderSubsection = {
    status_subsection_unknown: 0,
    status_subsection_ok: 1,
    status_subsection_warn: 2,
    status_subsection_fail: 3,
};
/*----------  Utility methods  ----------*/
function getTimeFromNowInMinutes(ts) {
    var now = globalMethods.NowUnix();
    return (now - ts) / 60;
}
/*----------  This is the core logic for red/yellow/green  ----------*/
/*
  Refresher
  ---
  WARN:
  - last refresh timestamp > 5 min
*/
function computeRefresherDotState(ambientStatus) {
    var state = 'status_subsection_ok';
    var fas = ambientStatus.frontendambientstatus;
    if (getTimeFromNowInMinutes(fas.lastrefreshtimestamp) > 5) {
        console.log('Emitting a status warning because: Last Refresher run timestamp > 5 min');
        state = 'status_subsection_warn';
        return state;
    }
    return state;
}
/*
  Inflights
  ---
  WARN:
  - inflights queue length > 10
*/
function computeInflightsDotState(ambientStatus) {
    var state = 'status_subsection_ok';
    var iflCount = 0;
    iflCount = iflCount + ambientStatus.inflights.boardsList.length;
    iflCount = iflCount + ambientStatus.inflights.threadsList.length;
    iflCount = iflCount + ambientStatus.inflights.postsList.length;
    iflCount = iflCount + ambientStatus.inflights.votesList.length;
    iflCount = iflCount + ambientStatus.inflights.keysList.length;
    iflCount = iflCount + ambientStatus.inflights.truststatesList.length;
    if (iflCount > 10) {
        console.log('Emitting a status warning because: Inflights queue length > 10');
        state = 'status_subsection_warn';
        return state;
    }
    return state;
}
/*
  Network
  ---
  WARN:
  - inboundscount15 < 2
  - last inbound conn timestamp > 10 min
  - upnp status != 'successful' or 'in progress' (actually this doesn't matter all that much so long as inbounds are coming in. NBD, let's skip this one for now.)
*/
function computeNetworkDotState(ambientStatus) {
    var state = 'status_subsection_ok';
    var bas = ambientStatus.backendambientstatus;
    if (bas.inboundscount15 < 2) {
        console.log('Emitting a status warning because: Less than 2 inbounds in the last 15 minutes');
        state = 'status_subsection_warn';
        return state;
    }
    if (getTimeFromNowInMinutes(bas.lastinboundconntimestamp) > 10) {
        console.log('Emitting a status warning because: Last inbound conn timestamp > 10 min');
        state = 'status_subsection_warn';
        return state;
    }
    return state;
}
/*
  Db
  ---
  WARN:
  - lastinsertdurationseconds > 360
*/
function computeDbDotState(ambientStatus) {
    var state = 'status_subsection_ok';
    var bas = ambientStatus.backendambientstatus;
    if (bas.lastinsertdurationseconds > 360) {
        console.log('Emitting a status warning because: Last db insert took > 5 min');
        state = 'status_subsection_warn';
        return state;
    }
    return state;
}
/*
  Caching
  ---
  WARN:
  - Last cache generation timestamp > 12 hours
*/
function computeCachingDotState(ambientStatus) {
    var state = 'status_subsection_ok';
    var bas = ambientStatus.backendambientstatus;
    if (getTimeFromNowInMinutes(bas.lastcachegenerationtimestamp) > 720 &&
        getTimeFromNowInMinutes(bas.lastcachegenerationtimestamp) != 0) {
        // ^ != 0 because cache generation starts a half hour after the app starts running. We don't want to show yellows in that time zone, it is working as intended.
        console.log('Emitting a status warning because: last cache generation timestamp is > 12h');
        state = 'status_subsection_warn';
        return state;
    }
    return state;
}
/*
No FAILs yet, because we don't know which one of these are critical until we actually run the real network. If we do it now, we might accidentally mark things as FAILs while in reality they're mild inconveniences.
*/
// The logic below is just to check the error states, the higher the number, the more severe the error state. We just show the highest error state as the parent error state.
function computeBackendDotState(ds) {
    var highestErrorState = 0;
    if (dotStateOrderSubsection[ds.networkDotState] > highestErrorState) {
        highestErrorState = dotStateOrderSubsection[ds.networkDotState];
    }
    if (dotStateOrderSubsection[ds.dbDotState] > highestErrorState) {
        highestErrorState = dotStateOrderSubsection[ds.dbDotState];
    }
    if (dotStateOrderSubsection[ds.cachingDotState] > highestErrorState) {
        highestErrorState = dotStateOrderSubsection[ds.cachingDotState];
    }
    var result = 'status_section_unknown';
    Object.keys(dotStateOrderSection).forEach(function (key) {
        if (dotStateOrderSection[key] === highestErrorState) {
            result = key;
        }
    });
    return result;
}
function computeFrontendDotState(ds) {
    var highestErrorState = 0;
    if (dotStateOrderSubsection[ds.refresherDotState] > highestErrorState) {
        highestErrorState = dotStateOrderSubsection[ds.refresherDotState];
    }
    if (dotStateOrderSubsection[ds.inflightsDotState] > highestErrorState) {
        highestErrorState = dotStateOrderSubsection[ds.inflightsDotState];
    }
    var result = 'status_section_unknown';
    Object.keys(dotStateOrderSection).forEach(function (key) {
        if (dotStateOrderSection[key] === highestErrorState) {
            result = key;
        }
    });
    return result;
}
/*----------  Higher level API  ----------*/
function computeDotStates(ambientStatus) {
    var dotStates = {
        /*----------  Main dot statuses  ----------*/
        backendDotState: 'status_section_unknown',
        frontendDotState: 'status_section_unknown',
        /*----------  Sub dot states  ----------*/
        refresherDotState: 'status_subsection_unknown',
        inflightsDotState: 'status_subsection_unknown',
        networkDotState: 'status_subsection_unknown',
        dbDotState: 'status_subsection_unknown',
        cachingDotState: 'status_subsection_unknown',
    };
    dotStates.refresherDotState = computeRefresherDotState(ambientStatus);
    dotStates.inflightsDotState = computeInflightsDotState(ambientStatus);
    dotStates.networkDotState = computeNetworkDotState(ambientStatus);
    dotStates.dbDotState = computeDbDotState(ambientStatus);
    dotStates.cachingDotState = computeCachingDotState(ambientStatus);
    dotStates.backendDotState = computeBackendDotState(dotStates);
    dotStates.frontendDotState = computeFrontendDotState(dotStates);
    return dotStates;
}
var statusActions = {
    // setDotStates updates the dot states we have based on the data that comes in from the frontend about other parts of the system.
    setDotStates: function (context, ambientStatus) {
        var ds = computeDotStates(ambientStatus);
        context.commit('SET_DOT_STATES', ds);
    },
};
exports.default = statusActions;
//# sourceMappingURL=statuslights.js.map