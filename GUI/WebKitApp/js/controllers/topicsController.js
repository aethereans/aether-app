function TopicsController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices) {




    $scope.setSubjectsToSingleColumn = function() {
        $rootScope.userProfile.subjectsSingleColumnLayout = true

    }
    $scope.setSubjectsToMultiColumn = function() {
        $rootScope.userProfile.subjectsSingleColumnLayout = false
    }


    $scope.$on('refreshPage', function() {
        gateReaderServices.getUppermostTopics(topicsArrived)
    })

    gateReaderServices.getNewborn(function(answer) {$scope.newborn = answer})

    gateReaderServices.getUppermostTopics(topicsArrived)
    function topicsArrived(posts)
    {

        // Sort topics alphabetically before posting them to DOM.
        posts.sort(function(a,b) {
            var aTopic = a.Subject.toUpperCase()
            var bTopic = b.Subject.toUpperCase()
            return (aTopic<bTopic) ? -1 : (aTopic>bTopic) ? 1 : 0
        })
        $scope.topics = []
        for (var i=0; i<posts.length; i++) {
            // check if posts[i] exists in userProfile.selectedTopics.
            if($rootScope.userProfile.selectedTopics
                .indexOf(posts[i].PostFingerprint) > -1 ) {
                // If exists in userProfile
                $scope.topics.push(posts[i])

            }
        }
        //console.log('these are topics: ', $scope.topics)

        $scope.$watch('requestedId', function(){
            if ($rootScope.requestedId === undefined && $scope.topics[0] != undefined)
            {
                $rootScope.requestedId = $scope.topics[0].PostFingerprint
            }
        })

    }

    $scope.topicClick = function(PostFingerprint) {
        $rootScope.thereIsASelectedSubject = false
        $scope.changeState('subjectsFeed', '', PostFingerprint)
    }

    $scope.returnToSubjects = function () {
        frameViewStateBroadcast.receiveState('', 'subjectsFeedLite', '')
    }

    $scope.toggleFindOrCreateTopic = function() {
        frameViewStateBroadcast.secondFrame == 'findOrCreateTopic' ?
        $scope.changeState('subjectsFeed', '', '') :
        $scope.changeState('findOrCreateTopic', '', '')
        $rootScope.thereIsASelectedSubject = false
    }

    $scope.isSelected = function(fingerprint) {
        if ($rootScope.thereIsASelectedSubject && $rootScope.parentId === fingerprint) {
            return true
        }
        else if (fingerprint == $scope.requestedId &&
        frameViewStateBroadcast.secondFrame != 'findOrCreateTopic') {
            return true
        }
        else {
            return false
        }
    }




}
TopicsController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']