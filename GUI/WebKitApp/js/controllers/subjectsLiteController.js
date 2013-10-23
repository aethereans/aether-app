function SubjectsLiteController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices) {
    $scope.$watch('requestedId', function(){
        if ($rootScope.requestedId !== undefined) {
            gateReaderServices.getParentPost(parentIdArrived, $rootScope.requestedId)
        }
    })

    function parentIdArrived(data) {
        $rootScope.thereIsASelectedSubject = true
        $rootScope.parentId = data.PostFingerprint
        $scope.selectedTopic = data
        gateReaderServices.getAllSpecificDepthPosts(subjectsArrived, $rootScope.parentId, 1)


        function subjectsArrived(data) {
            // Remove the first element, because it is the
            // subject itself.
            data.splice(0,1)
            $scope.subjects = data
        }
    }

    $scope.returnToTopics = function() {
        frameViewStateBroadcast.receiveState('', 'topicsFeedLite', '')
    }

    // Learned hard way: if you capture an event on a higher scope, it won't go
    // deeper. this means any $scope.on(..) that has already been captured at
    // frame controllers won't even get triggered here.

    // TODO: Here, check if the controller can access 2ndframeid. If it can, it's
    // at 3rd frame. In that case, write a ng-class directive that can look for
    // item id and activate on match so the one lights up.

    // I also need to determine here if this is in 2nd or 3rd scope, or write a
    // different controller for lite version.
}
SubjectsLiteController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']