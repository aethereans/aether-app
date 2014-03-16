function SavedItemsController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices, refreshService) {
    gateReaderServices.getSavedPosts(dataArrived)
    function dataArrived(data) {
        $scope.savedPosts = data
        //console.log($scope.savedPosts)
    }

    $scope.exportAllPostsClick = function() {
        gateReaderServices.exportAllPosts(answerArrived, $scope.savedPosts)
        function answerArrived(answer) {
            if (answer) {
                $scope.alreadySaved = true
            }

        }
    }

    $scope.cleanAllSavedItems = function() {
        console.log('remove all saved marks.')
        gateReaderServices.markAllSavedsAsNotSaved(answerArrived)
        function answerArrived(answer) {
            refreshService()
        }
    }

    $scope.$on('refreshPage', function() {
        gateReaderServices.getSavedPosts(dataArrived)
        function dataArrived(data) {
            $scope.savedPosts = data
            //console.log($scope.savedPosts)
        }
    })
}
SavedItemsController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices', 'refreshService']