// This is the third deepest API level, gateService. This receives messages from
// Charon, and decodes them into Javascript objects.

angular.module('aether.services')
    .service('gateService', function() {

    this.readFetchAllResponse = function(callback, postFingerprint) {
        var returnedBlob = JSON.parse(Charon.getAllPosts(postFingerprint))
        var resultArray = []
        for (var i = 0; i<returnedBlob.length; i++) {
            var returnedPost = JSON.parse(returnedBlob[i])
            resultArray.push(returnedPost)
        }
        callback(resultArray)
    }

    this.readFetchSingleResponse = function(callback, postFingerprint) {
        var returnedPost = JSON.parse(Charon.getSinglePost(postFingerprint))
        callback(returnedPost)
    }

    this.readFetchSpecificDepthResponse = function(callback, postFingerprint, depth) {
        var returnedBlob = JSON.parse(Charon.getSpecificDepthPosts(postFingerprint, depth))
        var resultArray = []
        for (var i = 0; i<returnedBlob.length; i++) {
            var returnedPost = JSON.parse(returnedBlob[i])
            resultArray.push(returnedPost)
        }
        callback(resultArray)
    }

    this.readFetchDirectDescendantsResponse = function(callback, postFingerprint) {
        var returnedBlob = JSON.parse(Charon.getDirectDescendantPosts(postFingerprint))
        var resultArray = []
        for (var i = 0; i<returnedBlob.length; i++) {
            var returnedPost = JSON.parse(returnedBlob[i])
            resultArray.push(returnedPost)
        }
        callback(resultArray)
    }

    this.readFetchDirectDescendantsWithTimespanResponse =
        function(callback, postFingerprint, daysToSubstract) {
            var returnedBlob = JSON.parse(Charon.getDirectDescendantPostsWithTimespan(
                postFingerprint, daysToSubstract))
            var resultArray = []
            for (var i = 0; i<returnedBlob.length; i++) {
                var returnedPost = JSON.parse(returnedBlob[i])
                resultArray.push(returnedPost)
            }
            callback(resultArray)
    }


    this.readCountDirectDescendantsResponse = function(callback, postFingerprint) {
        callback(Charon.countDirectDescendantPosts(postFingerprint))
    }

    this.readCountAllDescendantsResponse = function(callback, item, postFingerprint) {
        callback(item, Charon.countAllDescendantPosts(postFingerprint))
    }

    this.readFetchUppermostTopicsResponse = function(callback) {
        var returnedBlob = JSON.parse(Charon.getUppermostTopics())
                var resultArray = []
                for (var i = 0; i<returnedBlob.length; i++) {
                    var returnedPost = JSON.parse(returnedBlob[i])
                    resultArray.push(returnedPost)
                }
                callback(resultArray)
    }

    this.readFetchTopmostCommentResponse = function(callback, postFingerprint) {
        callback(JSON.parse(Charon.getTopmostComment(postFingerprint)))
    }

    this.readGetUserProfile = function(callback) {
        if (Charon.getUserProfile() === '') {
            callback(JSON.parse('{}'))
        }
        callback(JSON.parse(Charon.getUserProfile()))
        // It would be saner if I carried handling empty returns here than
        // having it in Python.
    }

    this.readGetSubjectCountsWithTimespanResponse = function(callback, days, fingerprintContainerObject) {
        callback(JSON.parse(Charon.getSubjectCounts(days,
            JSON.stringify(fingerprintContainerObject))))
    }

    this.readGetSavedPosts = function(callback) {
        var returnedBlob = JSON.parse(Charon.getSavedPosts())
        var resultArray = []
        for (var i = 0; i<returnedBlob.length; i++) {
            var returnedPost = JSON.parse(returnedBlob[i])
            resultArray.push(returnedPost)
        }
        callback(resultArray)
    }

    this.readCountConnectedNodes = function(callback) {
        callback(Charon.countConnectedNodes())
    }

    this.readGetLocallyCreatedPosts = function(callback) {
        var returnedBlob = JSON.parse(Charon.getLocallyCreatedPosts())
        var resultArray = []
        for (var i = 0; i<returnedBlob.length; i++) {
            var returnedPost = JSON.parse(returnedBlob[i])
            resultArray.push(returnedPost)
        }
        callback(resultArray)
    }

    this.readCountReplies = function(callback) {
         callback(Charon.countReplies())
    }

    this.readGetReplies = function(callback) {
        var returnedBlob = JSON.parse(Charon.getReplies())
                var resultArray = []
                for (var i = 0; i<returnedBlob.length; i++) {
                    var returnedPost = JSON.parse(returnedBlob[i])
                    resultArray.push(returnedPost)
                }
                callback(resultArray)
    }

    this.readGetUnregisteredUserPosts = function(callback, userName)  {
            var resultArray = __multiParsify(Charon.getUnregisteredUserPosts(userName))
            callback(resultArray)
    }

    this.readGetParentSubjectOfGivenPost = function(callback, postFingerprint) {
        var resultArray = __multiParsify(Charon.getParentSubjectOfGivenPost(postFingerprint))
        callback(resultArray)
    }

    this.readNewborn = function(callback) {
        callback(Charon.newborn())
    }

    this.readOnboardingComplete = function(callback) {
        callback(Charon.onboardingComplete())
    }

    this.readAppVersion = function(callback) {
        callback(Charon.appVersion())
    }

    this.readGetOperatingSystem = function(callback) {
        callback(Charon.getOperatingSystem())
    }

    this.readUpdateAvailable = function(callback) {
        callback(Charon.updateAvailable())
    }

    this.readGetSubjects = function(callback, postFingerprint, daysToSubtract) {
        var returnedBlob = JSON.parse(Charon.getSubjects(postFingerprint, daysToSubtract))
        var resultArray = []
        for (var i = 0; i<returnedBlob.length; i++) {
            var returnedPost = JSON.parse(returnedBlob[i])
            resultArray.push(returnedPost)
        }
        callback(resultArray)
    }

    this.readGetHomeScreen = function(callback, numberOfTopics, numberOfSubjects) {
        var returnedBlob = JSON.parse(Charon.getHomeScreen(numberOfTopics, numberOfSubjects))
        var resultArray = []
        for (var i=0;i<returnedBlob.length;i++) {
            var returnedDict = returnedBlob[i]
            var returnedSubjects = returnedDict.Subjects
            var resultTopic = {
                'TopicName': returnedDict.TopicName,
                'Subjects': []
            }
            for (var j=0;j<returnedSubjects.length;j++) {
                resultTopic.Subjects.push(JSON.parse(returnedDict.Subjects[j]))
            }
            resultArray.push(resultTopic)
        }
        callback(resultArray)
    }

    this.readGetLastConnectionTime = function(callback) {
        callback(JSON.parse(Charon.getLastConnectionTime()))
    }

    this.readCountTopics = function(callback) {
        callback(Charon.countUppermostTopics())
    }

    this.readCountSubjects = function(callback) {
        callback(Charon.countSubjects())

    }

    function __multiParsify(func, cb) {
        var returnedBlob = JSON.parse(func)
            var resultArray = []
            for (var i = 0; i<returnedBlob.length; i++) {
                var returnedPost = JSON.parse(returnedBlob[i])
                resultArray.push(returnedPost)
            }
            return resultArray
    }



    // Below are write services.

    this.writeCreateTopic = function(callback, topicName) {
        callback(Charon.createTopic(topicName))
    }

    this.writeCreatePost = function(callback, postSubject, postText, parentPostFingerprint,
        ownerUsername, postLanguage) {
        callback(Charon.createPost(postSubject, postText, parentPostFingerprint, ownerUsername, postLanguage))
    }

    this.writeCreateSignedPost = function(callback, postSubject, postText, parentPostFingerprint,
        ownerFingerprint) {
        callback(Charon.createSignedPost(postSubject, postText, parentPostFingerprint,
            ownerFingerprint))
    }

    this.writeSaveUserProfile = function(callback, userProfileAsObject) {
        callback(Charon.writeUserProfile(JSON.stringify(userProfileAsObject)))
    }

    this.writeVotePost = function(callback, postFingerprint, voteDirection) {
        callback(Charon.votePost(postFingerprint, voteDirection))
    }

    this.writeSavePost = function(callback, postFingerprint) {
        callback(Charon.savePost(postFingerprint))
    }

    this.writeMarkAllRepliesAsRead = function(callback) {
        callback(Charon.markAllRepliesAsRead())
    }

    this.writeMarkAllSavedsAsNotSaved = function(callback) {
        callback(Charon.markAllSavedsAsNotSaved())
    }

    this.writeQuitApp = function() {
        Charon.quitApp()
    }

    this.writeSetOnboardingComplete = function() {
        Charon.setOnboardingComplete()
    }

    this.writeConnectToNodeWithIP = function(ip, port) {
        Charon.connectToNodeWithIP(ip, port)
    }

    this.writeExportSinglePost = function(callback, post) {
        post = JSON.stringify(post)
        callback(Charon.exportSinglePost(post))
    }

    this.writeExportAllPosts = function(callback, posts) {
        posts = JSON.stringify(posts)
        callback(Charon.exportAllPosts(posts))
    }

})

// OK, so I don't know the difference between callback(resultArray) and
// the one without return.. It works without return, too, so I'm removing item
// but I would like to learn what exactly is happening there.