angular
    .module('RDash')
    .directive('nebDifflist', nebDifflist);

function nebDifflist() {
    var directive = {
        restrict: 'E',
        scope: {
            diffs: '='
        },
        templateUrl: 'templates/updates/diffList.html'
    };
    return directive;
};