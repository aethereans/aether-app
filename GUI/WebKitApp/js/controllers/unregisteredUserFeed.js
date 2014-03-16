function UnregisteredUserFeed($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices) {

    gateReaderServices.getUnregisteredUserPosts(dataArrived, $rootScope.requestedUsername)
    function dataArrived(data) {
        var leftFeed = []
        var rightFeed = []
        for (var i =0 ; i<data.length; i++) {
            if (i%2===0) {
                leftFeed.push(data[i])
            }
            else {
                rightFeed.push(data[i])
            }
        }
        $scope.leftFeed = leftFeed
        $scope.rightFeed = rightFeed
    }

    $scope.decideNgInclude = function(postSubject) {
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
        $rootScope.thirdFrameCSSStyle = {
            'display':'block'
        }
    }

}
UnregisteredUserFeed.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']