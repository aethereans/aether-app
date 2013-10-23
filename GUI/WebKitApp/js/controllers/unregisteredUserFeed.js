function UnregisteredUserFeed($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices) {

    gateReaderServices.getUnregisteredUserPosts(dataArrived, $rootScope.requestedUsername)
    function dataArrived(data) {
        $scope.feed = data
        console.log(data)
    }

    $scope.decideNgInclude = function(postSubject) {
        console.log(postSubject)
        if (postSubject === '') {
            // It is a post
            return 'contentBlocks/userPostItem.html'
        }
        else {
            // It is a subject
            return 'contentBlocks/userSubjectItem.html'
        }
    }

    $scope.clickToPostItem = function(postFingerprint) {
        $rootScope.changeState('singleReply', '', postFingerprint)

    }

    $scope.clickToSubjectItem = function(postFingerprint) {
        $rootScope.changeState('postsFeed', 'subjectsFeedLite', postFingerprint)
        $rootScope.secondFrameCSSStyle = {}
        $rootScope.thirdFrameCSSStyle = {}
    }

}
UnregisteredUserFeed.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']