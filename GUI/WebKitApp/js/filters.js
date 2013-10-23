angular.module('aether.filters',[]).
    filter('positiveNumber', function() {
        return function(num) {
            if (num <= 0) { return 0 }
            else { return num }
        }
    })


.filter('newlines', function() {
        return function(text) {
            return text.replace(/\n/g, '<br>')
        }
    })


.filter('truncate', function() {
        return function(text, length, end) {
            if (isNaN(length)) {
                length = 10
            }
            if (end === undefined) {
                end = '...'
            }
            if (text === undefined) {
                return text
            }
            if (text.length - end.length <= length) {
                return text
            }
            else {
                return String(text).substring(0, length - end.length) + end
            }
        }
    })