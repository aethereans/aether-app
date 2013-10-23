function FindOrCreateTopicController($scope, $rootScope, $timeout, frameViewStateBroadcast,
    gateReaderServices, refreshService) {

    $scope.postText = ''
    $scope.postButtonDisabled = true

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
        $scope.allAvailableTopics = filteredTopics

    }

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
        enforceHeaderMaxLength($scope.postText, 60)
        $scope.topicLetterCount = $scope.postText.length
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
        var inputText = angular.element(document.getElementsByClassName('topic-name-entry')).text()
        var topicFound = false
        for (var i = 0; i<$scope.allAvailableTopics.length; i++) {
            if (inputText === $scope.allAvailableTopics[i].Subject) {
                $rootScope.userProfile.selectedTopics.push(
                    $scope.allAvailableTopics[i].PostFingerprint)
                topicFound = true
                break
            }
        }
        if (topicFound === false) {
            gateReaderServices.createTopic(answerArrived, inputText)
            function answerArrived(data) {
                $rootScope.userProfile.selectedTopics.push(data)
            }
        }
        angular.element(document.getElementsByClassName('topic-name-entry')).text('')
        refreshService()


        // This for some reason does not remove the add topic button to its
        // original state. Investigate.
        $scope.createFieldsActive = false

    }

    $scope.$on('refreshPage', function() {
        gateReaderServices.getUppermostTopics(topicsArrived)
    })
    $scope.showCreateFields = function() {
        $scope.createFieldsActive = true
    }

}

FindOrCreateTopicController.$inject = ['$scope', '$rootScope', '$timeout', 'frameViewStateBroadcast',
    'gateReaderServices', 'refreshService']
