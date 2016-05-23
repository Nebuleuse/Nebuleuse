angular.module('RDash').directive('nebCommit', nebCommit);


function nebCommit() {
    var directive = {
        restrict: 'E',
        scope: {
            commit: '='
        },
        templateUrl: 'templates/updates/commits.html'
    };
    return directive;
};

