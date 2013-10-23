angular.module('aether.services')
    .factory('frameViewStateBroadcast', function($rootScope) {
    var frameViewState = {}
    frameViewState.secondFrame = ''
    frameViewState.thirdFrame = ''
    frameViewState.id = ''

    frameViewState.receiveState = function(secondFrame, thirdFrame, id) {
        this.secondFrame = secondFrame
        this.thirdFrame = thirdFrame
        this.id = id
        this.broadcastState()
    }

    frameViewState.broadcastState = function() {
        $rootScope.$broadcast('frameViewStateChanged')
    }

    return frameViewState
})

// angular.module('aether.services')
//     .factory('frameViewStateBroadcastWithoutIdRequest', function($rootScope) {
//         var frameViewState = {
//             secondFrame: '',
//             thirdFrame: '',
//             receiveState: function(secondFrame, thirdFrame) {
//                 this.secondFrame = secondFrame
//                 this.thirdFrame = thirdFrame
//                 $rootScope.$broadcast('frameViewStateWithoutIdChanged')
//             }
//          }

//     })

// angular.module('aether.services')
//     .factory('frameSizeBroadcast', function($rootScope) {
//     var frameSize = {}
//     frameSize.secondFrameSize = ''

//     frameSize.receiveSize = function(secondFrameSize) {
//         this.secondFrameSize = secondFrameSize
//         this.broadcastSize()
//     }

//     frameSize.broadcastSize = function() {
//         $rootScope.$broadcast('frameSizeChanged')
//     }

//     return frameSize
// })

angular.module('aether.services').factory('refreshService',
    function($rootScope) {
        var refresh = function() {
            $rootScope.$broadcast('refreshPage')
        }
        return refresh
    })