"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var globalMethods = require('../services/globals/methods');
var mixins = {
    localUserMixin: {
        computed: {
            localUserReadOnly: function () {
                if (globalMethods.IsUndefined(this.$store.state.localUser)) {
                    return true;
                }
                if (globalMethods.IsEmptyObject(this.$store.state.localUser)) {
                    return true;
                }
                if (this.$store.state.localUser.fingerprint.length === 0) {
                    return true;
                }
                return false;
            },
            localUserExists: function () {
                return this.$store.state.localUserExists;
            },
            localUserArrived: function () {
                return this.$store.state.localUserArrived;
            },
        },
    },
};
module.exports = mixins;
//# sourceMappingURL=mixins.js.map