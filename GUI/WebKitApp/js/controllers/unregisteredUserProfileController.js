function UnregisteredUserProfileController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices) {

    gateReaderServices.getUnregisteredUserPosts(dataArrived, $rootScope.userProfile.UserDetails.Username)
    function dataArrived(data) {
        $scope.feed = data
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
UnregisteredUserProfileController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']