// This is the fourth deepest and uppermost API of the front-back end
// communication stack, gateReaderServices. This exposes the data received from
// database into the Angular.js application.

// The API layers are, from deepest to highest, are thus: Hermes, Charon,
// gateService, gateReaderServices.

// Remember, any functionality you're adding there will likely need to be added
// to all layers down below.

angular.module('aether.services')
    .service('gateReaderServices', function(gateService) {
        var parentThis = this
        this.getAllPosts = function(callback, postFingerprint) {
            // This fetches all posts in node in the tree which is not the end.
            // If the requested node is the end node, it behaves like
            // getSinglePost
            gateService.readFetchAllResponse(dataArrived, postFingerprint)

            function dataArrived(data) {
                if (Array.isArray(data)) {
                    var sortedData = aetherSortAlgorithm(data)
                    callback(sortedData)
                }
                else {
                    callback(data)
                }

            }
        }

        this.getAllSpecificDepthPosts = function(callback, postFingerprint, depth) {
            gateService.readFetchSpecificDepthResponse(dataArrived, postFingerprint, depth)

            function dataArrived(data) {
                if (Array.isArray(data)) {
                    var sortedData = aetherSortAlgorithm(data)
                    callback(sortedData)
                }
                else {
                    callback(data)
                }
            }
        }

        this.getSinglePost = function(callback, postFingerprint) {
            // This fetches the single post targeted. Use this to fetch singular
            // subject posts in post flow screen.
            gateService.readFetchSingleResponse(dataArrived, postFingerprint)
            console.log('this is what getsinglepost receives', postFingerprint)

            function dataArrived(data) {
                callback(data)
            }
        }

        this.getDirectDescendantPosts = function(callback, postFingerprint) {
            gateService.readFetchDirectDescendantsResponse(dataArrived, postFingerprint)
            function dataArrived(data) {
                if (data.isArray) {
                    var sortedData = aetherFlatSortAlgorithm(data)
                    callback(sortedData)
                }
                else {
                    callback(data)
                }

            }
        }

        this.getDirectDescendantPostsWithTimespan = function(callback, postFingerprint, daysToSubtract) {
                gateService.readFetchDirectDescendantsWithTimespanResponse(
                    dataArrived, postFingerprint, daysToSubtract)
                function dataArrived(data) {
                    if (data.isArray) {
                        var sortedData = aetherFlatSortAlgorithm(data)
                        callback(sortedData)
                    }
                    else {
                        callback(data)
                    }

                }
        }

        this.countDirectDescendantPosts = function(callback, postFingerprint) {
        gateService.readCountDirectDescendantsResponse(dataArrived, postFingerprint)

        function dataArrived(data) {
            callback(data)
           }
        }

        this.countTopics = function(callback) {
            gateService.readCountTopics(countArrived)
            function countArrived(count) {
                callback(count)
            }
        }

        this.countSubjects = function(callback) {
            gateService.readCountSubjects(countArrived)
            function countArrived(count) {
                callback(count)
            }
        }

        this.countAllDescendantPosts = function(callback, item, postFingerprint) {
        gateService.readCountAllDescendantsResponse(dataArrived, item, postFingerprint)

        function dataArrived(item, data) {
            callback(item, data)
            }
        }

        this.getUppermostTopics = function(callback) {
            gateService.readFetchUppermostTopicsResponse(dataArrived)

            function dataArrived(data) {
                if (data.isArray) {
                    var sortedData = aetherFlatSortAlgorithm(data)
                    callback(sortedData)
                }
                else {
                    callback(data)
                }
            }
        }

        this.getTopmostComment = function(callback, postFingerprint) {
            gateService.readFetchTopmostCommentResponse(dataArrived, postFingerprint)

            function dataArrived(data) {
                callback(data)
            }
        }

        this.getParentPost = function(callback, postFingerprint) {
            parentThis.getSinglePost(originPostArrived, postFingerprint)
            function originPostArrived(data) {
                parentThis.getSinglePost(parentPostArrived, data.ParentPostFingerprint)
                function parentPostArrived(data) {
                    callback(data)
                }
            }
        }

        this.getUserProfile = function(callback) {
            gateService.readGetUserProfile(userProfileArrived)
            function userProfileArrived(data) {
                callback(data)
            }
        }

        this.getSubjectCounts = function(callback, days, fingerprintContainerObject) {
            gateService.readGetSubjectCountsWithTimespanResponse(dataArrived,
                days, fingerprintContainerObject)
            function dataArrived(data) {
                callback(data)
            }
        }

        this.getSavedPosts = function(callback) {
            gateService.readGetSavedPosts(dataArrived)
            function dataArrived(data) {
                callback(data)
            }
        }

        this.countConnectedNodes = function(callback) {
            gateService.readCountConnectedNodes(countArrived)
            function countArrived(count) {
                callback(count)
            }
        }

        this.getLocallyCreatedPosts = function(callback) {
            gateService.readGetLocallyCreatedPosts(dataArrived)
            function dataArrived(data) {
                callback(data)
            }
        }

        this.countReplies = function(callback) {
            gateService.readCountReplies(countArrived)
            function countArrived(count) {
                callback(count)
            }
        }

        this.getReplies = function(callback) {
            gateService.readGetReplies(dataArrived)
            function dataArrived(data) {
                callback(data)
            }
        }

        this.getUnregisteredUserPosts = function(callback, userName) {
            gateService.readGetUnregisteredUserPosts(dataArrived, userName)
            function dataArrived(data) {
                callback(data)
            }

        }

        this.getParentSubjectOfGivenPost = function(callback, postFingerprint) {
            gateService.readGetParentSubjectOfGivenPost(dataArrived, postFingerprint)
            function dataArrived(data) {
                callback(data)
            }
        }

        this.getNewborn = function(callback) {
            gateService.readNewborn(replyArrived)
            function replyArrived(reply) {
                callback(reply)
            }
        }

        this.getOnboardingComplete = function(callback) {
            gateService.readOnboardingComplete(replyArrived)
            function replyArrived(reply) {
                callback(reply)
            }
        }

        this.getAppVersion = function(callback) {
            gateService.readAppVersion(replyArrived)
            function replyArrived(reply) {
                callback(reply)
            }
        }

        this.getUpdateAvailable = function(callback) {
            gateService.readUpdateAvailable(replyArrived)
            function replyArrived(reply) {
                callback(reply)
            }
        }

        this.getOperatingSystem = function(callback) {
            gateService.readGetOperatingSystem(replyArrived)
            function replyArrived(reply) {
                callback(reply)
            }
        }

        this.getSubjects = function(callback, postFingerprint, daysToSubtract) {
            gateService.readGetSubjects(replyArrived, postFingerprint, daysToSubtract)
            function replyArrived(reply) {
                callback(reply)
            }
        }

        this.getHomeScreen = function(callback, numberOfTopics, numberOfSubjects) {
            gateService.readGetHomeScreen(replyArrived, numberOfTopics, numberOfSubjects)
            function replyArrived(reply) {
                callback(reply)
            }
        }

        this.getLastConnectionTime = function(callback) {
            gateService.readGetLastConnectionTime(replyArrived)
            function replyArrived(reply) {
                callback(reply)
            }
        }

        // Below are writes.

        this.createTopic = function(callback, topicName) {
            gateService.writeCreateTopic(answerArrived, topicName)
            function answerArrived(answer) {
                callback(answer)
            }
        }

        this.createPost = function(callback, postSubject, postText, parentPostFingerprint,
            ownerUsername, postLanguage) {
            gateService.writeCreatePost(answerArrived, postSubject, postText, parentPostFingerprint,
                ownerUsername, postLanguage)

            function answerArrived(answer) {
                callback(answer)
                // This should return the fingerprint of the post just created.
            }
        }

        this.createSignedPost = function(callback, postSubject, postText, parentPostFingerprint,
            ownerFingerprint) {
            gateService.writeCreateSignedPost(answerArrived, postSubject, postText, parentPostFingerprint,
                ownerFingerprint)
        }

        this.saveUserProfile = function(callback, userProfileAsObject) {
            gateService.writeSaveUserProfile(answerArrived, userProfileAsObject)
            function answerArrived(answer) {
                callback(answer)
            }
        }

        this.votePost = function(callback, postFingerprint, voteDirection) {
            gateService.writeVotePost(answerArrived, postFingerprint, voteDirection)
            function answerArrived(answer) {
                callback(answer)
            }
        }

        this.savePost = function(callback, postFingerprint) {
            gateService.writeSavePost(answerArrived, postFingerprint)
            function answerArrived(answer) {
                callback(answer)
            }
        }

        this.markAllRepliesAsRead = function(callback) {
            gateService.writeMarkAllRepliesAsRead(answerArrived)
            function answerArrived(answer) {
                callback(answer)
            }
        }

        this.markAllSavedsAsNotSaved = function(callback) {
            gateService.writeMarkAllSavedsAsNotSaved(answerArrived)
            function answerArrived(answer) {
                callback(answer)
            }
        }

        this.quitApp = function() {
            gateService.writeQuitApp()
        }

        this.setOnboardingComplete = function() {
            gateService.writeSetOnboardingComplete()
        }

        this.connectToNodeWithIP = function(ip, port) {
            gateService.writeConnectToNodeWithIP(ip, port)
        }

        this.exportSinglePost = function(callback, post) {
            gateService.writeExportSinglePost(answerArrived, post)
            function answerArrived(answer) {
                callback(answer)
            }
        }

        this.exportAllPosts = function(callback, posts) {
            gateService.writeExportAllPosts(answerArrived, posts)
            function answerArrived(answer) {
                callback(answer)
            }
        }

        // Below are the functions used in aether sort algorithms.

        function numericSortAccordingToVotes(a,b) {
            var sumOfA = a.UpvoteCount - a.DownvoteCount
            var sumOfB = b.UpvoteCount - b.DownvoteCount
            return sumOfB - sumOfA
        }

        function aetherFlatSortAlgorithm(array) {
            // This sorts the input without regard to its depth. So if B is a
            // comment to A and if B has more votes, B will rank higher than A.
            // This is the dumbest sort possible. Use it if you are sure your
            // input is composed of only one level of depth.
            return array.sort(numericSortAccordingToVotes)
        }



        function aetherSortAlgorithm(array) {
            var depthedArray = resolveDepth(array)
            var maxDepth = getDepths(depthedArray)
            var depthSortedArray = []
            for (var i=0;i<=maxDepth;i++) {
                var subsetArray = getSortedSubsetByDepth(depthedArray, i)
                // append this to depth sorted array
                depthSortedArray = depthSortedArray.concat(subsetArray)
            }
            var finalArray = []
            for (var i=0;i<depthSortedArray.length;i++) {
                var indexOfParent = returnIndexOfParentInArray(finalArray,
                    depthSortedArray[i])

                // if parent exists in the keys added so far
                if (indexOfParent !== -1) {
                    finalArray.splice(indexOfParent+1,0,depthSortedArray[i])
                }
                else //if parent doesn't exist in data added so far
                {
                    // insert normally to the next available position.
                    finalArray.push(depthSortedArray[i])
                }
            }
            return finalArray
        }

        function returnIndexOfParentInArray(array, item) {
            var parentFingerprint = item.ParentPostFingerprint
            for (var i=0;i<array.length;i++) {
                if (array[i].PostFingerprint === parentFingerprint) {
                    return i
                }
            }
            return -1 // this only hits on the first run (when given array is empty)

        }

        function getDepths(array) {
            var maxDepth = 0
            for (var i=0;i<array.length;i++) {
                if (array[i].Depth > maxDepth) { maxDepth = array[i].Depth }
            }
            return maxDepth
        }

        function getSortedSubsetByDepth(array, depth) {
            // This needs to return reverse except for the outermost depth,
            // because all the other levels are inserted at [parent+1] loc,
            // which makes the lowest votes inserted last, therefore at top.
            var subset = []
            for (var i=0;i<array.length;i++) {
                if (array[i].Depth === depth) {
                    subset.push(array[i])
                }
            }
            if (depth === 0) { return aetherFlatSortAlgorithm(subset)}
                else { return aetherFlatSortAlgorithm(subset).reverse()}
        }

        function resolveDepth(array) {
            for (var i=0;i<array.length;i++) {
                if (array[i].Depth === undefined) { array[i].Depth = 0}
                markDirectDescendants(array[i].Depth,
                    array, array[i].PostFingerprint)
            }
            return array
        }

        function markDirectDescendants(priorDepth, array, fingerprint) {
        // When given a fingerprint and array and a prior depth, this function
        // searches for posts in array that point the fingerprinted post as
        // parent and increments their depth by 1.
            for (var i=0;i<array.length;i++) {
                if (array[i].ParentPostFingerprint === fingerprint) {
                    array[i].Depth = priorDepth + 1
                }
            }

        }
})

