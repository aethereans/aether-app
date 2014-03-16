function FindOrCreateTopicController($scope, $rootScope, $timeout, frameViewStateBroadcast,
    gateReaderServices, refreshService, $interval) {

    $scope.topicSearchBoxText = ''



    $scope.postText = ''
    $scope.trimmedPostText = ''

    $scope.postButtonDisabled = true

    var logo = document.getElementsByClassName('aether-logo-main')[0]
    var counter = 0
    function spinLogo() {
        logo.style.webkitTransform = 'rotate('+((counter+1)/3)+'deg)'
        if (counter/3 < 360) {
            counter ++
        }
        else {
            counter = 1
        }

    }
    var spinner = $interval(spinLogo, 16)

    $scope.$on("$destroy", function() {
        $interval.cancel(spinner)
    })

    $rootScope.countsAvailable = false

    function refreshTopics() {
        gateReaderServices.getUppermostTopics(topicsArrived)
        function topicsArrived(data) {
            data = data.sort() // First alphabetic order, then numeric order based on reply count.
            data = data.sort(compareByReplyCount)
            var filteredTopics = []
            for (var i=0; i<data.length; i++) {
                // If counter is zero, not locally created, and not a selected topic.
                if (data[i].ReplyCount === 1 &&
                    data[i].LocallyCreated === false &&
                    $rootScope.userProfile.selectedTopics.indexOf(data[i].PostFingerprint) === -1) {
                    continue
                }
                else {
                    filteredTopics.push(data[i])
                }
            }
            // Filter topics according to their posts in it.

            function compareByReplyCount(b, a) {
                 return a.ReplyCount - b.ReplyCount
            }
            //filteredTopics = filteredTopics.sort(compareByReplyCount)
            $scope.allAvailableTopics = filteredTopics

            $scope.filterByBoxText = function() {
                var searchString = $scope.topicSearchBoxText.toLowerCase()
                var searchResultTopics = []
                if (searchString) {
                    for (var i=0; i<data.length;i++) {
                        if (data[i].Subject.toLowerCase().indexOf(searchString) != -1) {
                            searchResultTopics.push(data[i])
                        }
                    }
                    $scope.allAvailableTopics = searchResultTopics
                }
                else
                {
                    $scope.allAvailableTopics = filteredTopics
                }


            }
            console.log($scope.allAvailableTopics.length + 'topics available, and counts available:', $rootScope.countsAvailable)
            if (!$scope.allAvailableTopics.length) {
                // If there are no topics yet.
                $timeout(refreshTopics, 5000)
                return
            }
            else
            {
                // if there are topics, but not their counts.
                for (var i=0;i<$scope.allAvailableTopics.length;i++) {
                    if ($scope.allAvailableTopics[i].ReplyCount) { // if there is a reply count to any item, return
                        $rootScope.countsAvailable = true
                        return
                    }
                }
                // If no reply counts in the entire stack
                $timeout(refreshTopics, 5000)
            }
        }
    }

    refreshTopics()

    $scope.$watch('countsAvailable', function() {
        if ($scope.countsAvailable)
        {
            document.getElementById('find-or-create-topic-feed').style.height = 'auto'
        }
    })


    $scope.toggleAdditionToSelectedTopics = function(fingerprint) {
        if($rootScope.userProfile.selectedTopics
            .indexOf(fingerprint) > -1 ) {
            // If it exists, remove
            var index = $rootScope.userProfile.selectedTopics.indexOf(fingerprint)
            $rootScope.userProfile.selectedTopics.splice(index, 1)
        }
        else
        {
            $rootScope.userProfile.selectedTopics.push(fingerprint)
        }
        refreshService() // This is a broadcast to initiate refresh.
    }

    $scope.isCurrentlySelected = function(fingerprint) {
        return $rootScope.userProfile.selectedTopics
        .indexOf(fingerprint) > -1 ? true : false
    }

    $scope.$watch('postText', function() {
        $scope.trimmedPostText = $scope.postText.trim()
        enforceHeaderMaxLength($scope.trimmedPostText, 60)
        $scope.topicLetterCount = $scope.trimmedPostText.length
        if ($scope.postText.length > 3) {
            $scope.postButtonDisabled = false
        }
        else {
            $scope.postButtonDisabled = true
        }
    })

    function enforceHeaderMaxLength(text, maxLength) {
        $scope.headerLetterCount = text.length
        if (text.length > maxLength) {
            $scope.postText = text.slice(0, -(text.length - maxLength))
            // [0] of any jQuery or angular.element element is the element itself.
            setCaretToEndOfContenteditable(
                angular.element(document.getElementsByClassName('topic-name-entry'))[0])
        }

        function setCaretToEndOfContenteditable(target) {
            var range = document.createRange() // A range is an invisible selection
            range.selectNodeContents(target) // select node contents invisibly
            range.collapse(false) // collapse the selection to the end
            var selection = window.getSelection() // create a visible selection obj.
            selection.removeAllRanges() // remove the range within even if exists
            selection.addRange(range) // and make the range the selection's range.
        }

    }


    $scope.createTopic = function() {
        // here, set a scope variable that will modify the DOM to show a
        // creation screen, and upon the acceptance of that screen it will
        // commit the values to the database and refresh. the relatively hrd
        // part is to fit an input design that will fit this screen. I am also
        // feeling a little bit unsure about the green as highlight.. Too
        // bright, I would say. And the green hue doesn't fully reflect the
        // default aether green.
        var inputText = $scope.postText.trim()
        var topicFound = false
        for (var i = 0; i<$scope.allAvailableTopics.length; i++) {
            if (inputText === $scope.allAvailableTopics[i].Subject) {
                $rootScope.userProfile.selectedTopics.push(
                    $scope.allAvailableTopics[i].PostFingerprint)
                topicFound = true
                $scope.postText = ''
                $rootScope.scrollSecondFrameToTop()
                break
            }
        }
        if (topicFound === false) {
            gateReaderServices.createTopic(answerArrived, inputText)
            function answerArrived(data) {
                $rootScope.userProfile.selectedTopics.push(data)
                $scope.postText = ''
            }
        }
        angular.element(document.getElementsByClassName('topic-name-entry')).text('')
        refreshTopics()


        // This for some reason does not remove the add topic button to its
        // original state. Investigate.
        $scope.createFieldsActive = false

    }


    $scope.showCreateFields = function() {
        $scope.createFieldsActive = true
    }

}

FindOrCreateTopicController.$inject = ['$scope', '$rootScope', '$timeout', 'frameViewStateBroadcast',
    'gateReaderServices', 'refreshService', '$interval']
