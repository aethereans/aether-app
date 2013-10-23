function SinglePostController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices, $timeout) {

    $scope.postButtonDisabled = true
    $scope.postText = ''

    gateReaderServices.getSinglePost(postArrived, $rootScope.requestedId)
    function postArrived(p) {
        $scope.post = p
        gateReaderServices.getParentSubjectOfGivenPost(parentArrived, p.PostFingerprint)
        function parentArrived(parent) {
            $scope.subject = parent[0]
        }
        if (p.Body.length < 50) {
            $scope.centerLayout = true
        }
    }

    $scope.upvotePost = function(postFingerprint) {
        gateReaderServices.votePost(answerArrived, postFingerprint, 1)
        function answerArrived(answer) {
            if (answer === true) {
                var post = $scope.post
                if (!post.Upvoted) {
                    // If post is not already upvoted
                    post.UpvoteCount += 1
                    post.Upvoted = true
                    if (post.Downvoted) {
                        post.DownvoteCount -= 1
                        post.Downvoted = false
                    }
                }
                else {
                    // If it is already upvoted, remove the upvote.
                    gateReaderServices.votePost(answerArrived, postFingerprint, 2)
                    function answerArrived(answer) {
                        post.UpvoteCount -= 1
                        post.Upvoted = false
                    }
                }
            }
        }
    }

    $scope.downvotePost = function(postFingerprint) {
        gateReaderServices.votePost(answerArrived, postFingerprint, -1)
        function answerArrived(answer) {
            if (answer === true) {
                var post = $scope.post
                if (!post.Downvoted) {
                    // If this is a new downvote
                    post.DownvoteCount += 1
                    post.Downvoted = true
                    if (post.Upvoted) {
                        post.UpvoteCount -= 1
                        post.Upvoted = false
                    }
                }
                else {
                    // if already voted
                    gateReaderServices.votePost(answerArrived, postFingerprint, 2)
                    function answerArrived(answer) {
                        post.DownvoteCount -= 1
                        post.Downvoted = false
                    }
                }
            }
        }
    }

    $scope.savePost = function(postFingerprint) {
        gateReaderServices.savePost(answerArrived, postFingerprint)
        function answerArrived(answer) {
            if (answer === true) {
                var post = $scope.post
                post.Saved === true ? post.Saved = false : post.Saved = true
            }
        }
    }

    $scope.replyButtonClick = function() {
        $scope.replyPaneOpen = $scope.replyPaneOpen ? false : true
        var metaBox = angular.element(document.getElementsByClassName('meta-box'))
        $timeout(function() {
            metaBox[0].scrollIntoViewIfNeeded()
        }, 10)
        $scope.subjectBodyEntryStyle = {
            'min-height': ($scope.viewportHeight * 0.8) + 'px'
        }

    }

    $scope.submitButtonClick = function() {

        // check the username.
        if ($rootScope.userProfile.UserDetails.Username === '') {
            $rootScope.userProfile.UserDetails.Username = 'no name given'
        }

        $scope.replyPaneOpen = false
        var content = angular.element(document.getElementsByClassName('subject-body-entry'))[0].innerText
        // IMPORTANT: LANGUAGE SANITY CHECK, the lang needs to exist in selected user langs.

        gateReaderServices.createPost(answerArrived, '', content, $scope.post.PostFingerprint,
            $rootScope.userProfile.UserDetails.Username, $scope.post.Language)

        function answerArrived(createdPostFingerprint) {
            console.log('The user has replied to one of his / her replies. The fingerprint of the new reply is: ', createdPostFingerprint)
        }
    }

    $scope.$watch('postText', function() {
        if ($scope.postText.length > 5 && $scope.postText.length < 60000) {
            $scope.postButtonDisabled = false
        }
        else {
            $scope.postButtonDisabled = true
        }
    })

    var subjectBodyEntry = angular.element(document.getElementsByClassName('subject-body-entry'))
    var metaBox = angular.element(document.getElementsByClassName('meta-box'))
    $scope.subjectBodyEntryStyle = {}

    subjectBodyEntry.bind('paste keydown', function() {
        $timeout(function() {
            metaBox[0].scrollIntoViewIfNeeded()
            $scope.bodyLetterCount = subjectBodyEntry.text().length
            $scope.bodyWordCount = subjectBodyEntry.text().split(/\s+/).length - 1
        }, 10)
    })

}
SinglePostController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
    'gateReaderServices', '$timeout']
