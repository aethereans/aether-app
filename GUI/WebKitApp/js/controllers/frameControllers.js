// First, second and third frame controllers handle the resize events.
// They also handle global events that result in change of pages, buttons etc.

function FirstFrameController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices) {


    // Methods below are available to all scopes. They're inside the first frame
    // controller, because it seems I cannot place rootScope
    // objects outside one! They're still available everywhere though.

    // Yes, I know rootScope is evil. Let me know if you have a better solution.

    $rootScope.log = function(variable) {
        console.log(variable)
    }
    $rootScope.alert = function(text) {
        alert(text)
    }
    // These two methods are available globally: they allow basic frame change
    // methods to work. Frame change changes the partial views, and frame size
    // change changes the sizes of the second and third frames.

    // IMPORTANT: If you pass "" as the 3rd value of changeState, it will RETAIN
    // the value of requestedId in memory so you can still access it from the
    // new partial. If you pass undefined as the 3rd value, it will FLUSH the
    // existing value from rootScope, so no one has access to the value anymore.

    // Ergo, be careful in passing undefined, as in one page there are probably
    // more than one partial making use of that value. The only use case I can
    // see for flushing the value so far is the entry, where the value is
    // empty from the start, but if left as "", it is mistaken as "retain"
    // rather than explicit empty.

    // If you want to use undefined, be careful, most controllers actively
    // filter out undefined to keep the values within safe from outside
    // interference.

    $rootScope.changeState = function(secondFrame, thirdFrame, id) {
        frameViewStateBroadcast.receiveState(secondFrame, thirdFrame, id)
    }

    // This is a general listener on view state change that makes requestedId
    // variable available on the scope on all possible controllers.

    $rootScope.$on('frameViewStateChanged', function() {
        if (frameViewStateBroadcast.id !== "") {
            $rootScope.requestedId = frameViewStateBroadcast.id
        }
    })

    $rootScope.$watch('viewportHeight', function() {
        $rootScope.viewportStyle = {
            'height': $scope.viewportHeight + 'px'
        }
    })

    $rootScope.clickToUsername = function(userName) {
        $rootScope.requestedUsername = userName
        $rootScope.thirdFrameCSSStyle = {
            "display":"none"
        }
        $scope.changeState('unregisteredUserFeed', '', '')
    }

    // These three things are how you manipulate the frames in DOM rather
    // than jQuery. Strictly speaking I don't need this, but it's mostly here
    // as a reminder. (Because evaluation on DOM is graceful to undefined)

    // These are default states at the beginning of the application.

    $rootScope.firstFrameCSSStyle = {}

    $rootScope.secondFrameCSSStyle = {
        'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
        (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2)) + 'px'
    }

    $rootScope.thirdFrameCSSStyle = {
        "display":"none"
    }

    // Below are the first call to userProfile object to make it globally
    // available at application start. The watch below will fire at every change.

    // The watch might need an object equality check rather than default check
    // in the case this does not update, add true as the last parameter of the
    // watch.

    // The fact that I could see this and get this to work before I made a million
    // save requests everywhere points to that being a rare moment of brilliance.

    gateReaderServices.getUserProfile(
        function(userProfile) {
            $rootScope.userProfile = userProfile
        })

    // If any part of the userProfile doesn't exist on first load, create it here.

    if ($rootScope.userProfile.selectedTopics === undefined) {
        $rootScope.userProfile.selectedTopics = []
    }

    if ($rootScope.userProfile.UserDetails === undefined) {
        $rootScope.userProfile.UserDetails = {
            Username: '',
            UserLanguages: ['English'],
            StartAtBoot: true,
            Logging: true,
        }
    }

    if ($rootScope.userProfile.UnreadReplies === undefined) {
        $rootScope.userProfile.UnreadReplies = []
    }

    if ($rootScope.userProfile.ReadReplies === undefined) {
        $rootScope.userProfile.ReadReplies = []
    }

    $rootScope.appIsPaused = false


    $scope.$watch('userProfile', function() {
        gateReaderServices.saveUserProfile(function() {}, $rootScope.userProfile)
        console.log('autocommit fired, new profile: ',$rootScope.userProfile)
    }, true)


    // This is a helper function to check if the pane queried is currently
    // active.

    $rootScope.isActive = function(paneName) {
        return frameViewStateBroadcast.secondFrame == paneName ||
        frameViewStateBroadcast.thirdFrame == paneName ?
        true : false
    }

    gateReaderServices.getAppVersion(function(reply) {
        $rootScope.appVersion = reply.slice(0,1) + '.' + reply.slice(1,2) + '.' + reply.slice(2,3)
    })

    $scope.homeButtonClick = function() {
        $scope.changeState('homeFeed', 'topicGroupsFeedLite', '')
        $rootScope.secondFrameCSSStyle = {
            'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
            (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2)) + 'px'
        }
        $rootScope.thirdFrameCSSStyle = {
            'display':'none'
        }
    }

    $scope.topicsButtonClick = function() {
        $scope.changeState('subjectsFeed', 'topicsFeedLite', undefined)
        $rootScope.thereIsASelectedSubject = false

        $rootScope.secondFrameCSSStyle = {}
        $rootScope.thirdFrameCSSStyle = {}
    }

    $scope.savedItemsButtonClick = function() {
        $scope.changeState('savedItemsFeed', 'savedItemsFeedLite', '')
        $rootScope.secondFrameCSSStyle = {
            'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
            (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2)) + 'px'
        }
        $rootScope.thirdFrameCSSStyle = {
            'display':'none'
        }
    }

    $scope.settingsButtonClick = function() {
        $scope.changeState('settingsFeed', 'settingsSelectorFeedLite', undefined)
        $rootScope.thirdFrameCSSStyle = {
            'width': '150px'
        }
        $rootScope.secondFrameCSSStyle = {
            'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
            (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2))
            - (angular.element(document.getElementById('third-frame')).css('width').slice(0, -2)) + 'px'
        }

    }


}

