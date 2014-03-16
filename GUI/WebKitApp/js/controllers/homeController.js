function HomeController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices) {

    // Here I need to get the list of selected topics, get their names, then for
    // each single topic, get 10 top most popular subjects.

    // okay.. what does this need? I need to get top x topics , the topmost being
    // 'popular' and in each I need to have top 10 most popular subjects.
    // I only need topic's name. so I can just have named arrays as topics and
    // posts under as array elements.

    // Get all the subjects user has selected,

    gateReaderServices.getLastConnectionTime(lastTimeArrived)
    function lastTimeArrived(reply) {
        var minutes = Math.round((Date.now()/1000 - reply)/60)
        if (minutes > 59) {
            $scope.lastTimeText = 'Last updated ' + Math.floor(minutes/60) + ' hours ago.'
        }
        else {
            $scope.lastTimeText = 'Last updated ' + minutes + ' minutes ago.'
        }

    }

    $scope.totalTopicCount = 0
    $scope.totalSubjectCount = 0
    //console.log('asking for topic count')
//    gateReaderServices.countTopics(topicCountArrived)
//    function topicCountArrived(count) {
//        $scope.totalTopicCount = count
//    }
    // This costs me extra 20ms compared to above. Determine if worth it. TODO
    gateReaderServices.getUppermostTopics(topicsArrived)
    function topicsArrived(data) {
        //$scope.subjectCountObject = {}
        var filteredTopics = []
        for (var i=0; i<data.length; i++) {
            if (data[i].ReplyCount === 1 && data[i].LocallyCreated === false) {
                continue
            }
            else {
                filteredTopics.push(data[i])
            }
        }
        $scope.totalTopicCount = filteredTopics.length

    }




    gateReaderServices.countSubjects(subjectCountArrived)
    function subjectCountArrived(count) {
        $scope.totalSubjectCount = count
    }

    gateReaderServices.getHomeScreen(homeScrArrived, 10,12)
    function homeScrArrived(data) {


        var filteredArray = []
        for (var i=0;i<data.length;i++) {
            if (data[i].Subjects.length > 0) {
                var topic = data[i]
                topic.HomeCol1 = []
                topic.HomeCol2 = []
                topic.HomeCol3 = []
                //console.log(topic.Subjects)
                for (var j=0;j<topic.Subjects.length;j++) {
                    if (j%3===0) {
                        topic.HomeCol1.push(topic.Subjects[j])
                    }
                    else if(j%3===1) {
                        topic.HomeCol2.push(topic.Subjects[j])
                    }
                    else if(j%3===2) {
                        topic.HomeCol3.push(topic.Subjects[j])
                    }
                }
                filteredArray.push(topic)
                // Do not show empty topics.
            }
        }
        $scope.topics = filteredArray
    }

    gateReaderServices.countConnectedNodes(function(count) {
        $scope.onlineNodeCount = count
    })


    gateReaderServices.getUpdateAvailable(function(reply) {
        $rootScope.updateAvailable = reply
    })

    $scope.replyButtonClick = function() {
        $rootScope.changeState('repliesFeed', '', '')
    }

    $scope.topicHeaderClick = function(postFingerprint) {
        $rootScope.changeState('postsFeed', 'subjectsFeedLite', postFingerprint)
        $rootScope.secondFrameCSSStyle = {
            'width': (getStyle(document.getElementById('root-body')), 'width').slice(0, -2) -
            (getStyle(document.getElementById('first-frame')), 'width').slice(0, -2) -
            (getStyle(document.getElementById('third-frame')), 'width').slice(0, -2) + 'px'
        }

        $rootScope.thirdFrameCSSStyle = {
            'display':'block'
        }
    }

    $scope.togglePauseClick = function() {
        if (!$rootScope.appIsPaused) {
            $rootScope.appIsPaused = true
            gateReaderServices.pauseApp()
        }
        else {
            $rootScope.appIsPaused = false
            gateReaderServices.resumeApp()
        }
    }

}
HomeController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
    'gateReaderServices']
