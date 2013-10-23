function SettingsController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices) {

    // This two are on rootScope because the third frame needs to access this too.
    $rootScope.settingsSubmenuTemplates = [
    {
        Name: 'Basics',
        Url: 'partials/settings/basics.html'
    }, {
        Name: 'Identity',
        Url: 'partials/settings/identity.html'
    }, {
        Name: 'About',
        Url: 'partials/settings/about.html'
    }, {
        Name: 'Licenses',
        Url: 'partials/settings/licenses.html'
    },

    ]
    // This two are on rootScope because the third frame needs to access this too.
    $rootScope.settingsSelectedTemplate = $scope.settingsSubmenuTemplates[0]

    $scope.ClickSettingsSubmenu = function(targetMenuName) {
        for (var i=0; i<$scope.settingsSubmenuTemplates.length;i++) {
            if ($scope.settingsSubmenuTemplates[i].Name === targetMenuName) {
                $rootScope.settingsSelectedTemplate = $scope.settingsSubmenuTemplates[i]
            }
        }
    }

    for (var i=0;i<$rootScope.userProfile.UserDetails.UserLanguages.length;i++) {
        var language = $rootScope.userProfile.UserDetails.UserLanguages[i]
        $scope[language+'Selected'] = true
    }

    // Start at boot radio button.

    if ($rootScope.userProfile.UserDetails.StartAtBoot) {
        $scope.startAtBootRadioState = [
        {'checked':'checked', 'value':1, 'name':'YES'},
        {'checked':'', 'value':0, 'name':'NO'} ]
        // Setting the current state retrieved from JSON.
    }
    else {
        $scope.startAtBootRadioState = [
        {'checked':'', 'value':1, 'name':'YES'},
        {'checked':'checked', 'value':0, 'name':'NO'} ]
    }

    $scope.setStartAndBootStatus = function(value) {
        $rootScope.userProfile.UserDetails.StartAtBoot = !!value
    }

    // Logs radio button.

    if ($rootScope.userProfile.UserDetails.Logging) {
        $scope.loggingRadioState = [
        {'checked':'checked', 'value':1, 'name':'YES'},
        {'checked':'', 'value':0, 'name':'NO'} ]
        // Setting the current state retrieved from JSON.
    }
    else {
        $scope.loggingRadioState = [
        {'checked':'', 'value':1, 'name':'YES'},
        {'checked':'checked', 'value':0, 'name':'NO'} ]
    }

    $scope.setLoggingStatus = function(value) {
        $rootScope.userProfile.UserDetails.Logging = !!value
    }


    $scope.langClick = function(language) {
        $scope[language+'Selected']
        console.log($scope[language+'Selected'])
        if ($scope[language+'Selected']) {
            if ($rootScope.userProfile.UserDetails.UserLanguages.length > 1) {
                $scope[language+'Selected'] = false
                var index = $rootScope.userProfile.UserDetails.UserLanguages.indexOf(language)
                if (index > -1) {
                    $rootScope.userProfile.UserDetails.UserLanguages.splice(index, 1);
                }
            }
        }
        else
        {
            $scope[language+'Selected'] = true
            $rootScope.userProfile.UserDetails.UserLanguages.push(language)
        }
    }
    $scope.okButtonDisabled = false

    $scope.connectToNode = function() {
        var ipText = angular.element(document.getElementsByClassName('text-entry-ip'))[0].value
        var portText = angular.element(document.getElementsByClassName('text-entry-port'))[0].value

        console.log('connecttonode is called')
        if (ipText && portText) {
            gateReaderServices.connectToNodeWithIP(ipText, portText)
            $scope.okButtonDisabled = true
            angular.element(document.getElementsByClassName('text-entry-ip'))[0].value = ''
            angular.element(document.getElementsByClassName('text-entry-port'))[0].value = ''

        }
    }

}
SettingsController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']