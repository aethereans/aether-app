function CreateLiteController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices, refreshService) {

    $scope.postButtonDisabled = true
    $scope.postText = ''

    var cachedSecondFrameCSSStyle, cachedThirdFrameCSSStyle

    $scope.closeReplyPane = function() {
        frameViewStateBroadcast.receiveState("", "subjectsFeedLite", $rootScope.requestedId)

        $rootScope.secondFrameCSSStyle = {
            'width': '850px'
        }
        $rootScope.thirdFrameCSSStyle = {
            'width': '243px'
        }
    }

    $scope.postReply = function(postText) {

        // check the username.
        if ($rootScope.userProfile.UserDetails.Username === '') {
            $rootScope.userProfile.UserDetails.Username = 'no name given'
        }

        var content = $scope.postText
        gateReaderServices.createPost(answerArrived, '', content, $rootScope.targetPost,
            $rootScope.userProfile.UserDetails.Username, $scope.targetPostObject.Language)
        function answerArrived(createdPostFingerprint) {
            console.log("this is the id of the currently created post: ",createdPostFingerprint)
            frameViewStateBroadcast.receiveState("", "subjectsFeedLite", $rootScope.requestedId)
            $rootScope.secondFrameCSSStyle = {
                'width': '850px'
            }
            $rootScope.thirdFrameCSSStyle = {
                'width': '243px'
            }
            refreshService()
        }
    }

    $scope.togglePaneSize = function () {

        if ($rootScope.createPaneIsFullscreen === undefined ||
            $rootScope.createPaneIsFullscreen === false) {
            $rootScope.createPaneIsFullscreen = true
            // First, cache the styles,
            cachedSecondFrameCSSStyle = $rootScope.secondFrameCSSStyle
            cachedThirdFrameCSSStyle = $rootScope.thirdFrameCSSStyle
            // then assign new ones.
            $rootScope.secondFrameCSSStyle = {
                'display':'none'
            }
            $rootScope.thirdFrameCSSStyle = {
                'width':'100%'
            }
        }
        else {
            $rootScope.createPaneIsFullscreen = false
            // Returning to normal, call the old data from caches,
            $rootScope.secondFrameCSSStyle = cachedSecondFrameCSSStyle
            $rootScope.thirdFrameCSSStyle = cachedThirdFrameCSSStyle
            // and invalidate the caches.
            cachedSecondFrameCSSStyle = {}
            cachedThirdFrameCSSStyle = {}
        }
    }

    $scope.$watch('targetPost', function() {
        gateReaderServices.getSinglePost(subjectArrived, $rootScope.targetPost)
    })

    $scope.$watch('postText', function() {
        if ($scope.postText.length > 5 && $scope.postText.length < 60000) {
            $scope.postButtonDisabled = false
        }
        else {
            $scope.postButtonDisabled = true
        }
    })

    function subjectArrived(data) {
        $scope.targetPostObject = data
        if (data.Body === "") {
            $scope.parentIsSubject = true
        }
        else {
            $scope.parentIsSubject = false
        }
    }
}

CreateLiteController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
    'gateReaderServices', 'refreshService']