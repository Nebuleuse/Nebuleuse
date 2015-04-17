angular.module('RDash')
    .controller('LogCtrl', ['$scope', '$http', LogCtrl]);

function LogCtrl($scope, $http) {
	$scope.logLines = "";
	$scope.setPageTitle("Live Log");
	$scope.subscribeTo("logEvent");
	$scope.$on("logEvent", function (event, arg) {
		$scope.logLines += arg;
	})
}