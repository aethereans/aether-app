"use strict";
// Globals > Methods
Object.defineProperty(exports, "__esModule", { value: true });
var vs = require('../../store/index');
// ^ For some reason, having default above does not work, returns undefined. We're just adding the default below and it should work.
var exportedMethods = {
    GetUserName: function (owner) {
        var vuexStore = vs.default;
        if (exportedMethods.IsUndefined(owner)) {
            return '';
        }
        if (typeof owner === 'string') {
            // We were given a fingerprint. check if the current user entity's fingerprint. If that's a match, run the function again with that user's entity and return the result. Otherwise, without making a call to frontend or backend, we can't resolve this to a user entity. This can happen when showing uncompiled entities.
            if (vuexStore.state.currentUserEntity.fingerprint === owner) {
                return this.GetUserName(vuexStore.state.currentUserEntity);
            }
            if (vuexStore.state.localUser.fingerprint === owner) {
                return this.GetUserName(vuexStore.state.localUser);
            }
            if (owner.length === 0) {
                return 'Not found';
            }
            return '@' + owner;
        }
        if (exportedMethods.IsUndefined(owner.fingerprint)) {
            // Necessary because the 'observer' object is not undefined, but also not what we want.
            return '';
        }
        if (owner.compiledusersignals.cnamesourcefingerprint.length > 0 &&
            owner.compiledusersignals.canonicalname.length > 0) {
            return '@' + owner.compiledusersignals.canonicalname;
        }
        if (owner.noncanonicalname.length > 0) {
            return '@' + owner.noncanonicalname;
        }
        if (owner.fingerprint.length > 0) {
            return '@' + owner.fingerprint;
        }
        return 'Not found';
    },
    TimeSince: function (timestamp) {
        var now = new Date();
        var ts = new Date(timestamp * 1000);
        var secondsPast = Math.floor((now.getTime() - ts.getTime()) / 1000);
        if (secondsPast === 0) {
            return 'just now';
        }
        if (secondsPast < 60) {
            return secondsPast + 's ago';
        }
        if (secondsPast < 3600) {
            return Math.floor(secondsPast / 60) + 'm ago';
        }
        if (secondsPast <= 86400) {
            return Math.floor(secondsPast / 3600) + 'h ago';
        }
        // If older than a day
        var day;
        var month;
        var year;
        day = ts.getDate();
        var tsds = ts.toDateString().match(/ [a-zA-Z]*/);
        if (tsds !== null) {
            month = tsds[0].replace(' ', '');
        }
        year = ts.getFullYear() === now.getFullYear() ? '' : ' ' + ts.getFullYear();
        return day + ' ' + month + year;
    },
    NowUnix: function () {
        return Math.floor(new Date().getTime() / 1000);
    },
    IsEmptyObject: function (obj) {
        if (exportedMethods.IsUndefined(obj)) {
            return false;
        }
        if (Object.keys(obj).length === 0 && obj.constructor === Object) {
            return true;
        }
        return false;
    },
    IsUndefined: function (obj) {
        if (typeof obj === 'undefined') {
            return true;
        }
        return false;
    },
};
module.exports = exportedMethods;
//# sourceMappingURL=methods.js.map