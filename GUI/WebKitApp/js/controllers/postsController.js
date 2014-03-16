function PostsController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices) {

    $scope.$watch('requestedId', function() {
        if ($rootScope.requestedId !== undefined) {
            getPosts($rootScope.requestedId)
            getSubject($rootScope.requestedId)
        }
    })

    $scope.$on('refreshPage', function() {
        // This doesn't seem to work.. old stuff doesn't get cleared.
        getPosts($rootScope.requestedId)
        getSubject($rootScope.requestedId)
    })

    var getPosts = function(requestedId) {
        gateReaderServices.getAllPosts(postsArrived, $rootScope.requestedId)
        function postsArrived(data) {
            $scope.feed = data
        }
    }

    var getSubject = function(requestedId) {
        gateReaderServices.getSinglePost(subjectArrived, $rootScope.requestedId)
        function subjectArrived(data) {
            $scope.subject = data
            }
    }

    $scope.openReplyPane = function(fingerprint) {
        console.log("openreplypane has been called")
        $rootScope.targetPost = fingerprint
        frameViewStateBroadcast.receiveState("", "createLite", $rootScope.requestedId)
        $rootScope.thirdFrameCSSStyle = {
            'display':'block',
            'width':'393px'
        }

        // The value from the third frame is called from the root scope element
        // because angular.element call isn't fast enough to reflect the new
        // value.

        $rootScope.secondFrameCSSStyle = {
            'width': getStyle(document.getElementById('root-body'),'width').slice(0, -2) -
            $rootScope.thirdFrameCSSStyle.width.slice(0, -2) -
            getStyle(document.getElementById('first-frame'), 'width').slice(0, -2)
            + 'px'
            // This looks a little broken.. convert this to getstyle
        }
        $rootScope.firstFrameCSSStyle = {}


    }

    // You see these two methods below? Be VERY careful in changing them, these
    // methods have exact duplicates in the backend. Because of the architecture
    // I can send the signal upon user input and make it affect the database
    // record, but it is nigh impossible to send a signal from backend that
    // will tell the frontend to update a certain part. The only way backend
    // can do stuff in the frontend is inbetween page switches (reloads), and I
    // don't want to reload on every vote, so these functions do the manual
    // modification upon receiving the OK signal from the backend.

    // This is a pile of if ... s, esp. in the backend, so don't touch this
    // unless you have an hour to work on that.

    $scope.upvotePost = function(postFingerprint) {
        gateReaderServices.votePost(answerArrived, postFingerprint, 1)
        function answerArrived(answer) {
            if (answer === true) {
                var post = getPostFromScope(postFingerprint)
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
                var post = getPostFromScope(postFingerprint)
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
                var post = getPostFromScope(postFingerprint)
                post.Saved === true ? post.Saved = false : post.Saved = true
            }
        }
    }

    function getPostFromScope(postFingerprint) {
        // This looks at the current posts **available in frontend** and tries
        // to find it without communicating with the backend.
        var post
        for (var i = 0; i<$scope.feed.length; i++) {
            (function(i) {
                if ($scope.feed[i].PostFingerprint === postFingerprint) {
                    post = $scope.feed[i]
                }})(i)
        }
        if (post === undefined) {
            // If not found, i.e. it's a subject
            post = $scope.subject
        }
        return post
    }

    $scope.toggleSidebar = function() {

        // When the screen is made full screen, the return screen should be either
        // cached, so it will return to that after fullscreen, or it should show
        // the vanilla subjects feed.

        // For now, I'll just show the default view (subjects lite feed) but
        // in the future I should copy those CSS style objects and return to the
        // cached copy when exiting fullscreen.

        $scope.changeState('','subjectsFeedLite', '')

        // This view will allow arbitrary class additions in the DOM without
        // mudding the controller here. I also probably should keep track of
        // the root scope, I have hella lot of variables mingling out in the
        // wild..

        if ($rootScope.isFullscreen === undefined || $rootScope.isFullscreen === false) {
            $rootScope.isFullscreen = true

            $rootScope.firstFrameCSSStyle = {"display": "none"}
            $rootScope.secondFrameCSSStyle = {"width": angular.element(
                document.getElementsByTagName('body')).css("width")}
            $rootScope.thirdFrameCSSStyle = {"display": "none"}


        }
        else if ($rootScope.isFullscreen === true) {
            $rootScope.isFullscreen = false

            $rootScope.firstFrameCSSStyle = {}
            $rootScope.secondFrameCSSStyle = {}
            $rootScope.thirdFrameCSSStyle = {"display":"block"}

        }
    }

}
PostsController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
    'gateReaderServices']


