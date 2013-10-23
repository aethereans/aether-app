function SavedItemsController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices) {
    gateReaderServices.getSavedPosts(dataArrived)
    function dataArrived(data) {
        $scope.savedPosts = data
        //console.log($scope.savedPosts)
    }
}
SavedItemsController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']