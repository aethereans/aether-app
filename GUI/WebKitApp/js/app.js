//Main application starting point for frontend aether user interface.

angular.module('aether', [
    'aether.filters',
    'aether.services',
    'aether.directives'])


angular.module('aether.directives', [])
angular.module('aether.services', [])
angular.module('aether.filters', [])

// This is basically modifying angular outside of itself and I have a weird
// feeling about having to do this, but this seems to be the only way of
// being aware of the height change of the window. The watcher for this element
// injected into the scope is on the frameControllers.js.

// I can actually use to get hell a lot of values about the DOM. COOL. COOOOOOL.

function tellAngularOfCurrentDimensions() {
    var domElement = angular.element(document.getElementsByTagName('html'))
    var scope = domElement.scope()
    scope.$apply(function() {
        scope.viewportHeight = document.documentElement.clientHeight
        scope.currentSecondFrameWidth = getStyle(document.getElementById('second-frame'), 'width')
        scope.currentSecondFrameContentsWidth = getStyle(document.getElementById('second-frame-contents'), 'width')
    })
}

document.addEventListener("DOMContentLoaded", tellAngularOfCurrentDimensions, false);
window.onresize = tellAngularOfCurrentDimensions

function getStyle(element, attr) {
    return window.getComputedStyle(element).getPropertyValue(attr)
    //.match(/\d+/)
    // Last part removes the 'px' ending, if you want to use it.

    // So that's how everyone ends up implementing their own jquery...

    // GOTCHA: this doesn't work with getElementsByTagName, for reasons unknown.
}