FirstFrameController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
'gateReaderServices']



function SecondFrameController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices) {

    // This is the list of partials this scope can load. This is also one of
    // the key reasons the first / second / third frame controllers cannot be
    // removed.

    // REGIONS describe the states application is in. Region change is
    // required for first frame buttons' awareness, but could be useful in other
    // places.

    $scope.templates = [{
        name: 'postsFeed',
        url: 'partials/postsFeed.html',
        region: 'Details'
    }, {
        name: 'subjectsFeed',
        url: 'partials/subjectsFeed.html',
        region: 'Details'
    }, {
        name: 'topicsFeed',
        url: 'partials/topicsFeed.html',
        region: 'Details'
    }, {
        name: 'savedItemsFeed',
        url: 'partials/savedItemsFeed.html',
        region: 'SavedItems'
    }, {
        name: 'homeFeed',
        url: 'partials/homeFeed.html',
        region: 'Home'
    }, {
        name: 'userProfile',
        url: 'partials/userProfile.html'
    }, {
        name: 'createSubject',
        url: 'partials/createSubject.html',
        region: 'Details'
    }, {
        name: 'findOrCreateTopic',
        url: 'partials/findOrCreateTopic.html',
        region: 'Details'
    }, {
        name: 'settingsFeed',
        url: 'partials/settingsBase.html',
        region: 'Settings'
    }, {
        name: 'singlePost',
        url: 'partials/singlePost.html',
        region: 'SavedItems'
    }, {
        name: 'repliesFeed',
        url: 'partials/repliesFeed.html',
        region: 'Home'
    }, {
        name: 'singleReply',
        url: 'partials/singlePost.html',
        region: 'Home'
    }, {
        name: 'unregisteredUserFeed',
        url: 'partials/unregisteredUserItemsFeed.html',
        region: 'Home' // This is for now..
    },{
        name:'onboarding',
        url: 'partials/onboarding.html',
        region: 'NoRegion'
    }]

    gateReaderServices.getOnboardingComplete(onboardingCompleteArrived)
    function onboardingCompleteArrived(onboardingComplete) {
        $rootScope.onboardingComplete = onboardingComplete
        console.log('onboardingComplete status:', onboardingComplete)
        if (onboardingComplete === true) {
            // If not newborn
            $scope.selectedTemplate = $scope.templates[4] //default 4
            $rootScope.currentApplicationRegion = 'Home'
        }
        else
        {
            $scope.selectedTemplate = $scope.templates[13] //default 4
            $rootScope.currentApplicationRegion = 'NoRegion'
        }
    }


    // This part below fires on frame state change and greps for the incoming
    // frame name in templates. If I built this logic, this is smart, but I have
    // no recollection of writing it. Oh god.

    $scope.$on('frameViewStateChanged', function() {
        var searchResult
        for (var i = 0; i<$scope.templates.length; i++) {
            if ($scope.templates[i].name === frameViewStateBroadcast.secondFrame) {
                searchResult = $scope.templates[i]
            }
        }

        // Change only if 2nd variable exists i.e. it isn't ""
        if (frameViewStateBroadcast.secondFrame !== "") {
            $scope.selectedTemplate = searchResult
            $rootScope.currentApplicationRegion = searchResult.region
        }

    })
}

SecondFrameController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
'gateReaderServices']

function ThirdFrameController($scope, $rootScope, frameViewStateBroadcast) {
    $scope.templates = [{
        name: 'subjectsFeedLite',
        url: 'partials/subjectsFeedLite.html'
    }, {
        name: 'topicsFeedLite',
        url: 'partials/topicsFeedLite.html'
    }, {
        name: 'createLite',
        url: 'partials/createLite.html'
    }, {
        name: 'settingsSelectorFeedLite',
        url: 'partials/settingsSelectorFeedLite.html'
    }  ]

    $scope.selectedTemplate = $scope.templates[0]

    $scope.$on('frameViewStateChanged', function() {
        var searchResult
        for (var i = 0; i<$scope.templates.length; i++) {
            if ($scope.templates[i].name === frameViewStateBroadcast.thirdFrame) {
                searchResult = $scope.templates[i]
            }
        }

        // Change only if 2nd variable exists i.e. it isn't ""
        if (frameViewStateBroadcast.thirdFrame !== "") {
            $scope.selectedTemplate = searchResult
        }
    })
}

ThirdFrameController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast']