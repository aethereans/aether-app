export { }

var globalMethods = require('../services/globals/methods')

let mixins = {
  localUserMixin: {
    computed: {
      localUserReadOnly(this: any) {
        if (globalMethods.IsUndefined(this.$store.state.localUser)) {
          return true
        }
        if (globalMethods.IsEmptyObject(this.$store.state.localUser)) {
          return true
        }
        if (this.$store.state.localUser.fingerprint.length === 0) {
          return true
        }
        return false
      },
      localUserExists(this: any) {
        return this.$store.state.localUserExists
      },
      localUserArrived(this: any) {
        return this.$store.state.localUserArrived
      }
    }
  },

}
module.exports = mixins