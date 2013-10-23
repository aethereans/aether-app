function SubjectsController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices) {

    $scope.$watch('requestedId + timespan', function() {
        if ($rootScope.requestedId !== undefined) {
            // This if is needed because otherwise the first load calls an undefined
            // requestedId, forcing the app to load topics to memory when not needed.

            // Here we are doing a second level of data insertion, calculation of
            // comment counts for each of the subjects. This could be handled in
            // the API, but then all calls would have this property, needlessly
            // slowing things down when it is not needed.
            if ($scope.timespan === undefined || $scope.timespan === 0) {
                //console.log("timespan undefined or zero")
                gateReaderServices.getSubjects(subjectsArrived,
                    $rootScope.requestedId, 180) // six months.

                $scope.timespanIsFiltered = false
            }
            else {
                //console.log("timespan NOT undefined or zero")
                gateReaderServices.getSubjects(subjectsArrived,
                    $rootScope.requestedId, $scope.timespan)

                $scope.timespanIsFiltered = true
            }

            function subjectsArrived(data) {
                console.log('subjects arrived!')
                console.log(data)
                var column1 = []
                var column2 = []
                var column3 = []
                for (var i=0;i<data.length;i++) {
                    if (i%3===0) {
                        column1.push(data[i])
                    }
                    else if(i%3===1) {
                        column2.push(data[i])
                    }
                    else {
                        column3.push(data[i])
                    }
                }
                $scope.subjectsCol1 = column1
                $scope.subjectsCol2 = column2
                $scope.subjectsCol3 = column3

                var length = column1.length + column2.length + column3.length
                // This allows me to hide the unsightly page in case a topic has
                // no available subjects in it to show a good sadface.
                console.log("subjects in scope ", $scope.subjects)
                if(length === 0) {
                    $rootScope.noSubjectsAvailable = true
                }
                else {
                    $rootScope.noSubjectsAvailable = false
                }

            }

            gateReaderServices.getSinglePost(topicArrived, $rootScope.requestedId)
            function topicArrived(data) {
                $scope.selectedTopic = data
            }
        }
        else
        {
            $rootScope.noSubjectsAvailable = true
        }
    })


    // TODO: Here, check if the controller can access 2ndframeid. If it can, it's
    // at 3rd frame. In that case, write a ng-class directive that can look for
    // item id and activate on match so the one lights up.

    // I also need to determine here if this is in 2nd or 3rd scope, or write a
    // different controller for lite version.

}
SubjectsController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']
